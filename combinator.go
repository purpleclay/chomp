package chomp

import (
	"fmt"
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

const truncateErrAt = 50

// CombinatorParseError defines an error that is raised when a combinator
// fails to parse the input text under its expected condition.
type CombinatorParseError struct {
	// Input to the [Combinator]. This can be empty, as a combinator may
	// not require any input to parse the text.
	Input string

	// State at the point of failure, capturing the absolute position
	// within the original input alongside the unparsed suffix.
	State State

	// Type of [Combinator] that failed.
	Type string
}

// Error returns a friendly string representation of the current error.
func (e CombinatorParseError) Error() string {
	text := e.State.Rest()
	if len(text) > truncateErrAt {
		text = fmt.Sprintf("%s...(truncated)", text[:truncateErrAt])
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "(%s) combinator failed to parse text '%s'", e.Type, text)

	if e.Input != "" {
		fmt.Fprintf(&buf, " with input '%s'", e.Input)
	}

	return buf.String()
}

// ParserError defines an error that is raised when a parser
// fails to parse the input text due to a failed [Combinator].
type ParserError struct {
	// Err contains the [CombinatorParseError] that caused the parser to fail.
	Err error

	// Type of [Parser] that failed.
	Type string
}

// Error returns a friendly string representation of the current error.
func (e ParserError) Error() string {
	return fmt.Sprintf("(%s) parser failed. %v", e.Type, e.Err)
}

// Unwrap returns the inner [CombinatorParseError].
func (e ParserError) Unwrap() error {
	return e.Err
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

// Error returns a friendly string representation of the current error.
func (e RangedParserError) Error() string {
	return fmt.Sprintf("(%s) parser failed %s. %v", e.Type, e.Exec, e.Err)
}

// Unwrap returns the inner [CombinatorParseError].
func (e RangedParserError) Unwrap() error {
	return e.Err
}
