package chomp

import (
	"fmt"
	"unicode"
)

// Predicate defines an expression that will return either true or false
type Predicate interface {
	// Match a rune against a defined expression, returning true
	// if the condition is met
	Match(r rune) bool

	// Returns the name of the predicate for error handling
	fmt.Stringer
}

type isDigit struct{}

func (isDigit) Match(r rune) bool {
	return unicode.IsDigit(r)
}

func (isDigit) String() string {
	return "is_digit"
}

type isLetter struct{}

func (isLetter) Match(r rune) bool {
	return unicode.IsLetter(r)
}

func (isLetter) String() string {
	return "is_letter"
}

type isAlphanumeric struct{}

func (isAlphanumeric) Match(r rune) bool {
	return unicode.IsDigit(r) || unicode.IsLetter(r)
}

func (isAlphanumeric) String() string {
	return "is_alphanumeric"
}

type isLineEnding struct{}

func (isLineEnding) Match(r rune) bool {
	return r == '\n' || r == '\r'
}

func (isLineEnding) String() string {
	return "is_line_ending"
}

var (
	// IsDigit determines whether a rune is a decimal digit. A rune is classed
	// as a digit if it is between the ASCII range of '0' or '9', or if it belongs
	// within the Unicode [Nd] category.
	//
	// [Nd]: https://www.fileformat.info/info/unicode/category/Nd/list.htm
	IsDigit = isDigit{}

	// IsLetter determines if a rune is a letter. A rune is classed as a letter
	// if it is between the ASCII range of 'a' and 'z' (including its uppercase
	// equivalents), or it belongs within any of the Unicode letter categories:
	// [Lu] [LI] [Lt] [Lm] [Lo].
	//
	// [Lu]: https://www.fileformat.info/info/unicode/category/Lu/list.htm
	// [LI]: https://www.fileformat.info/info/unicode/category/Ll/list.htm
	// [Lt]: https://www.fileformat.info/info/unicode/category/Lt/list.htm
	// [Lm]: https://www.fileformat.info/info/unicode/category/Lm/list.htm
	// [Lo]: https://www.fileformat.info/info/unicode/category/Lo/list.htm
	IsLetter = isLetter{}

	// IsAlphanumeric determines whether a rune is a decimal digit or a letter.
	// This convenience method wraps the existing [IsDigit] and [IsLetter]
	// predicates.
	IsAlphanumeric = isAlphanumeric{}

	// IsLineEnding determines whether a rune is one of the following ASCII
	// line ending characters '\r' or '\n'.
	IsLineEnding = isLineEnding{}
)

// While will scan the input text, testing each character against the provided
// [Predicate]. The [Predicate] must match at least one character.
//
//	chomp.While(chomp.IsLetter)("Hello, World!")
//	// (", World!", "Hello", nil)
func While(p Predicate) Combinator[string] {
	return WhileN(p, 1)
}

// WhileN will scan the input text, testing each character against the provided
// [Predicate]. The [Predicate] must match at least n characters. If n is zero,
// this becomes an optional combinator.
//
//	chomp.WhileN(chomp.IsLetter, 1)("Hello, World!")
//	// (", World!", "Hello", nil)
//
//	chomp.WhileN(chomp.IsDigit, 0)("Hello, World!")
//	// ("Hello, World!", "", nil)
func WhileN(p Predicate, n uint) Combinator[string] {
	return func(s string) (string, string, error) {
		pos := 0
		for _, c := range s {
			if !p.Match(c) {
				break
			}
			pos += len(string(c))
		}

		if uint(pos) < n {
			return s, "", RangedParserError{
				Err:  CombinatorParseError{Text: s, Type: p.String()},
				Exec: RangeExecution(uint(pos), n),
				Type: "while_n",
			}
		}

		return s[pos:], s[:pos], nil
	}
}

// WhileNM will scan the input text, testing each character against the provided
// [Predicate]. The [Predicate] must match a minimum of n and upto a maximum
// of m characters. If n is zero, this becomes an optional combinator.
//
//	chomp.WhileNM(chomp.IsLetter, 1, 8)("Hello, World!")
//	// (", World!", "Hello", nil)
func WhileNM(p Predicate, n, m uint) Combinator[string] {
	return func(s string) (string, string, error) {
		pos := 0
		for _, c := range s {
			if !p.Match(c) {
				break
			}
			pos += len(string(c))
		}

		if uint(pos) < n || uint(pos) > m {
			return s, "", RangedParserError{
				Err:  CombinatorParseError{Text: s, Type: p.String()},
				Exec: RangeExecution(uint(pos), n, m),
				Type: "while_n_m",
			}
		}

		return s[pos:], s[:pos], nil
	}
}

// WhileNot will scan the input text, testing each character against the provided
// [Predicate]. The [Predicate] must not match at least one character. It has
// the inverse behavior of [While].
//
//	chomp.WhileNot(chomp.IsDigit)("Hello, World!")
//	// ("", "Hello, World!", nil)
func WhileNot(p Predicate) Combinator[string] {
	return WhileNotN(p, 1)
}

// WhileNotN will scan the input text, testing each character against the provided
// [Predicate]. The [Predicate] must not match at least n characters. If n is
// zero, this becomes an optional combinator. It has the inverse behavior of [WhileN].
//
//	chomp.WhileNotN(chomp.IsDigit, 1)("Hello, World!")
//	// ("", "Hello, World!", nil)
//
//	chomp.WhileNotN(chomp.IsLetter, 0)("Hello, World!")
//	// ("Hello, World!", "", nil)
func WhileNotN(p Predicate, n uint) Combinator[string] {
	return func(s string) (string, string, error) {
		pos := 0
		for _, c := range s {
			if p.Match(c) {
				break
			}
			pos += len(string(c))
		}

		if uint(pos) < n {
			return s, "", RangedParserError{
				Err:  CombinatorParseError{Text: s, Type: p.String()},
				Exec: RangeExecution(uint(pos), n),
				Type: "while_not_n",
			}
		}

		return s[pos:], s[:pos], nil
	}
}

// WhileNotNM will scan the input text, testing each character against the provided
// [Predicate]. The [Predicate] must not match a minimum of n and upto a maximum of
// m characters. If n is zero, this becomes an optional combinator. It has the
// inverse behavior of [WhileNM].
//
//	chomp.WhileNotNM(chomp.IsLetter, 1, 8)("20240709 was a great day")
//	// (" was a great day", "20240709", nil)
func WhileNotNM(p Predicate, n, m uint) Combinator[string] {
	return func(s string) (string, string, error) {
		pos := 0
		for _, c := range s {
			if p.Match(c) {
				break
			}
			pos += len(string(c))
		}

		if uint(pos) < n || uint(pos) > m {
			return s, "", RangedParserError{
				Err:  CombinatorParseError{Text: s, Type: p.String()},
				Exec: RangeExecution(uint(pos), n, m),
				Type: "while_not_n_m",
			}
		}

		return s[pos:], s[:pos], nil
	}
}
