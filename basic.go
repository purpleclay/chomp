package chomp

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Char matches a specific single character at the beginning of the input text.
//
//	chomp.Char(',')(",,rest")
//	// (",rest", ",", nil)
func Char(c rune) Combinator[string] {
	return func(s string) (string, string, error) {
		if s == "" {
			return s, "", CombinatorParseError{Input: string(c), Text: s, Type: "char"}
		}

		r, size := utf8.DecodeRuneInString(s)
		if r == c {
			return s[size:], s[:size], nil
		}

		return s, "", CombinatorParseError{Input: fmt.Sprintf("%c", c), Text: s, Type: "char"}
	}
}

// AnyChar matches any single character at the beginning of the input text.
//
//	chomp.AnyChar()("Hello")
//	// ("ello", "H", nil)
func AnyChar() Combinator[string] {
	return func(s string) (string, string, error) {
		if runes := []rune(s); len(runes) > 0 {
			matched := string(runes[0])
			return s[len(matched):], matched, nil
		}

		return s, "", CombinatorParseError{Text: s, Type: "any_char"}
	}
}

// Satisfy matches a single character at the beginning of the input text that
// satisfies the given predicate function.
//
//	chomp.Satisfy(func(r rune) bool { return r >= 'A' && r <= 'Z' })("Hello")
//	// ("ello", "H", nil)
func Satisfy(pred func(rune) bool) Combinator[string] {
	return func(s string) (string, string, error) {
		if runes := []rune(s); len(runes) > 0 && pred(runes[0]) {
			matched := string(runes[0])
			return s[len(matched):], matched, nil
		}

		return s, "", CombinatorParseError{Text: s, Type: "satisfy"}
	}
}

// Tag must match a series of characters at the beginning of the input text
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

// TagNoCase must match a series of characters at the beginning of the input text
// in the exact order provided, but ignoring case. The matched text from the input
// is returned (preserving the original casing).
//
//	chomp.TagNoCase("hello")("HELLO, World!")
//	// (", World!", "HELLO", nil)
func TagNoCase(str string) Combinator[string] {
	return func(s string) (string, string, error) {
		if len(s) >= len(str) && strings.EqualFold(s[:len(str)], str) {
			return s[len(str):], s[:len(str)], nil
		}

		return s, "", CombinatorParseError{Input: str, Text: s, Type: "tag_no_case"}
	}
}

// Any must match at least one character from the provided sequence at the
// beginning of the input text. Parsing stops upon the first unmatched character.
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

// Not must not match at least one character at the beginning of the input text
// from the provided sequence. Parsing stops upon the first matched character.
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

// OneOf must match a single character at the beginning of the text from
// the provided sequence.
//
//	chomp.OneOf("!,eH")("Hello, World!")
//	// ("ello, World!", "H", nil)
func OneOf(str string) Combinator[string] {
	return func(s string) (string, string, error) {
		if s == "" {
			return s, "", CombinatorParseError{Input: str, Text: s, Type: "one_of"}
		}

		r, size := utf8.DecodeRuneInString(s)
		for _, strc := range str {
			if r == strc {
				return s[size:], s[:size], nil
			}
		}

		return s, "", CombinatorParseError{Input: str, Text: s, Type: "one_of"}
	}
}

// NoneOf must not match a single character at the beginning of the text
// from the provided sequence.
//
//	chomp.NoneOf("loWrd!e")("Hello, World!")
//	// ("ello, World!", "H", nil)
func NoneOf(str string) Combinator[string] {
	return func(s string) (string, string, error) {
		if s == "" {
			return s, "", CombinatorParseError{Input: str, Text: s, Type: "none_of"}
		}

		r, size := utf8.DecodeRuneInString(s)
		for _, strc := range str {
			if r == strc {
				return s, "", CombinatorParseError{Input: str, Text: s, Type: "none_of"}
			}
		}

		return s[size:], s[:size], nil
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

// Take will consume exactly n characters from the beginning of the input text.
// Unicode characters are handled correctly by counting runes, not bytes.
//
//	chomp.Take(5)("Hello, World!")
//	// (", World!", "Hello", nil)
func Take(n uint) Combinator[string] {
	return func(s string) (string, string, error) {
		runes := []rune(s)
		if uint(len(runes)) < n {
			return s, "", CombinatorParseError{Text: s, Type: "take"}
		}

		taken := string(runes[:n])
		return s[len(taken):], taken, nil
	}
}

// TakeUntil1 will scan the input text for the first occurrence of the provided
// series of characters, requiring at least one character to be matched before
// the delimiter. Everything until that point in the text will be matched.
//
//	chomp.TakeUntil1(",")("Hello, World!")
//	// (", World!", "Hello", nil)
//
//	chomp.TakeUntil1(",")(",World!")
//	// Error: must match at least one character
func TakeUntil1(str string) Combinator[string] {
	return func(s string) (string, string, error) {
		if idx := strings.Index(s, str); idx > 0 {
			return s[idx:], s[:idx], nil
		}

		return s, "", CombinatorParseError{Input: str, Text: s, Type: "take_until_1"}
	}
}

// Escaped parses a string containing escape sequences. It takes a normal content
// combinator, an escape character, and a combinator that matches valid characters
// after the escape. The escape sequences are preserved in the output as-is.
//
//	chomp.Escaped(chomp.While(chomp.IsLetter), '\\', chomp.OneOf(`"n\`))(`Hello\"World`)
//	// ("", `Hello\"World`, nil)
func Escaped(normal Combinator[string], escape rune, escapable Combinator[string]) Combinator[string] {
	return func(s string) (string, string, error) {
		pos := 0
		rem := s

		for rem != "" {
			if newRem, ext, err := normal(rem); err == nil && ext != "" {
				pos += len(ext)
				rem = newRem
				continue
			}

			runes := []rune(rem)
			if len(runes) > 0 && runes[0] == escape {
				escLen := len(string(escape))
				if len(rem) <= escLen {
					break
				}

				if _, ext, err := escapable(rem[escLen:]); err == nil && ext != "" {
					pos += escLen + len(ext)
					rem = rem[escLen+len(ext):]
					continue
				}
			}

			break
		}

		if pos == 0 {
			return s, "", CombinatorParseError{Text: s, Type: "escaped"}
		}

		return s[pos:], s[:pos], nil
	}
}

// EscapedTransform parses a string containing escape sequences and transforms them.
// It takes a normal content combinator, an escape character, and a transform function
// that converts escape sequences to their actual values.
//
//	transform := func(s string) (string, string, error) {
//	    switch s[0] {
//	    case 'n':
//	        return s[1:], "\n", nil
//	    case '"':
//	        return s[1:], "\"", nil
//	    case '\\':
//	        return s[1:], "\\", nil
//	    }
//	    return s, "", errors.New("invalid escape")
//	}
//	chomp.EscapedTransform(chomp.While(chomp.IsLetter), '\\', transform)(`Hello\nWorld`)
//	// ("", "Hello\nWorld", nil)
func EscapedTransform(normal Combinator[string], escape rune, transform Combinator[string]) Combinator[string] {
	return func(s string) (string, string, error) {
		var result strings.Builder
		rem := s

		for rem != "" {
			if newRem, ext, err := normal(rem); err == nil && ext != "" {
				result.WriteString(ext)
				rem = newRem
				continue
			}

			runes := []rune(rem)
			if len(runes) > 0 && runes[0] == escape {
				escLen := len(string(escape))
				if len(rem) <= escLen {
					break
				}

				if newRem, transformed, err := transform(rem[escLen:]); err == nil && transformed != "" {
					result.WriteString(transformed)
					rem = newRem
					continue
				}
			}

			break
		}

		if result.Len() == 0 {
			return s, "", CombinatorParseError{Text: s, Type: "escaped_transform"}
		}

		return rem, result.String(), nil
	}
}
