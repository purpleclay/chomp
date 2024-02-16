/*
Copyright (c) 2023 - 2024 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
// a combinator will return both the unparsed and parsed text. All combinators are
// strict and must parse its input. Any failure to do so should raise a [CombinatorParseError].
type Combinator[T Result] func(string) (string, T, error)

const truncateErrAt = 50

// CombinatorParseError defines an error that is raised when a combinator
// fails to parse the input text under its expected condition.
type CombinatorParseError struct {
	// Input to the [Combinator]. This can be empty, as a combinator may
	// not require any input to parse the text.
	Input string

	// Text that was being parsed by the [Combinator]. This will be truncated
	// in the error message.
	Text string

	// Type of [Combinator] that failed.
	Type string
}

// Error returns a friendly string representation of the current error.
func (e CombinatorParseError) Error() string {
	text := e.Text
	if len(text) > truncateErrAt {
		text = fmt.Sprintf("%s...(truncated)", text[:truncateErrAt])
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("(%s) combinator failed to parse text '%s'", e.Type, text))

	if e.Input != "" {
		buf.WriteString(fmt.Sprintf(" with input '%s'", e.Input))
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

// RangedParserExec details how a ranged [Combinator] was exeucted.
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
	buf.WriteString(fmt.Sprintf("[count: %d", e.Count))
	if e.Min > 0 {
		buf.WriteString(fmt.Sprintf(" min: %d", e.Min))
	}

	if e.Max > 0 {
		buf.WriteString(fmt.Sprintf(" max: %d", e.Max))
	}
	buf.WriteString("]")
	return buf.String()
}

// RangeExecution ...
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
