package chomp

import "strings"

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
		return Suffixed(WhileNotN(IsLineEnding, 0), Opt(Crlf()))(s)
	}
}
