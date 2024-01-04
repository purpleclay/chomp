/*
Copyright (c) 2023 Purple Clay

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
)

// Parser is a higher-ordered function capable of parsing text under a defined condition.
// Parsers can be combined to form more complex parsers. A parser differs to a [Combinator]
// due to its ability to return multiple parsed values within a slice. Upon success, a parser
// will return both the unparsed and parsed text. All parsers are strict and must parse its
// input. Any failure to do so should raise a [ParserError].
type Parser func(string) (string, []string, error)

// ParserError defines an error that is raised when a parser
// fails to parse the input text due to a failed [Combinator]
type ParserError struct {
	// Err contains the [CombinatorParseError] that caused the parser to fail.
	Err error

	// Type of [Parser] that failed.
	Type string
}

// Error returns a friendly string representation of the current error.
func (e ParserError) Error() string {
	return fmt.Sprintf("%s parser failed. %v", e.Type, e.Err)
}

// Unwrap returns the inner [CombinatorParseError].
func (e ParserError) Unwrap() error {
	return e.Err
}

// Pair will scan the input text and match each [Combinator] in turn. Both combinators
// must match. The result of each will be returned in the slice in execution order.
//
//	chomp.Pair(chomp.Tag("Hello,"), chomp.Tag(" World"))("Hello, World!")
//	// ("!", []string{"Hello,", " World"}, nil)
func Pair(c1, c2 Combinator) Parser {
	return func(s string) (string, []string, error) {
		rem, out1, err := c1(s)
		if err != nil {
			return rem, nil, ParserError{Err: err, Type: "pair"}
		}

		rem, out2, err := c2(rem)
		if err != nil {
			return rem, nil, ParserError{Err: err, Type: "pair"}
		}

		return rem, []string{out1, out2}, nil
	}
}

// SepPair will scan the input text and match each [Combinator] in turn. All combinators
// must match. The result of the separator combinator is discarded and not included
// within the returned slice.
//
//	chomp.SepPair(
//		chomp.Tag("Hello"),
//		chomp.Tag(", "),
//		chomp.Tag("World"))("Hello, World!")
//	// ("!", []string{"Hello", "World"}, nil)
func SepPair(c1, sep, c2 Combinator) Parser {
	return func(s string) (string, []string, error) {
		rem, out1, err := c1(s)
		if err != nil {
			return rem, nil, ParserError{Err: err, Type: "pair_sep"}
		}

		rem, _, err = sep(rem)
		if err != nil {
			return rem, nil, ParserError{Err: err, Type: "pair_sep"}
		}

		rem, out2, err := c2(rem)
		if err != nil {
			return rem, nil, ParserError{Err: err, Type: "pair_sep"}
		}

		return rem, []string{out1, out2}, nil
	}
}

// Repeat will scan the input text and repeat the [Combinator] the defined number
// of times. Each combinator must match, with the output of each contained in
// the returned slice.
//
//	chomp.Repeat(chomp.Parentheses(), 2)("(Hello)(World)(!)")
//	// ("(!)", []string{"(Hello)", "(World)"}, nil)
func Repeat(c Combinator, n int) Parser {
	return func(s string) (string, []string, error) {
		var ext []string
		var err error

		rem := s
		for i := 0; i < n; i++ {
			var out string
			if rem, out, err = c(rem); err != nil {
				return rem, nil, ParserError{Err: err, Type: "repeat"}
			}

			ext = append(ext, out)
		}

		return rem, ext, nil
	}
}
