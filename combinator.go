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
	"strings"
)

// Combinator is a higher-order function capable of parsing text under a defined
// condition. Combinators can be combined to form more complex parsers. Upon success,
// a combinator will return both the unparsed and parsed text. All combinators are
// strict and must parse its input. Any failure to do so will raise a [CombinatorParseError].
type Combinator func(string) (string, string, error)

// CombinatorParseError defines an error that is raised when a combinator
// fails to parse the input text under its expected condition.
type CombinatorParseError struct {
	Input string
	Type  string
}

// Error returns a friendly string representation of the current error.
func (e CombinatorParseError) Error() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s combinator failed to parse text", e.Type))

	if e.Input != "" {
		buf.WriteString(fmt.Sprintf(" using input '%s'", e.Input))
	}

	return buf.String()
}

// Tag must match a series of characters at the beginning of the input text,
// in the exact order and case provided.
//
//	chomp.Tag("Hello")("Hello, World!")
//	// (", World!", "Hello", nil)
func Tag(str string) Combinator {
	return func(s string) (string, string, error) {
		if strings.HasPrefix(s, str) {
			return s[len(str):], str, nil
		}

		return s, "", CombinatorParseError{Input: str, Type: "tag"}
	}
}

// Any must match at least one character at the beginning of the input text,
// from the provided sequence. Parsing immediately stops upon the first
// unmatched character.
//
//	chomp.Any("eH")("Hello, World!")
//	// ("llo, World!", "He", nil)
func Any(str string) Combinator {
	return func(s string) (string, string, error) {
		pos := 0

	match:
		for _, sc := range s {
			for _, strc := range str {
				if sc == strc {
					pos = pos + len(string(strc))
					continue match
				}
			}

			break match
		}

		if pos == 0 {
			return s, "", CombinatorParseError{Input: str, Type: "any"}
		}

		return s[pos:], s[:pos], nil
	}
}

// Not must not match at least one character at the beginning of the input
// text from the provided sequence. Parsing immediately stops upon the
// first matched character.
//
//	chomp.Not("ol")("Hello, World!")
//	// ("llo, World!", "He", nil)
func Not(str string) Combinator {
	return func(s string) (string, string, error) {
		pos := 0

	match:
		for _, sc := range s {
			for _, strc := range str {
				if sc == strc {
					break match
				}
			}

			pos = pos + len(string(sc))
		}

		if pos == 0 {
			return s, "", CombinatorParseError{Input: str, Type: "not"}
		}

		return s[pos:], s[:pos], nil
	}
}

// Crlf must match either a CR or CRLF line ending.
//
//	chomp.Crlf()("\r\nHello")
//	// ("Hello", "\r\n", nil)
func Crlf() Combinator {
	return func(s string) (string, string, error) {
		idx := strings.Index(s, "\n")
		if idx == 0 || (idx == 1 && s[0] == '\r') {
			return s[idx+1:], s[:idx+1], nil
		}

		return "", "", CombinatorParseError{Type: "crlf"}
	}
}

// OneOf must match a single character at the beginning of the text from the
// provided sequence.
//
//	chomp.OneOf("!,eH")("Hello, World!")
//	// ("ello, World!", "H", nil)
func OneOf(str string) Combinator {
	return func(s string) (string, string, error) {
		if txt := []rune(s); len(txt) > 0 {
			for _, strc := range str {
				if txt[0] == strc {
					pos := len(string(txt[0]))
					return s[pos:], s[:pos], nil
				}
			}
		}

		return "", "", CombinatorParseError{Input: str, Type: "oneof"}
	}
}

// NoneOf must not match a single character at the beginning of the text from
// the provided sequence.
//
//	chomp.NoneOf("loWrd!e")("Hello, World!")
//	// ("ello, World!", "H", nil)
func NoneOf(str string) Combinator {
	return func(s string) (string, string, error) {
		if txt := []rune(s); len(txt) > 0 {
			for _, strc := range str {
				if txt[0] == strc {
					break
				}
			}

			pos := len(string(txt[0]))
			return s[pos:], s[:pos], nil
		}

		return "", "", CombinatorParseError{Input: str, Type: "noneof"}
	}
}

// Until will scan the input text for the first occurrence of the provided series
// of characters. Everything until that point in the text will be matched.
//
//	chomp.Until("World")("Hello, World!")
//	// ("World!", "Hello, ", nil)
func Until(str string) Combinator {
	return func(s string) (string, string, error) {
		if idx := strings.Index(s, str); idx != -1 {
			return s[idx:], s[:idx], nil
		}

		return "", "", CombinatorParseError{Input: str, Type: "until"}
	}
}
