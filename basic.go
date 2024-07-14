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

		return s, "", CombinatorParseError{Text: s, Type: "crlf"}
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

		return s, "", CombinatorParseError{Input: str, Text: s, Type: "one_of"}
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

		return s, "", CombinatorParseError{Input: str, Text: s, Type: "none_of"}
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

		return s, "", CombinatorParseError{Input: str, Text: s, Type: "until"}
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

// Prefixed will firstly scan the input text for a defined prefix and discard it.
// The remaining input text will be matched against the [Combinator] and returned
// if successful. Both combinators must match.
//
//	chomp.Prefixed(
//		chomp.Tag("Hello"),
//		chomp.Tag(`"`))(`"Hello, World!"`)
//	// (`, World!"`, "Hello", nil)
func Prefixed(c, pre Combinator[string]) Combinator[string] {
	return func(s string) (string, string, error) {
		rem, _, err := pre(s)
		if err != nil {
			return rem, "", err
		}

		return c(rem)
	}
}

// Suffixed will firstly scan the input text and match it against the [Combinator].
// The remaining text will be scanned for a defined suffix and discarded. Both
// combinators must match.
//
//	chomp.Suffixed(
//		chomp.Tag("Hello"),
//		chomp.Tag(", "))("Hello, World!")
//	// ("World!", "Hello", nil)
func Suffixed(c, suf Combinator[string]) Combinator[string] {
	return func(s string) (string, string, error) {
		rem, ext, err := c(s)
		if err != nil {
			return rem, "", err
		}

		rem, _, err = suf(rem)
		if err != nil {
			return rem, "", err
		}

		return rem, ext, nil
	}
}

// Eol will scan the text until it encounters any ASCII line ending characters
// identified by the [IsLineEnding] predicate. All text before the line ending
// will be returned. The line ending, if detected, will be discarded.
//
//	chomp.Eol()(`Hello, World!\nIt's a great day!`)
//	// ("It's a great day!", "Hello, World!", nil)
func Eol() Combinator[string] {
	return func(s string) (string, string, error) {
		return Suffixed(WhileNotN(IsLineEnding, 0), Opt(Crlf()))(s)
	}
}

// Peek will scan the text and apply the parser without consuming any of the input.
// Useful if you need to lookahead.
//
//	chomp.Peek(chomp.Tag("Hello"))("Hello, World!")
//	// ("Hello, World!", "Hello", nil)
//
//	chomp.Peek(
//		chomp.Many(chomp.Suffixed(chomp.Tag(" "), chomp.Until(" "))),
//	)("Hello and Good Morning!")
//	// ("Hello and Good Morning!", []string{"Hello", "and", "Good"}, nil)
func Peek[T Result](c Combinator[T]) Combinator[T] {
	return func(s string) (string, T, error) {
		_, ext, err := c(s)
		return s, ext, err
	}
}
