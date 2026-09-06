package chomp

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// Predicate defines an expression that will return either true or false
type Predicate interface {
	// Match a rune against a defined expression, returning true
	// if the condition is met
	Match(r rune) bool

	// Returns the name of the predicate for error handling
	fmt.Stringer
}

// namedPredicate adapts a plain predicate function into a [Predicate],
// pairing it with a name used for error messages.
type namedPredicate struct {
	name string
	fn   func(rune) bool
}

func (p namedPredicate) Match(r rune) bool { return p.fn(r) }
func (p namedPredicate) String() string    { return p.name }

// Named builds a [Predicate] from a plain function and a name, without
// requiring a hand-written struct implementing Match and String.
//
//	chomp.While(chomp.Named("vowel", func(r rune) bool {
//		return strings.ContainsRune("aeiouAEIOU", r)
//	})).Run("hello")
func Named(name string, f func(rune) bool) Predicate { //nolint:ireturn // Predicate is this package's own exported interface, the documented return type for a predicate constructor.
	return namedPredicate{name: name, fn: f}
}

var (
	// IsDigit determines whether a rune is a decimal digit. A rune is classed
	// as a digit if it is between the ASCII range of '0' or '9', or if it belongs
	// within the Unicode [Nd] category.
	//
	// [Nd]: https://www.fileformat.info/info/unicode/category/Nd/list.htm
	IsDigit = Named("is_digit", unicode.IsDigit)

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
	IsLetter = Named("is_letter", unicode.IsLetter)

	// IsAlphanumeric determines whether a rune is a decimal digit or a letter.
	// This convenience method wraps the existing [IsDigit] and [IsLetter]
	// predicates.
	IsAlphanumeric = Named("is_alphanumeric", func(r rune) bool {
		return unicode.IsDigit(r) || unicode.IsLetter(r)
	})

	// IsLineEnding determines whether a rune is one of the following ASCII
	// line ending characters '\r' or '\n'.
	IsLineEnding = Named("is_line_ending", func(r rune) bool {
		return r == '\n' || r == '\r'
	})

	// IsSpace determines whether a rune is a space character. A rune is classed
	// as a space if it is either a space ' ' or a tab '\t'.
	IsSpace = Named("is_space", func(r rune) bool {
		return r == ' ' || r == '\t'
	})

	// IsMultispace determines whether a rune is a whitespace character. A rune
	// is classed as whitespace if it is a space ' ', tab '\t', newline '\n',
	// or carriage return '\r'.
	IsMultispace = Named("is_multispace", func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})

	// IsHexDigit determines whether a rune is a hexadecimal digit. A rune is
	// classed as a hex digit if it is between '0'-'9', 'a'-'f', or 'A'-'F'.
	IsHexDigit = Named("is_hex_digit", func(r rune) bool {
		return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
	})

	// IsOctalDigit determines whether a rune is an octal digit. A rune is classed
	// as an octal digit if it is between '0' and '7'.
	IsOctalDigit = Named("is_octal_digit", func(r rune) bool {
		return r >= '0' && r <= '7'
	})

	// IsBinaryDigit determines whether a rune is a binary digit. A rune is classed
	// as a binary digit if it is either '0' or '1'.
	IsBinaryDigit = Named("is_binary_digit", func(r rune) bool {
		return r == '0' || r == '1'
	})
)

// While will scan the input text, testing each character against the provided
// [Predicate]. The [Predicate] must match at least one character.
//
//	chomp.While(chomp.IsLetter).Run("Hello, World!")
//	// (", World!", "Hello", nil)
func While(p Predicate) Combinator[string] {
	return WhileN(p, 1)
}

// WhileN will scan the input text, testing each character against the provided
// [Predicate]. The [Predicate] must match at least n characters. If n is zero,
// this becomes an optional combinator.
//
//	chomp.WhileN(chomp.IsLetter, 1).Run("Hello, World!")
//	// (", World!", "Hello", nil)
//
//	chomp.WhileN(chomp.IsDigit, 0).Run("Hello, World!")
//	// ("Hello, World!", "", nil)
func WhileN(p Predicate, n int) Combinator[string] {
	return func(s State) (State, string, error) {
		if n < 0 {
			return s, "", ParserError{Err: fmt.Errorf("chomp: count must be non-negative, got %d", n), kind: "while_n"}
		}
		rest := s.Rest()
		pos, runes := 0, 0
		for pos < len(rest) {
			c, size := utf8.DecodeRuneInString(rest[pos:])
			if !p.Match(c) {
				break
			}
			pos += size
			runes++
		}

		if runes < n {
			return s, "", RangedParserError{
				Err:  CombinatorParseError{State: s, kind: p.String()},
				Exec: RangedParserExec{Count: runes, Min: n},
				kind: "while_n",
			}
		}

		return s.Advance(pos), rest[:pos], nil
	}
}

// WhileNM will scan the input text, testing each character against the provided
// [Predicate]. The [Predicate] must match a minimum of n and upto a maximum
// of m characters. If n is zero, this becomes an optional combinator.
//
//	chomp.WhileNM(chomp.IsLetter, 1, 8).Run("Hello, World!")
//	// (", World!", "Hello", nil)
func WhileNM(p Predicate, n, m int) Combinator[string] {
	return func(s State) (State, string, error) {
		if n < 0 || m < 0 {
			return s, "", ParserError{Err: fmt.Errorf("chomp: count must be non-negative, got n=%d, m=%d", n, m), kind: "while_n_m"}
		}
		rest := s.Rest()
		pos, runes := 0, 0
		for pos < len(rest) {
			c, size := utf8.DecodeRuneInString(rest[pos:])
			if !p.Match(c) {
				break
			}
			pos += size
			runes++
		}

		if runes < n || runes > m {
			return s, "", RangedParserError{
				Err:  CombinatorParseError{State: s, kind: p.String()},
				Exec: RangedParserExec{Count: runes, Min: n, Max: m},
				kind: "while_n_m",
			}
		}

		return s.Advance(pos), rest[:pos], nil
	}
}

// WhileNot will scan the input text, testing each character against the provided
// [Predicate]. The [Predicate] must not match at least one character. It has
// the inverse behavior of [While].
//
//	chomp.WhileNot(chomp.IsDigit).Run("Hello, World!")
//	// ("", "Hello, World!", nil)
func WhileNot(p Predicate) Combinator[string] {
	return WhileNotN(p, 1)
}

// WhileNotN will scan the input text, testing each character against the provided
// [Predicate]. The [Predicate] must not match at least n characters. If n is
// zero, this becomes an optional combinator. It has the inverse behavior of [WhileN].
//
//	chomp.WhileNotN(chomp.IsDigit, 1).Run("Hello, World!")
//	// ("", "Hello, World!", nil)
//
//	chomp.WhileNotN(chomp.IsLetter, 0).Run("Hello, World!")
//	// ("Hello, World!", "", nil)
func WhileNotN(p Predicate, n int) Combinator[string] {
	return func(s State) (State, string, error) {
		if n < 0 {
			return s, "", ParserError{Err: fmt.Errorf("chomp: count must be non-negative, got %d", n), kind: "while_not_n"}
		}
		rest := s.Rest()
		pos, runes := 0, 0
		for pos < len(rest) {
			c, size := utf8.DecodeRuneInString(rest[pos:])
			if p.Match(c) {
				break
			}
			pos += size
			runes++
		}

		if runes < n {
			return s, "", RangedParserError{
				Err:  CombinatorParseError{State: s, kind: p.String()},
				Exec: RangedParserExec{Count: runes, Min: n},
				kind: "while_not_n",
			}
		}

		return s.Advance(pos), rest[:pos], nil
	}
}

// WhileNotNM will scan the input text, testing each character against the provided
// [Predicate]. The [Predicate] must not match a minimum of n and upto a maximum of
// m characters. If n is zero, this becomes an optional combinator. It has the
// inverse behavior of [WhileNM].
//
//	chomp.WhileNotNM(chomp.IsLetter, 1, 9).Run("20240709 was a great day")
//	// ("was a great day", "20240709 ", nil)
func WhileNotNM(p Predicate, n, m int) Combinator[string] {
	return func(s State) (State, string, error) {
		if n < 0 || m < 0 {
			return s, "", ParserError{Err: fmt.Errorf("chomp: count must be non-negative, got n=%d, m=%d", n, m), kind: "while_not_n_m"}
		}
		rest := s.Rest()
		pos, runes := 0, 0
		for pos < len(rest) {
			c, size := utf8.DecodeRuneInString(rest[pos:])
			if p.Match(c) {
				break
			}
			pos += size
			runes++
		}

		if runes < n || runes > m {
			return s, "", RangedParserError{
				Err:  CombinatorParseError{State: s, kind: p.String()},
				Exec: RangedParserExec{Count: runes, Min: n, Max: m},
				kind: "while_not_n_m",
			}
		}

		return s.Advance(pos), rest[:pos], nil
	}
}

// Alpha matches one or more ASCII or Unicode letters.
// Equivalent to While(IsLetter).
//
//	chomp.Alpha().Run("Hello123")
//	// ("123", "Hello", nil)
func Alpha() Combinator[string] {
	return While(IsLetter)
}

// Alpha0 matches zero or more ASCII or Unicode letters.
// Equivalent to WhileN(IsLetter, 0).
//
//	chomp.Alpha0().Run("123Hello")
//	// ("123Hello", "", nil)
func Alpha0() Combinator[string] {
	return WhileN(IsLetter, 0)
}

// Digit matches one or more decimal digits.
// Equivalent to While(IsDigit).
//
//	chomp.Digit().Run("123abc")
//	// ("abc", "123", nil)
func Digit() Combinator[string] {
	return While(IsDigit)
}

// Digit0 matches zero or more decimal digits.
// Equivalent to WhileN(IsDigit, 0).
//
//	chomp.Digit0().Run("abc123")
//	// ("abc123", "", nil)
func Digit0() Combinator[string] {
	return WhileN(IsDigit, 0)
}

// Alphanumeric matches one or more alphanumeric characters.
// Equivalent to While(IsAlphanumeric).
//
//	chomp.Alphanumeric().Run("Hello123!")
//	// ("!", "Hello123", nil)
func Alphanumeric() Combinator[string] {
	return While(IsAlphanumeric)
}

// Alphanumeric0 matches zero or more alphanumeric characters.
// Equivalent to WhileN(IsAlphanumeric, 0).
//
//	chomp.Alphanumeric0().Run("!Hello123")
//	// ("!Hello123", "", nil)
func Alphanumeric0() Combinator[string] {
	return WhileN(IsAlphanumeric, 0)
}

// Space matches one or more space or tab characters.
// Equivalent to While(IsSpace).
//
//	chomp.Space().Run("   Hello")
//	// ("Hello", "   ", nil)
func Space() Combinator[string] {
	return While(IsSpace)
}

// Space0 matches zero or more space or tab characters.
// Equivalent to WhileN(IsSpace, 0).
//
//	chomp.Space0().Run("Hello")
//	// ("Hello", "", nil)
func Space0() Combinator[string] {
	return WhileN(IsSpace, 0)
}

// Multispace matches one or more whitespace characters (space, tab, newline, carriage return).
// Equivalent to While(IsMultispace).
//
//	chomp.Multispace().Run("  \n\tHello")
//	// ("Hello", "  \n\t", nil)
func Multispace() Combinator[string] {
	return While(IsMultispace)
}

// Multispace0 matches zero or more whitespace characters (space, tab, newline, carriage return).
// Equivalent to WhileN(IsMultispace, 0).
//
//	chomp.Multispace0().Run("Hello")
//	// ("Hello", "", nil)
func Multispace0() Combinator[string] {
	return WhileN(IsMultispace, 0)
}

// HexDigit matches one or more hexadecimal digits (0-9, a-f, A-F).
// Equivalent to While(IsHexDigit).
//
//	chomp.HexDigit().Run("1a2B3c rest")
//	// (" rest", "1a2B3c", nil)
func HexDigit() Combinator[string] {
	return While(IsHexDigit)
}

// HexDigit0 matches zero or more hexadecimal digits (0-9, a-f, A-F).
// Equivalent to WhileN(IsHexDigit, 0).
//
//	chomp.HexDigit0().Run("xyz")
//	// ("xyz", "", nil)
func HexDigit0() Combinator[string] {
	return WhileN(IsHexDigit, 0)
}

// OctalDigit matches one or more octal digits (0-7).
// Equivalent to While(IsOctalDigit).
//
//	chomp.OctalDigit().Run("0127 rest")
//	// (" rest", "0127", nil)
func OctalDigit() Combinator[string] {
	return While(IsOctalDigit)
}

// OctalDigit0 matches zero or more octal digits (0-7).
// Equivalent to WhileN(IsOctalDigit, 0).
//
//	chomp.OctalDigit0().Run("89")
//	// ("89", "", nil)
func OctalDigit0() Combinator[string] {
	return WhileN(IsOctalDigit, 0)
}

// BinaryDigit matches one or more binary digits (0-1).
// Equivalent to While(IsBinaryDigit).
//
//	chomp.BinaryDigit().Run("1010 rest")
//	// (" rest", "1010", nil)
func BinaryDigit() Combinator[string] {
	return While(IsBinaryDigit)
}

// BinaryDigit0 matches zero or more binary digits (0-1).
// Equivalent to WhileN(IsBinaryDigit, 0).
//
//	chomp.BinaryDigit0().Run("234")
//	// ("234", "", nil)
func BinaryDigit0() Combinator[string] {
	return WhileN(IsBinaryDigit, 0)
}

// Newline matches a single newline character '\n'.
//
//	chomp.Newline().Run("\nHello")
//	// ("Hello", "\n", nil)
func Newline() Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		if rest != "" && rest[0] == '\n' {
			return s.Advance(1), "\n", nil
		}
		return s, "", CombinatorParseError{State: s, kind: "newline"}
	}
}

// Tab matches a single tab character '\t'.
//
//	chomp.Tab().Run("\tHello")
//	// ("Hello", "\t", nil)
func Tab() Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		if rest != "" && rest[0] == '\t' {
			return s.Advance(1), "\t", nil
		}
		return s, "", CombinatorParseError{State: s, kind: "tab"}
	}
}

// NotLineEnding matches any characters until a line ending ('\n' or '\r').
// Requires at least one character to be matched.
//
//	chomp.NotLineEnding().Run("Hello, World!\nNext line")
//	// ("\nNext line", "Hello, World!", nil)
func NotLineEnding() Combinator[string] {
	return WhileNot(IsLineEnding)
}

// AnyDigit matches a single decimal digit (0-9).
//
//	chomp.AnyDigit().Run("123")
//	// ("23", "1", nil)
func AnyDigit() Combinator[string] {
	return Satisfy(IsDigit.Match)
}

// AnyLetter matches a single ASCII or Unicode letter.
//
//	chomp.AnyLetter().Run("Hello")
//	// ("ello", "H", nil)
func AnyLetter() Combinator[string] {
	return Satisfy(IsLetter.Match)
}

// AnyAlphanumeric matches a single alphanumeric character.
//
//	chomp.AnyAlphanumeric().Run("a1!")
//	// ("1!", "a", nil)
func AnyAlphanumeric() Combinator[string] {
	return Satisfy(IsAlphanumeric.Match)
}

// AnyHexDigit matches a single hexadecimal digit (0-9, a-f, A-F).
//
//	chomp.AnyHexDigit().Run("fF0")
//	// ("F0", "f", nil)
func AnyHexDigit() Combinator[string] {
	return Satisfy(IsHexDigit.Match)
}

// AnyOctalDigit matches a single octal digit (0-7).
//
//	chomp.AnyOctalDigit().Run("752")
//	// ("52", "7", nil)
func AnyOctalDigit() Combinator[string] {
	return Satisfy(IsOctalDigit.Match)
}

// AnyBinaryDigit matches a single binary digit (0-1).
//
//	chomp.AnyBinaryDigit().Run("101")
//	// ("01", "1", nil)
func AnyBinaryDigit() Combinator[string] {
	return Satisfy(IsBinaryDigit.Match)
}
