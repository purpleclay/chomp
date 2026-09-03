# Chomp

![Nix](https://img.shields.io/badge/Nix-5277C3?logo=nixos&logoColor=white)
![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)
[![MIT](https://img.shields.io/badge/MIT-gray?logo=github&logoColor=white)](LICENSE)

A parser combinator library for Go that makes parsing text intuitive and maintainable. Stop wrestling with regex and start writing parsers that read like natural grammar.

> Inspired by [nom](https://github.com/rust-bakery/nom) 💜.

## Why Chomp?

Parser combinators offer significant advantages over regular expressions:

| | Chomp | Regex |
|---|-------|-------|
| **Readability** | Reads like grammar rules | Often "write-only" patterns |
| **Composability** | Build complex parsers from simple, reusable pieces | Monolithic patterns that resist reuse |
| **Error Messages** | Clear context on what failed and where | Generic "no match" or cryptic positions |
| **Maintainability** | Easy to modify and extend | Small changes can break everything |
| **Nested Structures** | Natural support for recursion | Struggles or impossible |
| **Type Safety** | Compile-time guarantees | Runtime string manipulation |

## Installation

```sh
go get github.com/purpleclay/chomp
```

## How It Works

At the heart of `chomp` is the **combinator** - a function that attempts to parse a [State](https://pkg.go.dev/github.com/purpleclay/chomp#State) (the original input plus a cursor) and returns a tuple `(rem, ext, err)`. `Run` is the string-in/string-out entry point for a top-level parse:

```
                       input
                         │
                         ▼
              ┌─────────────────────┐
              │     Combinator      │
              └─────────────────────┘
                         │
          ┌──────────────┼──────────────┐
          ▼              ▼              ▼
    ┌───────────┐  ┌───────────┐  ┌───────────┐
    │    rem    │  │    ext    │  │    err    │
    └───────────┘  └───────────┘  └───────────┘
      remaining      extracted    error (if any)
        text           text
```

```go
// Parse a simple tag
rem, ext, _ := chomp.Tag("Hello").Run("Hello, World!")
// ext: "Hello"
// rem: ", World!"
```

Combinators can be composed together to build sophisticated parsers:

```go
// Parse a key-value pair like "name=alice"
func KeyValue() chomp.Combinator[chomp.Tuple2[string, string]] {
    return chomp.SepPair(
        chomp.While(chomp.IsLetter),  // key: letters
        chomp.Tag("="),               // separator (discarded)
        chomp.While(chomp.IsLetter),  // value: letters
    )
}

rem, kv, _ := KeyValue().Run("name=alice&age=30")
// kv: {First: "name", Second: "alice"}
// rem: "&age=30"
```

## The Combinator Contract

Every combinator honours a single documented contract, so combinators compose predictably regardless of who wrote them:

1. **Failure is non-consuming.** On error, a combinator returns the `State` it was given unchanged (and the zero value for `ext`).
2. **Success extraction is a prefix.** On success, for a `Combinator[string]`, `ext` is exactly the consumed prefix: `input == ext + rem`. Combinators that transform their output, or intentionally discard part of the matched text (delimiters, prefixes, suffixes, separators), are documented as such and are exempt from this clause only.
3. **Zero-width success terminates repetition.** A repetition combinator stops iterating when an iteration succeeds without consuming input.

If you write your own combinators, follow the same rules — `First`, `Opt`, and the repetition combinators all rely on rule 1 to backtrack correctly.

## Examples

Real-world parser examples:

- [GPG Private Key Parser](examples/gpg/main.go) - Parse GPG key metadata
- [Git Diff Parser](examples/git-diff/main.go) - Parse unified diff output

## Documentation

- [Combinator Reference](docs/combinators.md) - All available combinators
- [Benchmarks](docs/benchmarks.md) - Performance benchmarks
