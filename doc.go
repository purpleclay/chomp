// Package chomp provides a parser combinator library for chomping strings
// (a rune at a time) in Go. A more intuitive way to parse text without having
// to write a single regex.
//
// # The combinator contract
//
// Every combinator in this package honours the same contract:
//
//  1. Failure is non-consuming. On error, a combinator returns the original
//     input string unchanged (and the zero value for ext).
//  2. Success extraction is a prefix. On success, for Combinator[string],
//     ext is exactly the consumed prefix: input == ext + rem. Combinators
//     that transform their output, or intentionally discard part of the
//     matched text (delimiters, prefixes, suffixes, separators), are
//     documented as such and are exempt from this clause only.
//  3. Zero-width success terminates repetition. A repetition combinator
//     stops iterating when an iteration succeeds without consuming input.
//
// Custom combinators composed from this package should honour the same
// contract to remain safe to use with First, Opt, and the repetition
// combinators, all of which rely on rule 1 to backtrack correctly.
package chomp
