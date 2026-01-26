package chomp

import (
	"strings"
	"unicode/utf8"
)

// Crlf must match either a CR '\r' or CRLF '\r\n' line ending.
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

// Eol will scan and return any text before any ASCII line ending
// characters. Line endings are discarded.
//
//	chomp.Eol()(`Hello, World!\nIt's a great day!`)
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
