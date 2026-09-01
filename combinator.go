package chomp

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

// Result is the expected output from a [Combinator].
type Result interface {
	string | []string
}

// Combinator is a higher-order function capable of parsing text under a defined
// condition. Combinators can be combined to form more complex parsers. Upon success,
// a combinator will return both the unparsed and parsed [State]. All combinators are
// strict and must parse its input. Any failure to do so should raise a [CombinatorParseError].
//
// Raw invocation c(state) is the composition surface used by combinators calling
// combinators (and custom drivers); its errors are unfinalised and may still
// reference the live input. Use [Combinator.Run] as the string-in/string-out entry
// point, which finalises any error via [Finalize] before returning it.
type Combinator[T Result] func(State) (State, T, error)

// run is shared by [Combinator.Run] and [MappedCombinator.Run].
func run[T any](c func(State) (State, T, error), input string) (string, T, error) { //nolint:ireturn // T is caller-constrained (Result or any wrapped by MappedCombinator), not an open interface.
	rem, ext, err := c(NewState(input))
	if err != nil {
		return rem.Rest(), ext, Finalize(err)
	}

	return rem.Rest(), ext, nil
}

// Run applies c against input, the sole string-in/string-out entry point.
// Unlike raw invocation c(state), any returned error is passed through
// [Finalize] first.
func (c Combinator[T]) Run(input string) (string, T, error) { //nolint:ireturn // T is closed over Result (string | []string), not an open interface.
	return run[T](c, input)
}

// Finalize prepares an error returned by a raw combinator invocation for use
// outside a parse. It is the boundary [Combinator.Run] applies errors to
// before returning them; custom drivers built on raw c(state) invocation
// should apply it at their own boundary. It is currently the identity
// function - a reserved seam for future work bounding the memory an escaped
// error retains.
func Finalize(err error) error {
	return err
}

// CombinatorParseError defines an error that is raised when a combinator
// fails to parse the input text under its expected condition.
type CombinatorParseError struct {
	// Expected describes, in human terms, what the combinator required to
	// succeed. Empty when no single literal applies, in which case
	// [CombinatorParseError.Error], [CombinatorParseError.Snippet] and
	// [CombinatorParseError.LogValue] fall back to a description derived
	// from Type.
	Expected string

	// State at the point of failure, capturing the absolute position
	// within the original input alongside the unparsed suffix.
	State State

	// Labels records the chain of grammar-level names (see the future
	// Label combinator) being parsed when the failure occurred, outermost
	// first. Always empty until a caller wraps a combinator with Label.
	Labels []string

	// Type of [Combinator] that failed.
	Type string
}

// typeFallback maps a Type with no natural literal to quote, and no "is_"
// predicate prefix, to a human phrase completing "expected ___".
var typeFallback = map[string]string{
	"satisfy":           "a matching character",
	"any_char":          "any character",
	"take":              "more input",
	"escaped":           "an escape sequence",
	"escaped_transform": "an escape sequence",
	"verify":            "a value satisfying the predicate",
	"eof":               "end of input",
	"all_consuming":     "end of input",
	"peek_not":          "no match",
	"first":             "one of the alternatives",
	"many_n":            "further progress",
	"many_till":         "further progress",
	"many_till_0":       "further progress",
	"fold_many":         "further progress",
	"many_count":        "further progress",
}

// expected resolves the human phrase completing "expected ___": Expected
// if set, otherwise a description derived from Type.
func (e CombinatorParseError) expected() string {
	if e.Expected != "" {
		return e.Expected
	}
	if after, ok := strings.CutPrefix(e.Type, "is_"); ok {
		return strings.ReplaceAll(after, "_", " ")
	}
	if phrase, ok := typeFallback[e.Type]; ok {
		return phrase
	}
	return strings.ReplaceAll(e.Type, "_", " ")
}

// Error returns a single-line, grep-stable string representation of the
// current error: "chomp: parse error at line %d, column %d (offset %d):
// expected %s". It never contains a newline. For a human-facing,
// caret-annotated view of the failure, see [CombinatorParseError.Snippet].
func (e CombinatorParseError) Error() string {
	line, col := e.State.Position()

	msg := fmt.Sprintf("chomp: parse error at line %d, column %d (offset %d): expected %s",
		line, col, e.State.Pos(), e.expected())

	if len(e.Labels) > 0 {
		msg += " while parsing " + strings.Join(e.Labels, " > ")
	}

	return msg
}

// snippetMaxWidth is the line length, in runes, beyond which [Snippet]
// truncates to a window around the failure column.
const snippetMaxWidth = 100

// snippetContext is the number of runes shown either side of the failure
// column once a line is truncated.
const snippetContext = 40

// Snippet renders a caret-annotated view of the source line containing the
// failure, for human-facing output such as a CLI - reached via [errors.As]
// rather than included in [CombinatorParseError.Error], so multi-line
// output only ever appears because a caller asked for it.
//
// Long lines are truncated to a window around the failure column. Tabs in
// the source line are replicated verbatim in the caret padding, so a
// terminal's own tab-stop expansion keeps both lines aligned. Alignment is
// rune-count based: double-width glyphs (e.g. CJK) may visually drift the
// caret, since that is inherently terminal-display-width dependent.
func (e CombinatorParseError) Snippet() string {
	line, col, text := e.State.lineInfo()
	runes := []rune(text)
	col0 := col - 1 // 0-based rune index into runes

	display, displayCol0 := runes, col0
	var prefix, suffix string
	if len(runes) > snippetMaxWidth {
		start, end := max(0, col0-snippetContext), min(len(runes), col0+snippetContext)
		if start > 0 {
			prefix = "… "
		}
		if end < len(runes) {
			suffix = " …"
		}
		display = runes[start:end]
		displayCol0 = col0 - start + len([]rune(prefix))
	}

	var pad strings.Builder
	for _, r := range display[:min(displayCol0, len(display))] {
		if r == '\t' {
			pad.WriteRune('\t')
		} else {
			pad.WriteRune(' ')
		}
	}

	gutter := strconv.Itoa(line)
	blank := strings.Repeat(" ", len(gutter))

	var buf strings.Builder
	fmt.Fprintf(&buf, "%s |\n", blank)
	fmt.Fprintf(&buf, "%s | %s%s%s\n", gutter, prefix, string(display), suffix)
	fmt.Fprintf(&buf, "%s | %s^ expected %s", blank, pad.String(), e.expected())

	return buf.String()
}

// LogValue implements [slog.LogValuer], emitting offset/line/column/
// expected (and context, once Labels is non-empty) as structured fields
// instead of a string to be parsed.
func (e CombinatorParseError) LogValue() slog.Value {
	line, col := e.State.Position()

	attrs := []slog.Attr{
		slog.Int("offset", e.State.Pos()),
		slog.Int("line", line),
		slog.Int("column", col),
		slog.String("expected", e.expected()),
	}

	if len(e.Labels) > 0 {
		attrs = append(attrs, slog.String("context", strings.Join(e.Labels, " > ")))
	}

	return slog.GroupValue(attrs...)
}

// ParserError defines an error that is raised when a parser
// fails to parse the input text due to a failed [Combinator].
type ParserError struct {
	// Err contains the [CombinatorParseError] that caused the parser to fail.
	Err error

	// Type of [Parser] that failed.
	Type string
}

// Error delegates to the inner error's Error(). Only the leaf
// [CombinatorParseError] carries a position; ParserError adds no prefix of
// its own so the rendered message stays the single, grep-stable line
// [CombinatorParseError.Error] produces, regardless of wrapping depth.
func (e ParserError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("chomp: %s parser error with no underlying cause", e.Type)
	}
	return e.Err.Error()
}

// Unwrap returns the inner [CombinatorParseError].
func (e ParserError) Unwrap() error {
	return e.Err
}

// LogValue delegates to the inner error's LogValue if it implements
// [slog.LogValuer], otherwise falls back to its Error() string.
func (e ParserError) LogValue() slog.Value {
	if lv, ok := e.Err.(slog.LogValuer); ok {
		return lv.LogValue()
	}
	return slog.StringValue(e.Error())
}

// RangedParserError defines an error that is raised when a ranged parser
// fails to parse the input text due to a failed [Combinator] within the
// expected execution range.
type RangedParserError struct {
	// Err contains the [CombinatorParseError] that caused the parser to fail.
	Err error

	// Range contains the execution details of the ranged parser.
	Exec RangedParserExec

	// Type of [Parser] that failed.
	Type string
}

// RangedParserExec details how a ranged [Combinator] was executed.
type RangedParserExec struct {
	// Min is the minimum number of expected executions.
	Min uint

	// Max is the maximum number of possible executions.
	Max uint

	// Count contains the number of executions.
	Count uint
}

// String returns a string representation of a [RangedParserExec].
func (e RangedParserExec) String() string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "[count: %d", e.Count)
	if e.Min > 0 {
		fmt.Fprintf(&buf, " min: %d", e.Min)
	}

	if e.Max > 0 {
		fmt.Fprintf(&buf, " max: %d", e.Max)
	}
	buf.WriteString("]")
	return buf.String()
}

// RangeExecution is a utility method for setting a [RangedParserExec].
//   - With one argument, the [RangeParserExec.Count] is set.
//   - With two arguments, the [RangeParserExec.Count] and [RangeParserExec.Min]
//     are set.
//   - With three arguments, the [RangeParserExec.Count]], [RangeParserExec.Min]
//     and [RangeParserExec.Max] are set.
//   - If four or more arguments are provided, a default [RangedParserExec] will
//     be returned.
func RangeExecution(i ...uint) RangedParserExec {
	exec := RangedParserExec{}

	switch len(i) {
	case 1:
		exec.Count = i[0]
	case 2:
		exec.Count = i[0]
		exec.Min = i[1]
	case 3:
		exec.Count = i[0]
		exec.Min = i[1]
		exec.Max = i[2]
	}

	return exec
}

// Error delegates to the inner error's Error(), for the same reason as
// [ParserError.Error]. [RangedParserError.Exec] remains available
// programmatically via [errors.As] for callers that want execution counts.
func (e RangedParserError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("chomp: %s parser error with no underlying cause", e.Type)
	}
	return e.Err.Error()
}

// Unwrap returns the inner [CombinatorParseError].
func (e RangedParserError) Unwrap() error {
	return e.Err
}

// LogValue delegates to the inner error's LogValue if it implements
// [slog.LogValuer], otherwise falls back to its Error() string.
func (e RangedParserError) LogValue() slog.Value {
	if lv, ok := e.Err.(slog.LogValuer); ok {
		return lv.LogValue()
	}
	return slog.StringValue(e.Error())
}
