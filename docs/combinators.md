# Combinator Reference

A comprehensive reference of all combinators provided by `chomp`.

Every combinator has the shape `func(chomp.State) (chomp.State, T, error)` - `State` threads the original input plus a cursor, so any point in a parse can recover its absolute byte position. `Run(input string)` is the string-in/string-out entry point used throughout this reference:

```go
chomp.Tag("Hello").Run("Hello, World!")
// rem: ", World!"
// ext: "Hello"
```

- [Character Combinators](#character-combinators)
- [Tag Combinators](#tag-combinators)
- [Sequence Combinators](#sequence-combinators)
- [Control Flow Combinators](#control-flow-combinators)
- [Predicate Combinators](#predicate-combinators)
- [Convenience Combinators](#convenience-combinators)
- [Modifier Combinators](#modifier-combinators)
- [Parsers](#parsers)
- [Predicates](#predicates)
- [Error Reporting](#error-reporting)

---

## Character Combinators

### Char

Matches a specific single character at the beginning of the input text.

```go
chomp.Char(',').Run(",,rest")
// rem: ",rest"
// ext: ","
```

### AnyChar

Matches any single character at the beginning of the input text.

```go
chomp.AnyChar().Run("Hello")
// rem: "ello"
// ext: "H"
```

### Satisfy

Matches a single character that satisfies the given predicate function.

```go
chomp.Satisfy(func(r rune) bool { return r >= 'A' && r <= 'Z' }).Run("Hello")
// rem: "ello"
// ext: "H"
```

### OneOf

Matches a single character from the provided sequence.

```go
chomp.OneOf("!,eH").Run("Hello, World!")
// rem: "ello, World!"
// ext: "H"
```

### NoneOf

Matches a single character NOT in the provided sequence.

```go
chomp.NoneOf("loWrd!e").Run("Hello, World!")
// rem: "ello, World!"
// ext: "H"
```

### Any

Matches one or more characters from the provided sequence. Stops at the first unmatched character.

```go
chomp.Any("eH").Run("Hello, World!")
// rem: "llo, World!"
// ext: "He"
```

### Not

Matches one or more characters NOT in the provided sequence. Stops at the first matched character.

```go
chomp.Not("ol").Run("Hello, World!")
// rem: "llo, World!"
// ext: "He"
```

### Take

Consumes exactly `n` characters. Unicode characters are handled correctly.

```go
chomp.Take(5).Run("Hello, World!")
// rem: ", World!"
// ext: "Hello"
```

---

## Tag Combinators

### Tag

Matches an exact sequence of characters (case-sensitive).

```go
chomp.Tag("Hello").Run("Hello, World!")
// rem: ", World!"
// ext: "Hello"
```

### TagNoCase

Matches a sequence of characters, ignoring case. Returns the original casing from input.

```go
chomp.TagNoCase("hello").Run("HELLO, World!")
// rem: ", World!"
// ext: "HELLO"
```

### Until

Scans until the first occurrence of the given string. Everything before is matched.

```go
chomp.Until("World").Run("Hello, World!")
// rem: "World!"
// ext: "Hello, "
```

### TakeUntil1

Like `Until`, but requires at least one character to be matched before the delimiter.

```go
chomp.TakeUntil1(",").Run("Hello, World!")
// rem: ", World!"
// ext: "Hello"
```

---

## Sequence Combinators

### Pair

Matches two combinators in sequence. Both must match.

```go
chomp.Pair(chomp.Tag("Hello,"), chomp.Tag(" World")).Run("Hello, World!")
// rem: "!"
// ext: ["Hello,", " World"]
```

### SepPair

Matches two combinators separated by a third (separator is discarded).

```go
chomp.SepPair(chomp.Tag("Hello"), chomp.Tag(", "), chomp.Tag("World")).Run("Hello, World!")
// rem: "!"
// ext: ["Hello", "World"]
```

### All

Matches all combinators in sequence. All must match in order.

```go
chomp.All(chomp.Tag("Hello"), chomp.Until("W"), chomp.Tag("World!")).Run("Hello, World!")
// rem: ""
// ext: ["Hello", ", ", "World!"]
```

### First

Tries each combinator in order, returning the first successful match.

```go
chomp.First(chomp.Tag("Good Morning"), chomp.Tag("Hello")).Run("Hello, World!")
// rem: ", World!"
// ext: "Hello"
```

### Repeat

Matches a combinator exactly `n` times.

```go
chomp.Repeat(chomp.Parentheses(), 2).Run("(Hello)(World)(!)")
// rem: "(!)"
// ext: ["Hello", "World"]
```

### RepeatRange

Matches a combinator between `n` (minimum) and `m` (maximum) times.

```go
chomp.RepeatRange(chomp.OneOf("Helo"), 1, 8).Run("Hello, World!")
// rem: ", World!"
// ext: ["H", "e", "l", "l", "o"]
```

### Many

Matches a combinator one or more times (greedy). Equivalent to `ManyN(c, 1)`.

```go
chomp.Many(chomp.OneOf("Ho")).Run("Hello, World!")
// rem: "ello, World!"
// ext: ["H"]
```

### ManyN

Matches a combinator at least `n` times (greedy).

```go
chomp.ManyN(chomp.OneOf("Helo"), 0).Run("Hello, World!")
// rem: ", World!"
// ext: ["H", "e", "l", "l", "o"]
```

### Delimited

Matches content between left and right delimiters (delimiters are discarded).

```go
chomp.Delimited(chomp.Tag("'"), chomp.Until("'"), chomp.Tag("'")).Run("'Hello, World!'")
// rem: ""
// ext: "Hello, World!"
```

### Prefixed

Matches a prefix (discarded) followed by content.

```go
chomp.Prefixed(chomp.Tag("Hello"), chomp.Tag(`"`)).Run(`"Hello, World!"`)
// rem: `, World!"`
// ext: "Hello"
```

### Suffixed

Matches content followed by a suffix (discarded).

```go
chomp.Suffixed(chomp.Tag("Hello"), chomp.Tag(", ")).Run("Hello, World!")
// rem: "World!"
// ext: "Hello"
```

### QuoteDouble

Matches text surrounded by double quotes.

```go
chomp.QuoteDouble().Run(`"Hello, World!"`)
// rem: ""
// ext: "Hello, World!"
```

### QuoteSingle

Matches text surrounded by single quotes.

```go
chomp.QuoteSingle().Run("'Hello, World!'")
// rem: ""
// ext: "Hello, World!"
```

### BracketSquare

Matches text surrounded by square brackets.

```go
chomp.BracketSquare().Run("[Hello, World!]")
// rem: ""
// ext: "Hello, World!"
```

### Parentheses

Matches text surrounded by parentheses.

```go
chomp.Parentheses().Run("(Hello, World!)")
// rem: ""
// ext: "Hello, World!"
```

### BracketAngled

Matches text surrounded by angled brackets.

```go
chomp.BracketAngled().Run("<Hello, World!>")
// rem: ""
// ext: "Hello, World!"
```

### SeparatedList

Matches elements separated by a delimiter. At least one element must match. The separator is discarded.

```go
chomp.SeparatedList(chomp.Alpha(), chomp.Tag(",")).Run("apple,banana,cherry,")
// rem: ","
// ext: ["apple", "banana", "cherry"]
```

### SeparatedList0

Matches elements separated by a delimiter. Zero or more elements may match. The separator is discarded.

```go
chomp.SeparatedList0(chomp.Alpha(), chomp.Tag(",")).Run("123")
// rem: "123"
// ext: []
```

### ManyTill

Matches elements repeatedly until a terminator is found. The terminator is consumed but not included in the result. At least one element must match.

```go
chomp.ManyTill(chomp.AnyChar(), chomp.Tag("END")).Run("abcEND")
// rem: ""
// ext: ["a", "b", "c"]
```

### ManyTill0

Matches elements repeatedly until a terminator is found. The terminator is consumed but not included in the result. Zero or more elements may match.

```go
chomp.ManyTill0(chomp.AnyChar(), chomp.Tag("END")).Run("END")
// rem: ""
// ext: []
```

### FoldMany

Matches a combinator repeatedly and accumulates results using a reducer function. At least one match is required.

```go
chomp.FoldMany(chomp.OneOf("123"), 0, func(acc int, val string) int {
    return acc + int(val[0]-'0')
}).Run("123abc")
// rem: "abc"
// ext: 6
```

### FoldMany0

Matches a combinator repeatedly and accumulates results using a reducer function. Zero or more matches allowed.

```go
chomp.FoldMany0(chomp.OneOf("123"), 0, func(acc int, val string) int {
    return acc + int(val[0]-'0')
}).Run("abc")
// rem: "abc"
// ext: 0
```

### ManyCount

Counts the number of times a combinator matches without storing results. At least one match is required. Memory efficient for counting.

```go
chomp.ManyCount(chomp.OneOf("abc")).Run("abc123")
// rem: "123"
// ext: 3
```

### ManyCount0

Counts the number of times a combinator matches without storing results. Zero or more matches allowed. Memory efficient for counting.

```go
chomp.ManyCount0(chomp.OneOf("abc")).Run("123")
// rem: "123"
// ext: 0
```

### LengthCount

Parses a length value first, then applies a combinator that exact number of times.

```go
chomp.LengthCount(
    chomp.Map(chomp.OneOf("0123456789"), func(s string) uint {
        return uint(s[0] - '0')
    }),
    chomp.OneOf("abc"),
).Run("3abcdef")
// rem: "def"
// ext: ["a", "b", "c"]
```

### Fill

Matches a combinator exactly `n` times, populating a result slice. All `n` matches must succeed.

```go
chomp.Fill(chomp.OneOf("abc"), 3).Run("abcdef")
// rem: "def"
// ext: ["a", "b", "c"]
```

---

## Control Flow Combinators

### Verify

Validates the parsed result against a predicate function without modifying the output. If the predicate returns false, the combinator fails.

```go
chomp.Verify(chomp.Alpha(), func(s string) bool {
    return len(s) >= 3
}).Run("Hello, World!")
// rem: ", World!"
// ext: "Hello"
```

### Recognize

Returns the consumed input as the output, regardless of the inner parser's result. Useful for capturing complex patterns as text.

```go
chomp.Recognize(chomp.SepPair(chomp.Alpha(), chomp.Tag(", "), chomp.Alpha())).Run("Hello, World!")
// rem: "!"
// ext: "Hello, World"
```

### Consumed

Provides both the raw consumed text and the parsed output. The first element is the raw consumed text, followed by the parsed result.

```go
chomp.Consumed(chomp.SepPair(chomp.Alpha(), chomp.Tag(", "), chomp.Alpha())).Run("Hello, World!")
// rem: "!"
// ext: ["Hello, World", "Hello", "World"]
```

### Eof

Matches only when at the end of input, returning an empty string on success. Prevents partial parsing.

```go
chomp.Eof().Run("")
// rem: ""
// ext: ""

chomp.Pair(chomp.Tag("Hello"), chomp.Eof()).Run("Hello")
// rem: ""
// ext: ["Hello", ""]
```

### AllConsuming

Ensures the entire input is consumed by the inner parser, failing if any text remains unparsed.

```go
chomp.AllConsuming(chomp.Tag("Hello")).Run("Hello")
// rem: ""
// ext: "Hello"

chomp.AllConsuming(chomp.Tag("Hello")).Run("Hello, World!")
// error: all_consuming failed
```

### Rest

Returns all remaining unconsumed input as a string value. Always succeeds, even with empty input.

```go
chomp.Rest().Run("Hello, World!")
// rem: ""
// ext: "Hello, World!"

chomp.Pair(chomp.Tag("Hello"), chomp.Rest()).Run("Hello, World!")
// rem: ""
// ext: ["Hello", ", World!"]
```

### Value

Returns a fixed value upon parser success, discarding the actual parse result. Useful for assigning semantic meaning to parsed tokens.

```go
chomp.Value(chomp.Tag("true"), true).Run("true")
// rem: ""
// ext: true

chomp.Value(chomp.Tag("null"), nil).Run("null")
// rem: ""
// ext: nil
```

### Cond

Conditionally applies a parser based on a boolean flag. If the condition is true, the parser is applied. Otherwise, returns an empty result without consuming input.

```go
chomp.Cond(true, chomp.Tag("Hello")).Run("Hello, World!")
// rem: ", World!"
// ext: "Hello"

chomp.Cond(false, chomp.Tag("Hello")).Run("Hello, World!")
// rem: "Hello, World!"
// ext: ""
```

### Cut

Converts recoverable parsing errors into fatal failures, preventing backtracking past decision points. Improves error messaging by committing to a parsing path. When used with `First`, a `CutError` stops backtracking immediately.

```go
// Without Cut, First would try the second alternative and succeed
// With Cut, once "if" matches, failure on "(" is fatal
chomp.First(
    chomp.Flatten(chomp.All(
        chomp.Tag("if"),
        chomp.Cut(chomp.Tag("(")))),
    chomp.Tag("if x")).Run("if x")
// error: CutError (no backtracking to "if x")
```

### PeekNot

Succeeds when the inner parser fails without consuming input. Implements negative lookahead for validation. Pairs with `Peek` for positive lookahead.

```go
chomp.PeekNot(chomp.Tag("Hello")).Run("World!")
// rem: "World!"
// ext: ""

chomp.PeekNot(chomp.Tag("Hello")).Run("Hello, World!")
// error: peek_not failed
```

---

## Predicate Combinators

### While

Matches one or more characters that satisfy a predicate.

```go
chomp.While(chomp.IsLetter).Run("Hello, World!")
// rem: ", World!"
// ext: "Hello"
```

### WhileN

Matches at least `n` characters that satisfy a predicate. If `n` is 0, becomes optional.

```go
chomp.WhileN(chomp.IsDigit, 2).Run("12345abc")
// rem: "abc"
// ext: "12345"
```

### WhileNM

Matches between `n` and `m` characters that satisfy a predicate.

```go
chomp.WhileNM(chomp.IsLetter, 1, 3).Run("Hel1lo")
// rem: "1lo"
// ext: "Hel"
```

### WhileNot

Matches one or more characters that do NOT satisfy a predicate.

```go
chomp.WhileNot(chomp.IsDigit).Run("Hello123")
// rem: "123"
// ext: "Hello"
```

### WhileNotN

Matches at least `n` characters that do NOT satisfy a predicate.

```go
chomp.WhileNotN(chomp.IsDigit, 1).Run("Hello123")
// rem: "123"
// ext: "Hello"
```

### WhileNotNM

Matches between `n` and `m` characters that do NOT satisfy a predicate.

```go
chomp.WhileNotNM(chomp.IsLetter, 1, 9).Run("20240709 was great")
// rem: "was great"
// ext: "20240709 "
```

---

## Convenience Combinators

These are shorthand combinators built on predicates.

| Combinator          | Equivalent                      | Description                   |
| ------------------- | ------------------------------- | ----------------------------- |
| `Alpha()`           | `While(IsLetter)`               | One or more letters           |
| `Alpha0()`          | `WhileN(IsLetter, 0)`           | Zero or more letters          |
| `Digit()`           | `While(IsDigit)`                | One or more digits            |
| `Digit0()`          | `WhileN(IsDigit, 0)`            | Zero or more digits           |
| `Alphanumeric()`    | `While(IsAlphanumeric)`         | One or more alphanumeric      |
| `Alphanumeric0()`   | `WhileN(IsAlphanumeric, 0)`     | Zero or more alphanumeric     |
| `Space()`           | `While(IsSpace)`                | One or more spaces/tabs       |
| `Space0()`          | `WhileN(IsSpace, 0)`            | Zero or more spaces/tabs      |
| `Multispace()`      | `While(IsMultispace)`           | One or more whitespace        |
| `Multispace0()`     | `WhileN(IsMultispace, 0)`       | Zero or more whitespace       |
| `HexDigit()`        | `While(IsHexDigit)`             | One or more hex digits        |
| `HexDigit0()`       | `WhileN(IsHexDigit, 0)`         | Zero or more hex digits       |
| `OctalDigit()`      | `While(IsOctalDigit)`           | One or more octal digits      |
| `OctalDigit0()`     | `WhileN(IsOctalDigit, 0)`       | Zero or more octal digits     |
| `BinaryDigit()`     | `While(IsBinaryDigit)`          | One or more binary digits     |
| `BinaryDigit0()`    | `WhileN(IsBinaryDigit, 0)`      | Zero or more binary digits    |
| `Newline()`         | -                               | Matches `\n`                  |
| `Tab()`             | -                               | Matches `\t`                  |
| `NotLineEnding()`   | `WhileNot(IsLineEnding)`        | Characters until line ending  |
| `AnyDigit()`        | `Satisfy(IsDigit.Match)`        | Single decimal digit          |
| `AnyLetter()`       | `Satisfy(IsLetter.Match)`       | Single letter                 |
| `AnyAlphanumeric()` | `Satisfy(IsAlphanumeric.Match)` | Single alphanumeric character |
| `AnyHexDigit()`     | `Satisfy(IsHexDigit.Match)`     | Single hexadecimal digit      |
| `AnyOctalDigit()`   | `Satisfy(IsOctalDigit.Match)`   | Single octal digit            |
| `AnyBinaryDigit()`  | `Satisfy(IsBinaryDigit.Match)`  | Single binary digit           |

---

## Modifier Combinators

### Map

Transforms the result of a combinator to any other type.

```go
chomp.Map(chomp.While(chomp.IsDigit), func(in string) int {
    n, _ := strconv.Atoi(in)
    return n
}).Run("123abc")
// rem: "abc"
// ext: 123
```

### S

Wraps a string result in a slice. Useful for chaining combinators with different return types.

```go
chomp.S(chomp.Until(",")).Run("Hello, World!")
// rem: ", World!"
// ext: ["Hello"]
```

### I

Extracts a single string from a slice result at index `i`.

```go
chomp.I(chomp.SepPair(chomp.Tag("Hello"), chomp.Tag(", "), chomp.Tag("World")), 1).Run("Hello, World!")
// rem: "!"
// ext: "World"
```

### Peek

Applies a combinator without consuming input (lookahead).

```go
chomp.Peek(chomp.Tag("Hello")).Run("Hello, World!")
// rem: "Hello, World!"
// ext: "Hello"
```

### Opt

Makes a combinator optional. On failure, returns empty without error.

```go
chomp.Opt(chomp.Tag("Hey")).Run("Hello, World!")
// rem: "Hello, World!"
// ext: ""
```

### Flatten

Joins all extracted values from a combinator into a single string.

```go
chomp.Flatten(chomp.Many(chomp.Parentheses())).Run("(H)(el)(lo)")
// rem: ""
// ext: "Hello"
```

### Escaped

Parses strings containing escape sequences, preserving them as-is.

```go
chomp.Escaped(chomp.While(chomp.IsLetter), '\\', chomp.OneOf(`"n\`)).Run(`Hello\"World`)
// rem: ""
// ext: `Hello\"World`
```

### EscapedTransform

Parses and transforms escape sequences to their actual values.

```go
transform := func(s chomp.State) (chomp.State, string, error) {
    switch s.Rest()[0] {
    case 'n':
        return s.Advance(1), "\n", nil
    case '"':
        return s.Advance(1), "\"", nil
    }
    return s, "", errors.New("invalid escape")
}
chomp.EscapedTransform(chomp.While(chomp.IsLetter), '\\', transform).Run(`Hello\nWorld`)
// rem: ""
// ext: "Hello\nWorld"  // actual newline character
```

---

## Parsers

### Crlf

Matches a strict CRLF (`\r\n`) line ending only. Does not match a bare LF or bare CR.

```go
chomp.Crlf().Run("\r\nHello")
// rem: "Hello"
// ext: "\r\n"
```

### LineEnding

Matches either `\n` (LF) or `\r\n` (CRLF). Does not match a bare CR; see `Eol` if a bare CR should be treated as a legacy-Mac line terminator.

```go
chomp.LineEnding().Run("\nHello")
// rem: "Hello"
// ext: "\n"
```

### Eol

Scans and returns text before any line ending. The line ending is discarded. Unlike `LineEnding`, a bare CR is also consumed here, treated as a legacy-Mac line terminator.

```go
chomp.Eol().Run("Hello, World!\nNext line")
// rem: "Next line"
// ext: "Hello, World!"
```

---

## Predicates

Predicates are used with `While`, `WhileN`, `WhileNM`, `WhileNot`, `WhileNotN`, and `WhileNotNM`.

| Predicate        | Description                                            |
| ---------------- | ------------------------------------------------------ |
| `IsDigit`        | Decimal digits (0-9) and Unicode Nd category           |
| `IsLetter`       | ASCII letters (a-z, A-Z) and Unicode letter categories |
| `IsAlphanumeric` | Combination of `IsDigit` and `IsLetter`                |
| `IsLineEnding`   | Line ending characters (`\n`, `\r`)                    |
| `IsSpace`        | Space (` `) or tab (`\t`)                              |
| `IsMultispace`   | Space, tab, newline, or carriage return                |
| `IsHexDigit`     | Hexadecimal digits (0-9, a-f, A-F)                     |
| `IsOctalDigit`   | Octal digits (0-7)                                     |
| `IsBinaryDigit`  | Binary digits (0-1)                                    |

---

## Error Reporting

Every combinator failure is a `chomp.CombinatorParseError`, reachable via `errors.As` through any wrapping error. Rendering is layered:

- **`Error()`** - a single line, safe for any log pipeline (`\n`-delimited collectors treat a newline as a record boundary):

```go
_, _, err := chomp.All(
    chomp.Until("."),
    chomp.Tag("."),
    chomp.Digit(),
    chomp.Tag("."),
    chomp.Digit(),
).Run("version = 1.4.x")
fmt.Println(err)
// chomp: parse error at line 1, column 15 (offset 14): expected digit
```

- **`Snippet()`** - opt-in, reached via `errors.As`, for a human at a terminal:

```go
var pe chomp.CombinatorParseError
errors.As(err, &pe)
fmt.Println(err.Error() + "\n" + pe.Snippet())
// chomp: parse error at line 1, column 15 (offset 14): expected digit
//   |
// 1 | version = 1.4.x
//   |               ^ expected digit
```

- **`LogValue()`** - implements `slog.LogValuer`, so a structured logger gets `offset`/`line`/`column`/`expected` (and `context`, once a failure occurred inside a labelled grammar rule) as fields, not a string to parse:

```go
logger.Error("failed to parse manifest", "err", err)
// {"level":"ERROR","msg":"failed to parse manifest","err":{"offset":14,"line":1,"column":15,"expected":"digit"}}
```

### Label

Attaches a grammar-level name to any failure beneath a combinator, so an error speaks in the grammar's own terms rather than internal combinator names. Nested `Label`s chain outermost-first.

```go
chomp.Label("manifest", chomp.Label("version", chomp.All(
    chomp.Until("."),
    chomp.Tag("."),
    chomp.Digit(),
    chomp.Tag("."),
    chomp.Digit(),
))).Run("version = 1.4.x")
fmt.Println(err)
// chomp: parse error at line 1, column 15 (offset 14): expected digit while parsing manifest > version
```

`Snippet()` and `LogValue()`'s `context` field render the same chain.
