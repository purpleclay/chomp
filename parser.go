package chomp

import (
	"unicode/utf8"
)

// Crlf must match a strict CRLF '\r\n' line ending. It never matches a bare
// LF or a bare CR; see [LineEnding] to also accept a bare LF, and [Eol] for
// how a bare CR is otherwise handled. Inspects at most the first two bytes
// of the input, regardless of its length.
//
//	chomp.Crlf()("\r\nHello")
//	// ("Hello", "\r\n", nil)
func Crlf() Combinator[string] {
	return func(s string) (string, string, error) {
		if len(s) >= 2 && s[0] == '\r' && s[1] == '\n' {
			return s[2:], s[:2], nil
		}

		return s, "", CombinatorParseError{Text: s, Type: "crlf"}
	}
}

// LineEnding must match either a LF '\n' or CRLF '\r\n' line ending. A bare
// CR '\r' is never matched; see [Eol] if you need to treat a bare CR as a
// legacy-Mac line terminator. Inspects at most the first two bytes of the
// input, regardless of its length.
//
//	chomp.LineEnding()("\nHello")
//	// ("Hello", "\n", nil)
//
//	chomp.LineEnding()("\r\nHello")
//	// ("Hello", "\r\n", nil)
func LineEnding() Combinator[string] {
	return func(s string) (string, string, error) {
		if len(s) >= 1 && s[0] == '\n' {
			return s[1:], s[:1], nil
		}
		if len(s) >= 2 && s[0] == '\r' && s[1] == '\n' {
			return s[2:], s[:2], nil
		}

		return s, "", CombinatorParseError{Text: s, Type: "line_ending"}
	}
}

// Eol will scan and return any text before any ASCII line ending
// characters. Line endings are discarded. Unlike [LineEnding], a bare CR
// '\r' is also consumed here, treated as a legacy-Mac line terminator.
//
//	chomp.Eol()("Hello, World!\nIt's a great day!")
//	// ("It's a great day!", "Hello, World!", nil)
func Eol() Combinator[string] {
	return func(s string) (string, string, error) {
		pos := 0
		for _, c := range s {
			if c == '\n' || c == '\r' {
				break
			}
			pos += utf8.RuneLen(c)
		}

		rem := s[pos:]
		matched := s[:pos]
		if rem != "" {
			if rem[0] == '\n' {
				rem = rem[1:]
			} else if len(rem) >= 2 && rem[0] == '\r' && rem[1] == '\n' {
				rem = rem[2:]
			} else if rem[0] == '\r' {
				rem = rem[1:]
			}
		}

		return rem, matched, nil
	}
}
