package chomp

import (
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
