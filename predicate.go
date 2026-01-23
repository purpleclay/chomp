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

type isSpace struct{}

func (isSpace) Match(r rune) bool {
	return r == ' ' || r == '\t'
}

func (isSpace) String() string {
	return "is_space"
}

type isMultispace struct{}

func (isMultispace) Match(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func (isMultispace) String() string {
	return "is_multispace"
}

type isHexDigit struct{}

func (isHexDigit) Match(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

func (isHexDigit) String() string {
	return "is_hex_digit"
}

type isOctalDigit struct{}

func (isOctalDigit) Match(r rune) bool {
	return r >= '0' && r <= '7'
}

func (isOctalDigit) String() string {
	return "is_octal_digit"
}

type isBinaryDigit struct{}

func (isBinaryDigit) Match(r rune) bool {
	return r == '0' || r == '1'
}

func (isBinaryDigit) String() string {
	return "is_binary_digit"
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

	// IsSpace determines whether a rune is a space character. A rune is classed
	// as a space if it is either a space ' ' or a tab '\t'.
	IsSpace = isSpace{}

	// IsMultispace determines whether a rune is a whitespace character. A rune
	// is classed as whitespace if it is a space ' ', tab '\t', newline '\n',
	// or carriage return '\r'.
	IsMultispace = isMultispace{}

	// IsHexDigit determines whether a rune is a hexadecimal digit. A rune is
	// classed as a hex digit if it is between '0'-'9', 'a'-'f', or 'A'-'F'.
	IsHexDigit = isHexDigit{}

	// IsOctalDigit determines whether a rune is an octal digit. A rune is classed
	// as an octal digit if it is between '0' and '7'.
	IsOctalDigit = isOctalDigit{}

	// IsBinaryDigit determines whether a rune is a binary digit. A rune is classed
	// as a binary digit if it is either '0' or '1'.
	IsBinaryDigit = isBinaryDigit{}
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

// Alpha matches one or more ASCII or Unicode letters.
// Equivalent to While(IsLetter).
//
//	chomp.Alpha()("Hello123")
//	// ("123", "Hello", nil)
func Alpha() Combinator[string] {
	return While(IsLetter)
}

// Alpha0 matches zero or more ASCII or Unicode letters.
// Equivalent to WhileN(IsLetter, 0).
//
//	chomp.Alpha0()("123Hello")
//	// ("123Hello", "", nil)
func Alpha0() Combinator[string] {
	return WhileN(IsLetter, 0)
}

// Digit matches one or more decimal digits.
// Equivalent to While(IsDigit).
//
//	chomp.Digit()("123abc")
//	// ("abc", "123", nil)
func Digit() Combinator[string] {
	return While(IsDigit)
}

// Digit0 matches zero or more decimal digits.
// Equivalent to WhileN(IsDigit, 0).
//
//	chomp.Digit0()("abc123")
//	// ("abc123", "", nil)
func Digit0() Combinator[string] {
	return WhileN(IsDigit, 0)
}

// Alphanumeric matches one or more alphanumeric characters.
// Equivalent to While(IsAlphanumeric).
//
//	chomp.Alphanumeric()("Hello123!")
//	// ("!", "Hello123", nil)
func Alphanumeric() Combinator[string] {
	return While(IsAlphanumeric)
}

// Alphanumeric0 matches zero or more alphanumeric characters.
// Equivalent to WhileN(IsAlphanumeric, 0).
//
//	chomp.Alphanumeric0()("!Hello123")
//	// ("!Hello123", "", nil)
func Alphanumeric0() Combinator[string] {
	return WhileN(IsAlphanumeric, 0)
}

// Space matches one or more space or tab characters.
// Equivalent to While(IsSpace).
//
//	chomp.Space()("   Hello")
//	// ("Hello", "   ", nil)
func Space() Combinator[string] {
	return While(IsSpace)
}

// Space0 matches zero or more space or tab characters.
// Equivalent to WhileN(IsSpace, 0).
//
//	chomp.Space0()("Hello")
//	// ("Hello", "", nil)
func Space0() Combinator[string] {
	return WhileN(IsSpace, 0)
}

// Multispace matches one or more whitespace characters (space, tab, newline, carriage return).
// Equivalent to While(IsMultispace).
//
//	chomp.Multispace()("  \n\tHello")
//	// ("Hello", "  \n\t", nil)
func Multispace() Combinator[string] {
	return While(IsMultispace)
}

// Multispace0 matches zero or more whitespace characters (space, tab, newline, carriage return).
// Equivalent to WhileN(IsMultispace, 0).
//
//	chomp.Multispace0()("Hello")
//	// ("Hello", "", nil)
func Multispace0() Combinator[string] {
	return WhileN(IsMultispace, 0)
}

// HexDigit matches one or more hexadecimal digits (0-9, a-f, A-F).
// Equivalent to While(IsHexDigit).
//
//	chomp.HexDigit()("1a2B3c rest")
//	// (" rest", "1a2B3c", nil)
func HexDigit() Combinator[string] {
	return While(IsHexDigit)
}

// HexDigit0 matches zero or more hexadecimal digits (0-9, a-f, A-F).
// Equivalent to WhileN(IsHexDigit, 0).
//
//	chomp.HexDigit0()("xyz")
//	// ("xyz", "", nil)
func HexDigit0() Combinator[string] {
	return WhileN(IsHexDigit, 0)
}

// OctalDigit matches one or more octal digits (0-7).
// Equivalent to While(IsOctalDigit).
//
//	chomp.OctalDigit()("0127 rest")
//	// (" rest", "0127", nil)
func OctalDigit() Combinator[string] {
	return While(IsOctalDigit)
}

// OctalDigit0 matches zero or more octal digits (0-7).
// Equivalent to WhileN(IsOctalDigit, 0).
//
//	chomp.OctalDigit0()("89")
//	// ("89", "", nil)
func OctalDigit0() Combinator[string] {
	return WhileN(IsOctalDigit, 0)
}

// BinaryDigit matches one or more binary digits (0-1).
// Equivalent to While(IsBinaryDigit).
//
//	chomp.BinaryDigit()("1010 rest")
//	// (" rest", "1010", nil)
func BinaryDigit() Combinator[string] {
	return While(IsBinaryDigit)
}

// BinaryDigit0 matches zero or more binary digits (0-1).
// Equivalent to WhileN(IsBinaryDigit, 0).
//
//	chomp.BinaryDigit0()("234")
//	// ("234", "", nil)
func BinaryDigit0() Combinator[string] {
	return WhileN(IsBinaryDigit, 0)
}

// Newline matches a single newline character '\n'.
//
//	chomp.Newline()("\nHello")
//	// ("Hello", "\n", nil)
func Newline() Combinator[string] {
	return func(s string) (string, string, error) {
		if len(s) > 0 && s[0] == '\n' {
			return s[1:], "\n", nil
		}
		return s, "", CombinatorParseError{Text: s, Type: "newline"}
	}
}

// Tab matches a single tab character '\t'.
//
//	chomp.Tab()("\tHello")
//	// ("Hello", "\t", nil)
func Tab() Combinator[string] {
	return func(s string) (string, string, error) {
		if len(s) > 0 && s[0] == '\t' {
			return s[1:], "\t", nil
		}
		return s, "", CombinatorParseError{Text: s, Type: "tab"}
	}
}

// NotLineEnding matches any characters until a line ending ('\n' or '\r').
// Requires at least one character to be matched.
//
//	chomp.NotLineEnding()("Hello, World!\nNext line")
//	// ("\nNext line", "Hello, World!", nil)
func NotLineEnding() Combinator[string] {
	return WhileNot(IsLineEnding)
}
