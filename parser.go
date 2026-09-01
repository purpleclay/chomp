package chomp

import (
	"unicode/utf8"
)

// Crlf must match a strict CRLF '\r\n' line ending. It never matches a bare
// LF or a bare CR; see [LineEnding] to also accept a bare LF, and [Eol] for
// how a bare CR is otherwise handled. Inspects at most the first two bytes
// of the input, regardless of its length.
//
//	chomp.Crlf().Run("\r\nHello")
//	// ("Hello", "\r\n", nil)
func Crlf() Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		if len(rest) >= 2 && rest[0] == '\r' && rest[1] == '\n' {
			return s.Advance(2), rest[:2], nil
		}

		return s, "", CombinatorParseError{State: s, Type: "crlf"}
	}
}

// LineEnding must match either a LF '\n' or CRLF '\r\n' line ending. A bare
// CR '\r' is never matched; see [Eol] if you need to treat a bare CR as a
// legacy-Mac line terminator. Inspects at most the first two bytes of the
// input, regardless of its length.
//
//	chomp.LineEnding().Run("\nHello")
//	// ("Hello", "\n", nil)
//
//	chomp.LineEnding().Run("\r\nHello")
//	// ("Hello", "\r\n", nil)
func LineEnding() Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		if len(rest) >= 1 && rest[0] == '\n' {
			return s.Advance(1), rest[:1], nil
		}
		if len(rest) >= 2 && rest[0] == '\r' && rest[1] == '\n' {
			return s.Advance(2), rest[:2], nil
		}

		return s, "", CombinatorParseError{State: s, Type: "line_ending"}
	}
}

// Eol will scan and return any text before any ASCII line ending
// characters. Line endings are discarded. Unlike [LineEnding], a bare CR
// '\r' is also consumed here, treated as a legacy-Mac line terminator.
//
//	chomp.Eol().Run("Hello, World!\nIt's a great day!")
//	// ("It's a great day!", "Hello, World!", nil)
func Eol() Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		pos := 0
		for _, c := range rest {
			if c == '\n' || c == '\r' {
				break
			}
			pos += utf8.RuneLen(c)
		}

		matched := rest[:pos]
		lineEnd := rest[pos:]
		skip := 0
		if lineEnd != "" {
			if lineEnd[0] == '\n' {
				skip = 1
			} else if len(lineEnd) >= 2 && lineEnd[0] == '\r' && lineEnd[1] == '\n' {
				skip = 2
			} else if lineEnd[0] == '\r' {
				skip = 1
			}
		}

		return s.Advance(pos + skip), matched, nil
	}
}
