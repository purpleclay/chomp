package chomp

import (
	"strings"
	"unicode"
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

		return s, "", CombinatorParseError{Input: string(c), Text: s, Type: "char"}
	}
}

// AnyChar matches any single character at the beginning of the input text.
//
//	chomp.AnyChar()("Hello")
//	// ("ello", "H", nil)
func AnyChar() Combinator[string] {
	return func(s string) (string, string, error) {
		if s == "" {
			return s, "", CombinatorParseError{Text: s, Type: "any_char"}
		}

		_, size := utf8.DecodeRuneInString(s)
		return s[size:], s[:size], nil
	}
}

// Satisfy matches a single character at the beginning of the input text that
// satisfies the given predicate function.
//
//	chomp.Satisfy(func(r rune) bool { return r >= 'A' && r <= 'Z' })("Hello")
//	// ("ello", "H", nil)
func Satisfy(pred func(rune) bool) Combinator[string] {
	return func(s string) (string, string, error) {
		if s == "" {
			return s, "", CombinatorParseError{Text: s, Type: "satisfy"}
		}

		r, size := utf8.DecodeRuneInString(s)
		if pred(r) {
			return s[size:], s[:size], nil
		}

		return s, "", CombinatorParseError{Text: s, Type: "satisfy"}
	}
}

// Tag must match a series of characters at the beginning of the input text
// in the exact order and case provided. An empty str matches trivially
// without consuming any input.
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
// in the exact order provided, using full Unicode case-folding (not just
// ASCII). Runes are compared via their case-fold orbit, so fold pairs that
// don't share the same encoded byte length (the Kelvin sign 'K' U+212A
// folds with 'k'/'K', despite being 3 bytes to their 1) still match.
// Malformed UTF-8 in either the input or str never matches, even against a
// literal U+FFFD on the other side, since invalid bytes decode to the same
// replacement-character sentinel as a genuine one. The matched text from
// the input is returned (preserving the original casing). An empty str
// matches trivially without consuming any input, the same as Tag.
//
//	chomp.TagNoCase("hello")("HELLO, World!")
//	// (", World!", "HELLO", nil)
func TagNoCase(str string) Combinator[string] {
	return func(s string) (string, string, error) {
		pos, tagPos := 0, 0
		for tagPos < len(str) {
			want, wantSize := utf8.DecodeRuneInString(str[tagPos:])
			if pos >= len(s) || (want == utf8.RuneError && wantSize <= 1) {
				return s, "", CombinatorParseError{Input: str, Text: s, Type: "tag_no_case"}
			}

			got, size := utf8.DecodeRuneInString(s[pos:])
			if (got == utf8.RuneError && size <= 1) || !foldEqual(got, want) {
				return s, "", CombinatorParseError{Input: str, Text: s, Type: "tag_no_case"}
			}
			pos += size
			tagPos += wantSize
		}

		return s[pos:], s[:pos], nil
	}
}

// foldEqual reports whether a and b are equal under simple Unicode case
// folding, walking a's case-fold orbit rather than assuming a fixed set of
// case pairs.
func foldEqual(a, b rune) bool {
	if a == b {
		return true
	}

	for f := unicode.SimpleFold(a); f != a; f = unicode.SimpleFold(f) {
		if f == b {
			return true
		}
	}

	return false
}

// charSetThreshold is the minimum charset size at which a map lookup
// becomes more efficient than linear scanning
const charSetThreshold = 8

// Any must match at least one character from the provided sequence at the
// beginning of the input text. Parsing stops upon the first unmatched character.
// An empty sequence can never satisfy "at least one", so this always fails.
//
//	chomp.Any("eH")("Hello, World!")
//	// ("llo, World!", "He", nil)
func Any(str string) Combinator[string] {
	runeCount := utf8.RuneCountInString(str)

	if runeCount >= charSetThreshold {
		charSet := make(map[rune]struct{}, runeCount)
		for _, r := range str {
			charSet[r] = struct{}{}
		}

		return func(s string) (string, string, error) {
			pos := 0
			for _, sc := range s {
				if _, ok := charSet[sc]; !ok {
					break
				}
				pos += utf8.RuneLen(sc)
			}

			if pos == 0 {
				return s, "", CombinatorParseError{Input: str, Text: s, Type: "any"}
			}

			return s[pos:], s[:pos], nil
		}
	}

	return func(s string) (string, string, error) {
		pos := 0

	match:
		for _, sc := range s {
			for _, strc := range str {
				if sc == strc {
					pos += utf8.RuneLen(sc)
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
// An empty sequence excludes nothing, so this matches the entire remaining
// input when at least one character remains; an empty input still fails.
//
//	chomp.Not("ol")("Hello, World!")
//	// ("llo, World!", "He", nil)
func Not(str string) Combinator[string] {
	runeCount := utf8.RuneCountInString(str)

	if runeCount >= charSetThreshold {
		charSet := make(map[rune]struct{}, runeCount)
		for _, r := range str {
			charSet[r] = struct{}{}
		}

		return func(s string) (string, string, error) {
			pos := 0
			for _, sc := range s {
				if _, ok := charSet[sc]; ok {
					break
				}
				pos += utf8.RuneLen(sc)
			}

			if pos == 0 {
				return s, "", CombinatorParseError{Input: str, Text: s, Type: "not"}
			}

			return s[pos:], s[:pos], nil
		}
	}

	return func(s string) (string, string, error) {
		pos := 0

	match:
		for _, sc := range s {
			for _, strc := range str {
				if sc == strc {
					break match
				}
			}
			pos += utf8.RuneLen(sc)
		}

		if pos == 0 {
			return s, "", CombinatorParseError{Input: str, Text: s, Type: "not"}
		}

		return s[pos:], s[:pos], nil
	}
}

// OneOf must match a single character at the beginning of the text from
// the provided sequence. An empty sequence can never be matched, so this
// always fails.
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
// from the provided sequence. An empty sequence excludes nothing, so this
// degenerates to matching any single character, the same as [AnyChar].
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
// An empty str matches at position 0, so this succeeds trivially without
// consuming any input.
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
		pos := 0
		for i := uint(0); i < n; i++ {
			if pos >= len(s) {
				return s, "", CombinatorParseError{Text: s, Type: "take"}
			}
			_, size := utf8.DecodeRuneInString(s[pos:])
			pos += size
		}
		return s[pos:], s[:pos], nil
	}
}

// TakeUntil1 will scan the input text for the first occurrence of the provided
// series of characters, requiring at least one character to be matched before
// the delimiter. Everything until that point in the text will be matched.
// An empty str always matches at position 0, which never satisfies "at
// least one", so this always fails.
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
			if newRem, _, err := normal(rem); err == nil && len(newRem) < len(rem) {
				pos += len(rem) - len(newRem)
				rem = newRem
				continue
			}

			r, _ := utf8.DecodeRuneInString(rem)
			if r == escape {
				escLen := utf8.RuneLen(escape)
				if len(rem) <= escLen {
					break
				}

				escInput := rem[escLen:]
				if newRem, _, err := escapable(escInput); err == nil && len(newRem) < len(escInput) {
					pos += escLen + (len(escInput) - len(newRem))
					rem = newRem
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
			if newRem, ext, err := normal(rem); err == nil && len(newRem) < len(rem) {
				result.WriteString(ext)
				rem = newRem
				continue
			}

			r, _ := utf8.DecodeRuneInString(rem)
			if r == escape {
				escLen := utf8.RuneLen(escape)
				if len(rem) <= escLen {
					break
				}

				escInput := rem[escLen:]
				if newRem, transformed, err := transform(escInput); err == nil && len(newRem) < len(escInput) {
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
