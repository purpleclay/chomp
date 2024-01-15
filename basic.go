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

// Tag must match a series of characters at the beginning of the input text,
// in the exact order and case provided.
//
//	chomp.Tag("Hello")("Hello, World!")
//	// (", World!", "Hello", nil)
func Tag(str string) Combinator[string] {
	return func(s string) (string, string, error) {
		if strings.HasPrefix(s, str) {
			return s[len(str):], str, nil
		}

		return s, "", CombinatorParseError{Input: str, Text: s, Type: "tag"}
	}
}

// Any must match at least one character at the beginning of the input text,
// from the provided sequence. Parsing immediately stops upon the first
// unmatched character.
//
//	chomp.Any("eH")("Hello, World!")
//	// ("llo, World!", "He", nil)
func Any(str string) Combinator[string] {
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
			return s, "", CombinatorParseError{Input: str, Text: s, Type: "any"}
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
func Not(str string) Combinator[string] {
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
			return s, "", CombinatorParseError{Input: str, Text: s, Type: "not"}
		}

		return s[pos:], s[:pos], nil
	}
}

// Crlf must match either a CR or CRLF line ending.
//
//	chomp.Crlf()("\r\nHello")
//	// ("Hello", "\r\n", nil)
func Crlf() Combinator[string] {
	return func(s string) (string, string, error) {
		idx := strings.Index(s, "\n")
		if idx == 0 || (idx == 1 && s[0] == '\r') {
			return s[idx+1:], s[:idx+1], nil
		}

		return "", "", CombinatorParseError{Text: s, Type: "crlf"}
	}
}

// OneOf must match a single character at the beginning of the text from the
// provided sequence.
//
//	chomp.OneOf("!,eH")("Hello, World!")
//	// ("ello, World!", "H", nil)
func OneOf(str string) Combinator[string] {
	return func(s string) (string, string, error) {
		if txt := []rune(s); len(txt) > 0 {
			for _, strc := range str {
				if txt[0] == strc {
					pos := len(string(txt[0]))
					return s[pos:], s[:pos], nil
				}
			}
		}

		return "", "", CombinatorParseError{Input: str, Text: s, Type: "one_of"}
	}
}

// NoneOf must not match a single character at the beginning of the text from
// the provided sequence.
//
//	chomp.NoneOf("loWrd!e")("Hello, World!")
//	// ("ello, World!", "H", nil)
func NoneOf(str string) Combinator[string] {
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

		return "", "", CombinatorParseError{Input: str, Text: s, Type: "none_of"}
	}
}

// Until will scan the input text for the first occurrence of the provided series
// of characters. Everything until that point in the text will be matched.
//
//	chomp.Until("World")("Hello, World!")
//	// ("World!", "Hello, ", nil)
func Until(str string) Combinator[string] {
	return func(s string) (string, string, error) {
		if idx := strings.Index(s, str); idx != -1 {
			return s[idx:], s[:idx], nil
		}

		return "", "", CombinatorParseError{Input: str, Text: s, Type: "until"}
	}
}

// Opt allows a combinator to be optional. Any error returned by the underlying
// combinator will be swallowed. The parsed text will not be modified if the
// underlying combinator did not run.
//
//	chomp.Opt(chomp.Tag("Hey"))("Hello, World!")
//	// ("Hello, World!", "", nil)
func Opt[T Result](c Combinator[T]) Combinator[T] {
	return func(s string) (string, T, error) {
		rem, out, _ := c(s)
		return rem, out, nil
	}
}

// S wraps the result of the inner combinator within a string slice.
// Combinators of differing return types can be successfully chained
// together while using this conversion combinator.
//
//	chomp.S(chomp.Until(","))("Hello, World!")
//	// (", World!", []string{"Hello"}, nil)
func S(c Combinator[string]) Combinator[[]string] {
	return func(s string) (string, []string, error) {
		rem, ext, err := c(s)
		if err != nil {
			return rem, nil, err
		}

		return rem, []string{ext}, err
	}
}

// I extracts and returns a single string from the result of the inner combinator.
// Combinators of differing return types can be successfully chained together while
// using this conversion combinator.
//
//	chomp.I(chomp.SepPair(
//		chomp.Tag("Hello"),
//		chomp.Tag(", "),
//		chomp.Tag("World")), 1)("Hello, World!")
//	// ("!", "World", nil)
func I(c Combinator[[]string], i int) Combinator[string] {
	return func(s string) (string, string, error) {
		rem, ext, err := c(s)
		if err != nil {
			return rem, "", err
		}

		if i < 0 || i >= len(ext) {
			return rem, "", ParserError{
				Err:  fmt.Errorf("index %d is out of bounds within string slice of %d elements", i, len(ext)),
				Type: "i",
			}
		}

		return rem, ext[i], nil
	}
}
