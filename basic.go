package chomp

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Char matches a specific single character at the beginning of the input text.
//
//	chomp.Char(',').Run(",,rest")
//	// (",rest", ",", nil)
func Char(c rune) Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		if rest == "" {
			return s, "", CombinatorParseError{Expected: fmt.Sprintf("%q", c), State: s, Type: "char"}
		}

		r, size := utf8.DecodeRuneInString(rest)
		if r == c {
			return s.Advance(size), rest[:size], nil
		}

		return s, "", CombinatorParseError{Expected: fmt.Sprintf("%q", c), State: s, Type: "char"}
	}
}

// AnyChar matches any single character at the beginning of the input text.
//
//	chomp.AnyChar().Run("Hello")
//	// ("ello", "H", nil)
func AnyChar() Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		if rest == "" {
			return s, "", CombinatorParseError{State: s, Type: "any_char"}
		}

		_, size := utf8.DecodeRuneInString(rest)
		return s.Advance(size), rest[:size], nil
	}
}

// Satisfy matches a single character at the beginning of the input text that
// satisfies the given predicate function.
//
//	chomp.Satisfy(func(r rune) bool { return r >= 'A' && r <= 'Z' }).Run("Hello")
//	// ("ello", "H", nil)
func Satisfy(pred func(rune) bool) Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		if rest == "" {
			return s, "", CombinatorParseError{State: s, Type: "satisfy"}
		}

		r, size := utf8.DecodeRuneInString(rest)
		if pred(r) {
			return s.Advance(size), rest[:size], nil
		}

		return s, "", CombinatorParseError{State: s, Type: "satisfy"}
	}
}

// Tag must match a series of characters at the beginning of the input text
// in the exact order and case provided. An empty str matches trivially
// without consuming any input.
//
//	chomp.Tag("Hello").Run("Hello, World!")
//	// (", World!", "Hello", nil)
func Tag(str string) Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		if strings.HasPrefix(rest, str) {
			return s.Advance(len(str)), str, nil
		}

		return s, "", CombinatorParseError{Expected: fmt.Sprintf("%q", str), State: s, Type: "tag"}
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
//	chomp.TagNoCase("hello").Run("HELLO, World!")
//	// (", World!", "HELLO", nil)
func TagNoCase(str string) Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		pos, tagPos := 0, 0
		for tagPos < len(str) {
			want, wantSize := utf8.DecodeRuneInString(str[tagPos:])
			if pos >= len(rest) || (want == utf8.RuneError && wantSize <= 1) {
				return s, "", CombinatorParseError{Expected: fmt.Sprintf("%q", str), State: s, Type: "tag_no_case"}
			}

			got, size := utf8.DecodeRuneInString(rest[pos:])
			if (got == utf8.RuneError && size <= 1) || !foldEqual(got, want) {
				return s, "", CombinatorParseError{Expected: fmt.Sprintf("%q", str), State: s, Type: "tag_no_case"}
			}
			pos += size
			tagPos += wantSize
		}

		return s.Advance(pos), rest[:pos], nil
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

// IsA must match at least one character from the provided sequence at the
// beginning of the input text. Parsing stops upon the first unmatched character.
// An empty sequence can never satisfy "at least one", so this always fails.
//
//	chomp.IsA("eH").Run("Hello, World!")
//	// ("llo, World!", "He", nil)
func IsA(str string) Combinator[string] {
	runeCount := utf8.RuneCountInString(str)

	if runeCount >= charSetThreshold {
		charSet := make(map[rune]struct{}, runeCount)
		for _, r := range str {
			charSet[r] = struct{}{}
		}

		return func(s State) (State, string, error) {
			rest := s.Rest()
			pos := 0
			for _, sc := range rest {
				if _, ok := charSet[sc]; !ok {
					break
				}
				pos += utf8.RuneLen(sc)
			}

			if pos == 0 {
				return s, "", CombinatorParseError{Expected: fmt.Sprintf("a character in %q", str), State: s, Type: "is_a"}
			}

			return s.Advance(pos), rest[:pos], nil
		}
	}

	return func(s State) (State, string, error) {
		rest := s.Rest()
		pos := 0

	match:
		for _, sc := range rest {
			for _, strc := range str {
				if sc == strc {
					pos += utf8.RuneLen(sc)
					continue match
				}
			}
			break match
		}

		if pos == 0 {
			return s, "", CombinatorParseError{Expected: fmt.Sprintf("a character in %q", str), State: s, Type: "is_a"}
		}

		return s.Advance(pos), rest[:pos], nil
	}
}

// IsNot must not match at least one character at the beginning of the input text
// from the provided sequence. Parsing stops upon the first matched character.
// An empty sequence excludes nothing, so this matches the entire remaining
// input when at least one character remains; an empty input still fails.
//
//	chomp.IsNot("ol").Run("Hello, World!")
//	// ("llo, World!", "He", nil)
func IsNot(str string) Combinator[string] {
	runeCount := utf8.RuneCountInString(str)

	if runeCount >= charSetThreshold {
		charSet := make(map[rune]struct{}, runeCount)
		for _, r := range str {
			charSet[r] = struct{}{}
		}

		return func(s State) (State, string, error) {
			rest := s.Rest()
			pos := 0
			for _, sc := range rest {
				if _, ok := charSet[sc]; ok {
					break
				}
				pos += utf8.RuneLen(sc)
			}

			if pos == 0 {
				return s, "", CombinatorParseError{Expected: fmt.Sprintf("a character not in %q", str), State: s, Type: "is_not"}
			}

			return s.Advance(pos), rest[:pos], nil
		}
	}

	return func(s State) (State, string, error) {
		rest := s.Rest()
		pos := 0

	match:
		for _, sc := range rest {
			for _, strc := range str {
				if sc == strc {
					break match
				}
			}
			pos += utf8.RuneLen(sc)
		}

		if pos == 0 {
			return s, "", CombinatorParseError{Expected: fmt.Sprintf("a character not in %q", str), State: s, Type: "is_not"}
		}

		return s.Advance(pos), rest[:pos], nil
	}
}

// OneOf must match a single character at the beginning of the text from
// the provided sequence. An empty sequence can never be matched, so this
// always fails.
//
//	chomp.OneOf("!,eH").Run("Hello, World!")
//	// ("ello, World!", "H", nil)
func OneOf(str string) Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		if rest == "" {
			return s, "", CombinatorParseError{Expected: fmt.Sprintf("a character in %q", str), State: s, Type: "one_of"}
		}

		r, size := utf8.DecodeRuneInString(rest)
		for _, strc := range str {
			if r == strc {
				return s.Advance(size), rest[:size], nil
			}
		}

		return s, "", CombinatorParseError{Expected: fmt.Sprintf("a character in %q", str), State: s, Type: "one_of"}
	}
}

// NoneOf must not match a single character at the beginning of the text
// from the provided sequence. An empty sequence excludes nothing, so this
// degenerates to matching any single character, the same as [AnyChar].
//
//	chomp.NoneOf("loWrd!e").Run("Hello, World!")
//	// ("ello, World!", "H", nil)
func NoneOf(str string) Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		if rest == "" {
			return s, "", CombinatorParseError{Expected: fmt.Sprintf("a character not in %q", str), State: s, Type: "none_of"}
		}

		r, size := utf8.DecodeRuneInString(rest)
		for _, strc := range str {
			if r == strc {
				return s, "", CombinatorParseError{Expected: fmt.Sprintf("a character not in %q", str), State: s, Type: "none_of"}
			}
		}

		return s.Advance(size), rest[:size], nil
	}
}

// Until will scan the input text for the first occurrence of the provided series
// of characters. Everything until that point in the text will be matched.
// An empty str matches at position 0, so this succeeds trivially without
// consuming any input.
//
//	chomp.Until("World").Run("Hello, World!")
//	// ("World!", "Hello, ", nil)
func Until(str string) Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		if idx := strings.Index(rest, str); idx != -1 {
			return s.Advance(idx), rest[:idx], nil
		}

		return s, "", CombinatorParseError{Expected: fmt.Sprintf("%q", str), State: s, Type: "until"}
	}
}

// Take will consume exactly n characters from the beginning of the input text.
// Unicode characters are handled correctly by counting runes, not bytes.
//
//	chomp.Take(5).Run("Hello, World!")
//	// (", World!", "Hello", nil)
func Take(n int) Combinator[string] {
	return func(s State) (State, string, error) {
		if n < 0 {
			return s, "", ParserError{Err: fmt.Errorf("chomp: count must be non-negative, got %d", n), Type: "take"}
		}
		rest := s.Rest()
		pos := 0
		for i := 0; i < n; i++ {
			if pos >= len(rest) {
				return s, "", CombinatorParseError{State: s, Type: "take"}
			}
			_, size := utf8.DecodeRuneInString(rest[pos:])
			pos += size
		}
		return s.Advance(pos), rest[:pos], nil
	}
}

// TakeUntil1 will scan the input text for the first occurrence of the provided
// series of characters, requiring at least one character to be matched before
// the delimiter. Everything until that point in the text will be matched.
// An empty str always matches at position 0, which never satisfies "at
// least one", so this always fails.
//
//	chomp.TakeUntil1(",").Run("Hello, World!")
//	// (", World!", "Hello", nil)
//
//	chomp.TakeUntil1(",").Run(",World!")
//	// Error: must match at least one character
func TakeUntil1(str string) Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		if idx := strings.Index(rest, str); idx > 0 {
			return s.Advance(idx), rest[:idx], nil
		}

		return s, "", CombinatorParseError{Expected: fmt.Sprintf("%q", str), State: s, Type: "take_until_1"}
	}
}

// Escaped parses a string containing escape sequences. It takes a normal content
// combinator, an escape character, and a combinator that matches valid characters
// after the escape. The escape sequences are preserved in the output as-is.
//
//	chomp.Escaped(chomp.While(chomp.IsLetter), '\\', chomp.OneOf(`"n\`)).Run(`Hello\"World`)
//	// ("", `Hello\"World`, nil)
func Escaped(normal Combinator[string], escape rune, escapable Combinator[string]) Combinator[string] {
	return func(s State) (State, string, error) {
		cur := s

		for cur.Rest() != "" {
			if newCur, _, err := normal(cur); err == nil && newCur.Pos() > cur.Pos() {
				cur = newCur
				continue
			}

			r, _ := utf8.DecodeRuneInString(cur.Rest())
			if r == escape {
				escLen := utf8.RuneLen(escape)
				if len(cur.Rest()) <= escLen {
					break
				}

				escState := cur.Advance(escLen)
				if newCur, _, err := escapable(escState); err == nil && newCur.Pos() > escState.Pos() {
					cur = newCur
					continue
				}
			}

			break
		}

		if cur.Pos() == s.Pos() {
			return s, "", CombinatorParseError{State: s, Type: "escaped"}
		}

		return cur, cur.since(s), nil
	}
}

// EscapedTransform parses a string containing escape sequences and transforms them.
// It takes a normal content combinator, an escape character, and a transform function
// that converts escape sequences to their actual values.
//
//	transform := func(s chomp.State) (chomp.State, string, error) {
//	    switch s.Rest()[0] {
//	    case 'n':
//	        return s.Advance(1), "\n", nil
//	    case '"':
//	        return s.Advance(1), "\"", nil
//	    case '\\':
//	        return s.Advance(1), "\\", nil
//	    }
//	    return s, "", errors.New("invalid escape")
//	}
//	chomp.EscapedTransform(chomp.While(chomp.IsLetter), '\\', transform).Run(`Hello\nWorld`)
//	// ("", "Hello\nWorld", nil)
func EscapedTransform(normal Combinator[string], escape rune, transform Combinator[string]) Combinator[string] {
	return func(s State) (State, string, error) {
		var result strings.Builder
		cur := s

		for cur.Rest() != "" {
			if newCur, ext, err := normal(cur); err == nil && newCur.Pos() > cur.Pos() {
				result.WriteString(ext)
				cur = newCur
				continue
			}

			r, _ := utf8.DecodeRuneInString(cur.Rest())
			if r == escape {
				escLen := utf8.RuneLen(escape)
				if len(cur.Rest()) <= escLen {
					break
				}

				escState := cur.Advance(escLen)
				if newCur, transformed, err := transform(escState); err == nil && newCur.Pos() > escState.Pos() {
					result.WriteString(transformed)
					cur = newCur
					continue
				}
			}

			break
		}

		if result.Len() == 0 {
			return s, "", CombinatorParseError{State: s, Type: "escaped_transform"}
		}

		return cur, result.String(), nil
	}
}
