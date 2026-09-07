# Benchmarks

Benchmarks are run on an Apple M4 Pro (darwin/arm64) using Go's built-in testing framework.

```
go test -bench=. -benchmem ./...
```

## Results

### Basic Combinators

| Benchmark            |  ns/op |  MB/s | B/op | allocs/op |
| --------------------- | -----: | ----: | ---: | --------: |
| Tag/Ascii             |   4.33 |  9937 |    0 |         0 |
| Tag/Unicode           |   4.24 | 14153 |    0 |         0 |
| TagNoCase/Ascii       |  46.74 |   920 |    0 |         0 |
| Char/Ascii            |   3.36 | 12782 |    0 |         0 |
| Char/Unicode          |   3.87 | 15525 |    0 |         0 |
| AnyChar/Ascii         |   3.34 | 12877 |    0 |         0 |
| AnyChar/Unicode       |   3.74 | 16038 |    0 |         0 |
| Take/Ascii            |  17.88 |  2405 |    0 |         0 |
| Take/Unicode          |  12.02 |  4991 |    0 |         0 |
| Until/Ascii           |   8.24 |  5222 |    0 |         0 |
| Until/Unicode         |  13.91 |  4314 |    0 |         0 |
| IsA/Small/Ascii       |  11.13 |  3863 |    0 |         0 |
| IsA/Large/Ascii       |  69.35 |   620 |    0 |         0 |
| IsA/Small/Unicode     |  30.33 |  1979 |    0 |         0 |
| IsA/Large/Unicode     |  62.65 |   958 |    0 |         0 |
| IsNot/Small/Ascii     |  55.09 |   781 |    0 |         0 |
| IsNot/Large/Ascii     | 115.40 |   373 |    0 |         0 |
| IsNot/Small/Unicode   |  99.37 |   604 |    0 |         0 |
| IsNot/Large/Unicode   | 115.60 |   519 |    0 |         0 |
| OneOf/Ascii           |  14.31 |  3006 |    0 |         0 |
| OneOf/Unicode         |   5.74 | 10462 |    0 |         0 |
| NoneOf/Ascii          |   5.47 |  7865 |    0 |         0 |
| NoneOf/Unicode        |   7.96 |  7537 |    0 |         0 |

### Predicate Combinators

| Benchmark              |  ns/op |  MB/s | B/op | allocs/op |
| ---------------------- | -----: | ----: | ---: | --------: |
| While/Digit            |  36.58 |  1449 |    0 |         0 |
| While/Letter/Ascii     |  14.55 |  2956 |    0 |         0 |
| While/Letter/Unicode   | 219.40 |   273 |    0 |         0 |
| While/Alphanumeric     |  57.46 |   278 |    0 |         0 |
| While/Space            |  22.36 |  2147 |    0 |         0 |
| WhileNot/Digit/Ascii   | 170.40 |   270 |    0 |         0 |
| WhileNot/Digit/Unicode | 185.80 |   339 |    0 |         0 |
| Satisfy/Ascii          |   3.96 | 10867 |    0 |         0 |
| Satisfy/Unicode        |   4.43 | 13532 |    0 |         0 |

### Sequence Combinators

| Benchmark             |  ns/op | MB/s | B/op | allocs/op |
| ---------------------- | -----: | ---: | ---: | --------: |
| Pair/Ascii             |   8.96 | 4801 |    0 |         0 |
| Pair/Unicode           |   9.18 | 6533 |    0 |         0 |
| Delimited/Parentheses  |  12.84 | 1635 |    0 |         0 |
| Delimited/Quotes       |  12.57 | 1670 |    0 |         0 |
| SepPair                |  37.17 |  377 |    0 |         0 |
| All/ThreeTags          |  68.36 |  629 |  112 |         3 |
| All/FiveTags           | 106.50 |  404 |  240 |         4 |

### Modifier Combinators

| Benchmark    |   ns/op |  MB/s |  B/op | allocs/op |
| ------------ | ------: | ----: | ----: | --------: |
| Opt/Match    |    5.19 |  8292 |     0 |         0 |
| Opt/NoMatch  |   80.34 |   573 |   133 |         3 |
| Map          |   20.96 |  2290 |     0 |         0 |
| Many/Small   |  194.80 |    72 |   372 |         7 |
| Many/Medium  |  703.50 |   148 |  2293 |        10 |
| Many/Large   | 4520.00 |   222 | 18943 |        13 |
| Peek/Ascii   |    9.26 |  4644 |     0 |         0 |
| Peek/Unicode |   12.63 |  4751 |     0 |         0 |
| Flatten      |   92.21 |   466 |   128 |         4 |

`Opt/NoMatch` and `Many/*` construct and discard at least one `CombinatorParseError` internally (the failed match `Opt` swallows, or the zero-width guard that stops a `Many` loop) - `CombinatorParseError` grew by one field (`Cause error`, for `MapRes`), which is why these rows carry 16 more bytes than before.

### Control Flow Combinators

| Benchmark         |  ns/op |  MB/s | B/op | allocs/op |
| ----------------- | -----: | ----: | ---: | --------: |
| First/FirstMatch  |   5.58 |  7704 |    0 |         0 |
| First/LastMatch   | 213.70 |    66 |  320 |         8 |
| Verify/Pass       |  21.16 |  2269 |    0 |         0 |
| Recognize/Ascii   |  14.18 |  3033 |    0 |         0 |
| Recognize/Unicode |  17.90 |  3351 |    0 |         0 |
| Consumed          |  17.59 |  2445 |    0 |         0 |
| Eof               |   2.09 |     - |    0 |         0 |
| AllConsuming      |   5.90 |  3222 |    0 |         0 |
| Rest/Ascii        |   2.67 | 16129 |    0 |         0 |
| Rest/Unicode      |   2.63 | 22783 |    0 |         0 |
| Value             |   5.16 |  9110 |    0 |         0 |
| Cond/True         |   5.05 |  8512 |    0 |         0 |
| Cond/False        |   1.74 | 24703 |    0 |         0 |
| Cut               |   5.15 |  8358 |    0 |         0 |

`First/LastMatch` discards two failed alternatives (`the`, `quick`) before matching `fox`, each constructing its own `CombinatorParseError` - now also collected into an `AlternativesError` rather than discarded, which is why this row's allocations increased over the previous count.

### Parser Combinators

| Benchmark   |  ns/op |  MB/s | B/op | allocs/op |
| ----------- | -----: | ----: | ---: | --------: |
| Crlf        |   2.49 | 18087 |    0 |         0 |
| Eol/Ascii   |  61.46 |   716 |    0 |         0 |
| Eol/Unicode |  59.99 |  1017 |    0 |         0 |

### Scaling Benchmarks

| Benchmark           |    ns/op |  MB/s | B/op | allocs/op |
| ------------------- | -------: | ----: | ---: | --------: |
| UntilScaling/Small  |     7.27 |  1238 |    0 |         0 |
| UntilScaling/Medium |     8.68 | 11983 |    0 |         0 |
| UntilScaling/Large  |   117.10 | 85461 |    0 |         0 |
| WhileScaling/Small  |    14.44 |   415 |    0 |         0 |
| WhileScaling/Medium |   327.10 |   315 |    0 |         0 |
| WhileScaling/Large  | 32818.00 |   305 |    0 |         0 |
