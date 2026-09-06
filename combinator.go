package chomp

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

// Combinator is a higher-order function capable of parsing text under a defined
// condition. Combinators can be combined to form more complex parsers. Upon success,
// a combinator will return both the unparsed and parsed [State]. All combinators are
// strict and must parse its input. Any failure to do so should raise a [CombinatorParseError].
//
// Raw invocation c(state) is the composition surface used by combinators calling
// combinators (and custom drivers); its errors are unfinalised and may still
// reference the live input. Use [Combinator.Run] as the string-in/string-out entry
// point, which finalises any error via [Finalize] before returning it.
type Combinator[T any] func(State) (State, T, error)

// Tuple2 holds the results of two combinators matched by [Pair], [SepPair],
// or [Consumed].
type Tuple2[A, B any] struct {
	First  A
	Second B
}

// run is shared by [Combinator.Run].
func run[T any](c func(State) (State, T, error), input string) (string, T, error) { //nolint:ireturn // T is caller-supplied, not an open interface.
	rem, ext, err := c(NewState(input))
	if err != nil {
		return rem.Rest(), ext, Finalize(err)
	}

	return rem.Rest(), ext, nil
}

// Run applies c against input, the sole string-in/string-out entry point.
// Unlike raw invocation c(state), any returned error is passed through
// [Finalize] first.
func (c Combinator[T]) Run(input string) (string, T, error) { //nolint:ireturn // T is caller-supplied, not an open interface.
	return run[T](c, input)
}

// Finalize prepares an error returned by a raw combinator invocation for use
// outside a parse. It is the boundary [Combinator.Run] applies errors to
// before returning them; custom drivers built on raw c(state) invocation
// should apply it at their own boundary.
//
// For every [CombinatorParseError] reached through err's chain (walking
// down through [ParserError], [RangedParserError], [CutError], and
// [AlternativesError], the only wrapper types this package produces), it
// captures the failure's position once and clones Snippet's bounded
// display window into its own backing array, then drops the reference to
// the original input. An escaped error therefore holds O(1) memory
// regardless of input size, rather than pinning the entire input for as
// long as the error lives. Errors outside this closed set are returned
// unchanged. Applying Finalize more than once is a no-op.
func Finalize(err error) error {
	switch e := err.(type) {
	case CombinatorParseError:
		return e.finalize()
	case ParserError:
		e.Err = Finalize(e.Err)
		return e
	case RangedParserError:
		e.Err = Finalize(e.Err)
		return e
	case CutError:
		e.Err = Finalize(e.Err)
		return e
	case AlternativesError:
		for i := range e.Errs {
			e.Errs[i] = Finalize(e.Errs[i])
		}
		return e
	default:
		return err
	}
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
	// within the original input alongside the unparsed suffix. Only
	// meaningful before the error is finalised (see [Finalize]): a
	// finalised error - any error returned by [Combinator.Run] - clears
	// State, since retaining it would keep the entire original input
	// reachable for as long as the error lives. Use
	// [CombinatorParseError.Offset] and [CombinatorParseError.Position]
	// instead, which work correctly either way.
	State State

	// Labels records the chain of grammar-level names (see the future
	// Label combinator) being parsed when the failure occurred, outermost
	// first. Always empty until a caller wraps a combinator with Label.
	Labels []string

	// Cause, when non-nil, is the underlying error a semantic check (see
	// [MapRes]) failed with. Reachable via errors.Is/errors.As through
	// [CombinatorParseError.Unwrap], so a caller can branch on a stdlib
	// sentinel with no chomp-specific vocabulary.
	Cause error

	// kind identifies which combinator produced this error; used
	// internally to derive a fallback for Expected. Not exported: it is
	// deliberately not part of the public error-classification surface.
	kind string

	// finalized caches a bounded, self-contained snapshot of the failure
	// position once Finalize runs, so Error/Snippet/LogValue no longer
	// need State (and the input it retains) to keep working.
	finalized *errSnapshot
}

// Unwrap returns Cause. Returns nil when unset, the same as every other
// leaf failure that doesn't wrap an underlying cause.
func (e CombinatorParseError) Unwrap() error {
	return e.Cause
}

// errSnapshot is the O(1), input-independent view [CombinatorParseError]
// renders from: either derived live from State on every call, or computed
// once and cached by [CombinatorParseError.finalize].
type errSnapshot struct {
	offset int
	line   int
	col    int

	// prefix/content/suffix are Snippet's line display, already windowed to
	// snippetMaxWidth; dispCol is the caret column within content (1-based).
	prefix, content, suffix string
	dispCol                 int
}

// snapshot resolves the current, O(1) view of where this error occurred:
// the cached finalized snapshot if one exists, otherwise computed live
// from State (cheap - the input is alive anyway during an in-flight parse).
func (e CombinatorParseError) snapshot() errSnapshot {
	if e.finalized != nil {
		return *e.finalized
	}

	line, col, text := e.State.lineInfo()
	prefix, content, suffix, dispCol := windowLine(text, col)

	return errSnapshot{
		offset:  e.State.Pos(),
		line:    line,
		col:     col,
		prefix:  prefix,
		content: content,
		suffix:  suffix,
		dispCol: dispCol,
	}
}

// windowLine bounds text to a window around col (1-based) once it exceeds
// snippetMaxWidth runes, inserting "… "/" …" ellipsis markers at
// truncation points. Returns the content to display (excluding markers),
// the markers themselves, and a caret column re-based to account for the
// full rendered line (markers + content).
func windowLine(text string, col int) (prefix, content, suffix string, dispCol int) {
	runes := []rune(text)
	col0 := col - 1
	if len(runes) <= snippetMaxWidth {
		return "", text, "", col
	}

	start, end := max(0, col0-snippetContext), min(len(runes), col0+snippetContext)
	if start > 0 {
		prefix = "… "
	}
	if end < len(runes) {
		suffix = " …"
	}
	content = string(runes[start:end])
	dispCol0 := col0 - start + len([]rune(prefix))

	return prefix, content, suffix, dispCol0 + 1
}

// finalize returns a copy of e that no longer references the live input:
// its position and Snippet's bounded display window are captured once and
// cloned into their own backing array, and State is cleared. Applying
// finalize more than once is a no-op.
func (e CombinatorParseError) finalize() CombinatorParseError {
	if e.finalized != nil {
		return e
	}

	snap := e.snapshot()
	snap.content = strings.Clone(snap.content)
	e.finalized = &snap
	e.State = State{}

	return e
}

// typeFallback maps a kind with no natural literal to quote, and no "is_"
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
	"many_n":            "further progress",
	"many_till":         "further progress",
	"many_till_0":       "further progress",
	"fold_many":         "further progress",
	"many_count":        "further progress",
}

// expected resolves the human phrase completing "expected ___": Expected
// if set, then Cause's own message, otherwise a description derived from
// kind.
func (e CombinatorParseError) expected() string {
	if e.Expected != "" {
		return e.Expected
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	if after, ok := strings.CutPrefix(e.kind, "is_"); ok {
		return strings.ReplaceAll(after, "_", " ")
	}
	if phrase, ok := typeFallback[e.kind]; ok {
		return phrase
	}
	return strings.ReplaceAll(e.kind, "_", " ")
}

// labelSuffix returns " while parsing A > B" when Labels is non-empty
// (outermost first), otherwise an empty string.
func (e CombinatorParseError) labelSuffix() string {
	if len(e.Labels) == 0 {
		return ""
	}
	return " while parsing " + strings.Join(e.Labels, " > ")
}

// Offset returns the byte offset into the original input where parsing
// failed. Unlike reading [CombinatorParseError.State] directly, this works
// correctly whether or not the error has been finalised (see [Finalize]).
func (e CombinatorParseError) Offset() int {
	return e.snapshot().offset
}

// Position returns the 1-based line and rune-aware column where parsing
// failed. Unlike reading [CombinatorParseError.State] directly, this works
// correctly whether or not the error has been finalised (see [Finalize]).
func (e CombinatorParseError) Position() (line, col int) {
	s := e.snapshot()
	return s.line, s.col
}

// Error returns a single-line, grep-stable string representation of the
// current error: "chomp: parse error at line %d, column %d (offset %d):
// expected %s". It never contains a newline. For a human-facing,
// caret-annotated view of the failure, see [CombinatorParseError.Snippet].
func (e CombinatorParseError) Error() string {
	s := e.snapshot()

	return fmt.Sprintf("chomp: parse error at line %d, column %d (offset %d): expected %s%s",
		s.line, s.col, s.offset, e.expected(), e.labelSuffix())
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
	s := e.snapshot()
	content := []rune(s.content)

	var pad strings.Builder
	for _, r := range content[:min(s.dispCol-1, len(content))] {
		if r == '\t' {
			pad.WriteRune('\t')
		} else {
			pad.WriteRune(' ')
		}
	}

	gutter := strconv.Itoa(s.line)
	blank := strings.Repeat(" ", len(gutter))

	var buf strings.Builder
	fmt.Fprintf(&buf, "%s |\n", blank)
	fmt.Fprintf(&buf, "%s | %s%s%s\n", gutter, s.prefix, s.content, s.suffix)
	fmt.Fprintf(&buf, "%s | %s^ expected %s%s", blank, pad.String(), e.expected(), e.labelSuffix())

	return buf.String()
}

// LogValue implements [slog.LogValuer], emitting offset/line/column/
// expected (and context, once Labels is non-empty) as structured fields
// instead of a string to be parsed.
func (e CombinatorParseError) LogValue() slog.Value {
	s := e.snapshot()

	attrs := []slog.Attr{
		slog.Int("offset", s.offset),
		slog.Int("line", s.line),
		slog.Int("column", s.col),
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

	// kind identifies which combinator produced this error. Not exported:
	// deliberately not part of the public error-classification surface.
	kind string
}

// Error delegates to the inner error's Error(). Only the leaf
// [CombinatorParseError] carries a position; ParserError adds no prefix of
// its own so the rendered message stays the single, grep-stable line
// [CombinatorParseError.Error] produces, regardless of wrapping depth.
func (e ParserError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("chomp: %s parser error with no underlying cause", e.kind)
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

	// kind identifies which combinator produced this error. Not exported:
	// deliberately not part of the public error-classification surface.
	kind string
}

// RangedParserExec details how a ranged [Combinator] was executed.
type RangedParserExec struct {
	// Min is the minimum number of expected executions.
	Min int

	// Max is the maximum number of possible executions.
	Max int

	// Count contains the number of executions.
	Count int
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

// Error delegates to the inner error's Error(), for the same reason as
// [ParserError.Error]. [RangedParserError.Exec] remains available
// programmatically via [errors.As] for callers that want execution counts.
func (e RangedParserError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("chomp: %s parser error with no underlying cause", e.kind)
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

// AlternativesError is returned by [First] when every alternative fails.
// Unwrap exposes each alternative's own error, in the order attempted, via
// the multi-error form (the same shape [errors.Join] uses) so
// errors.Is/errors.As can inspect any of them, not just the first.
type AlternativesError struct {
	// Errs contains each attempted alternative's own error, in the order
	// [First] tried them. Never includes an alternative past a [CutError]:
	// a cut exits immediately rather than continuing to the next
	// alternative.
	Errs []error
}

// Error renders the first attempted alternative's own message. Unlike
// [errors.Join], this stays single-line: joining every alternative's
// message would break Error's single-line guarantee.
func (e AlternativesError) Error() string {
	if len(e.Errs) == 0 {
		return "chomp: no alternatives matched"
	}
	return e.Errs[0].Error()
}

// Unwrap returns every attempted alternative's error, letting
// errors.Is/errors.As reach any of them.
func (e AlternativesError) Unwrap() []error {
	return e.Errs
}

// LogValue delegates to the first alternative's LogValue if it implements
// [slog.LogValuer], otherwise falls back to its Error() string.
func (e AlternativesError) LogValue() slog.Value {
	if len(e.Errs) == 0 {
		return slog.StringValue(e.Error())
	}
	if lv, ok := e.Errs[0].(slog.LogValuer); ok {
		return lv.LogValue()
	}
	return slog.StringValue(e.Error())
}
