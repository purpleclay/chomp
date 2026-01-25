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

At the heart of `chomp` is the **combinator** - a function that attempts to parse text and returns a tuple `(rem, ext, err)`:

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
rem, ext, _ := chomp.Tag("Hello")("Hello, World!")
// ext: "Hello"
// rem: ", World!"
```

Combinators can be composed together to build sophisticated parsers:

```go
// Parse a key-value pair like "name=alice"
func KeyValue() chomp.Combinator[[]string] {
    return chomp.SepPair(
        chomp.While(chomp.IsLetter),  // key: letters
        chomp.Tag("="),               // separator (discarded)
        chomp.While(chomp.IsLetter),  // value: letters
    )
}

rem, kv, _ := KeyValue()("name=alice&age=30")
// kv: ["name", "alice"]
// rem: "&age=30"
```

## Examples

Real-world parser examples:

- [GPG Private Key Parser](examples/gpg/main.go) - Parse GPG key metadata
- [Git Diff Parser](examples/git-diff/main.go) - Parse unified diff output

## Documentation

- [Combinator Reference](docs/combinators.md) - All available combinators
- [Benchmarks](docs/benchmarks.md) - Performance benchmarks
