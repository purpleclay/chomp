# Chomp

A parser combinator library for chomping strings (_a rune at a time_) in Go. A more intuitive way to parse text without having to write a single regex. Happy to chomp both ASCII and Unicode (_it all tastes the same_).

Inspired by [nom](https://github.com/rust-bakery/nom) 💜.

## Design

At the heart of `chomp` is a combinator. A higher-order function capable of parsing text under a defined condition and returning a tuple `(1,2,3)`:

- `1`: the remaining unparsed (_or unchomped_) text.
- `2`: the parsed (_or chomped_) text.
- `3`: an error if the combinator failed to parse.

Here's a sneak peek at its definition:

```go
type Result interface {
	string | []string
}

type Combinator[T Result] func(string) (string, T, error)
```

A combinator in its simplest form would look like this:

```go
func Tag(str string) chomp.Combinator[string] {
	return func(s string) (string, string, error) {
		if strings.HasPrefix(s, str) {
			// Return a tuple containing:
			// 1. the remaining string after the prefix
			// 2. the matched prefix
			// 3. no error
			return s[len(str):], str, nil
		}

		return s, "", chomp.CombinatorParseError{
			Input: str,
			Text: s,
			Type: "tag",
		}
	}
}
```

The true power of `chomp` comes from the ability to build parsers by chaining (_or combining_) combinators together.

## Writing a Parser Combinator

Take a look at one of the examples of how to write a parser combinator.

1. [GPG Private Key parser]()

## Why use Chomp?

- Combinators are very easy to write and combine into more complex parsers.
- Code written with chomp looks like natural grammar and is easy to understand, maintain and extend.
- It is incredibly easy to unit test.

## Badges

[![Build status](https://img.shields.io/github/actions/workflow/status/purpleclay/chomp/ci.yml?style=flat-square&logo=go)](https://github.com/purpleclay/chomp/actions?workflow=ci)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/purpleclay/chomp?style=flat-square)](https://goreportcard.com/report/github.com/purpleclay/chomp)
[![Go Version](https://img.shields.io/github/go-mod/go-version/purpleclay/chomp.svg?style=flat-square)](go.mod)
[![DeepSource](https://app.deepsource.com/gh/purpleclay/chomp.svg/?label=active+issues&show_trend=false&token=DFB8RRar8iHJrVaNF7e9JaVm)](https://app.deepsource.com/gh/purpleclay/chomp/)
