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

import "unicode"

// Predicate defines an expression that will return either true or false
type Predicate func(rune) bool

var (
	// IsDigit determines whether a rune is a decimal digit. A rune is classed
	// as a digit if it is between the ASCII range of '0' or '9', or it belongs
	// within the Unicode [Nd] category.
	//
	// [Nd]: https://www.fileformat.info/info/unicode/category/Nd/list.htm
	IsDigit = unicode.IsDigit

	// IsLetter determines whether a rune is a letter. A rune is classed as a
	// letter if it is between the ASCII range of 'a' and 'z' (including its
	// uppercase equivalents), or it belongs within any of the Unicode letter
	// categories: [Lu] [LI] [Lt] [Lm] [Lo]
	//
	// [Lu]: https://www.fileformat.info/info/unicode/category/Lu/list.htm
	// [LI]: https://www.fileformat.info/info/unicode/category/Ll/list.htm
	// [Lt]: https://www.fileformat.info/info/unicode/category/Lt/list.htm
	// [Lm]: https://www.fileformat.info/info/unicode/category/Lm/list.htm
	// [Lo]: https://www.fileformat.info/info/unicode/category/Lo/list.htm
	IsLetter = unicode.IsLetter
)

// While will scan the input text, testing each character against the provided
// [Predicate]. Everything until the predicate returns false will be matched.
//
//	chomp.While(chomp.IsLetter)("Hello, World!")
//	// (", World!", "Hello", nil)
func While(p Predicate) Combinator[string] {
	return func(s string) (string, string, error) {
		pos := 0
		for _, c := range s {
			if !p(c) {
				break
			}
			pos += len(string(c))
		}

		if pos == 0 {
			return "", "", CombinatorParseError{Text: s, Type: "while"}
		}

		return s[pos:], s[:pos], nil
	}
}

// WhileNot will scan the input text, testing each character against the provided
// [Predicate]. Everything until the predicate returns true will be matched. This
// is the inverse of [While].
//
//	chomp.WhileNot(chomp.IsDigit)("Hello, World!")
//	// ("", "Hello, World!", nil)
func WhileNot(p Predicate) Combinator[string] {
	return func(s string) (string, string, error) {
		pos := 0
		for _, c := range s {
			if p(c) {
				break
			}
			pos += len(string(c))
		}

		if pos == 0 {
			return "", "", CombinatorParseError{Text: s, Type: "while_not"}
		}

		return s[pos:], s[:pos], nil
	}
}
