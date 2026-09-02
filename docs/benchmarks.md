# Benchmarks

Benchmarks are run on an Apple M4 Pro (darwin/arm64) using Go's built-in testing framework.

```
go test -bench=. -benchmem ./...
```

## Results

### Basic Combinators

| Benchmark         |  ns/op |  MB/s | B/op | allocs/op |
| ----------------- | -----: | ----: | ---: | --------: |
| Tag/Ascii         |   4.24 | 10138 |    0 |         0 |
| Tag/Unicode       |   4.29 | 13999 |    0 |         0 |
| TagNoCase/Ascii   |  47.28 |   910 |    0 |         0 |
| Char/Ascii        |   3.34 | 12876 |    0 |         0 |
| Char/Unicode      |   3.94 | 15239 |    0 |         0 |
| AnyChar/Ascii     |   3.36 | 12813 |    0 |         0 |
| AnyChar/Unicode   |   3.73 | 16102 |    0 |         0 |
| Take/Ascii        |  17.58 |  2446 |    0 |         0 |
| Take/Unicode      |  11.88 |  5049 |    0 |         0 |
| Until/Ascii       |   7.70 |  5585 |    0 |         0 |
| Until/Unicode     |  14.16 |  4237 |    0 |         0 |
| Any/Small/Ascii   |  10.31 |  4171 |    0 |         0 |
| Any/Large/Ascii   |  69.70 |   617 |    0 |         0 |
| Any/Small/Unicode |  30.35 |  1977 |    0 |         0 |
| Any/Large/Unicode |  59.74 |  1004 |    0 |         0 |
| Not/Small/Ascii   |  50.72 |   848 |    0 |         0 |
| Not/Large/Ascii   | 112.60 |   382 |    0 |         0 |
| Not/Small/Unicode |  99.12 |   605 |    0 |         0 |
| Not/Large/Unicode | 111.80 |   537 |    0 |         0 |
| OneOf/Ascii       |  13.34 |  3224 |    0 |         0 |
| OneOf/Unicode     |   5.73 | 10480 |    0 |         0 |
| NoneOf/Ascii      |   5.40 |  7967 |    0 |         0 |
| NoneOf/Unicode    |   7.96 |  7538 |    0 |         0 |

### Predicate Combinators

| Benchmark              |  ns/op |  MB/s | B/op | allocs/op |
| ---------------------- | -----: | ----: | ---: | --------: |
| While/Digit            |  31.24 |  1697 |    0 |         0 |
| While/Letter/Ascii     |  12.37 |  3476 |    0 |         0 |
| While/Letter/Unicode   | 198.20 |   303 |    0 |         0 |
| While/Alphanumeric     |  57.00 |   281 |    0 |         0 |
| While/Space            |  20.77 |  2311 |    0 |         0 |
| WhileNot/Digit/Ascii   | 135.50 |   339 |    0 |         0 |
| WhileNot/Digit/Unicode | 184.10 |   342 |    0 |         0 |
| Satisfy/Ascii          |   3.93 | 10943 |    0 |         0 |
| Satisfy/Unicode        |   4.23 | 14199 |    0 |         0 |

### Sequence Combinators

| Benchmark             |  ns/op | MB/s | B/op | allocs/op |
| --------------------- | -----: | ---: | ---: | --------: |
| Pair/Ascii            |  43.76 |  983 |   48 |         2 |
| Pair/Unicode          |  44.32 | 1354 |   48 |         2 |
| Delimited/Parentheses |  12.32 | 1705 |    0 |         0 |
| Delimited/Quotes      |  12.36 | 1699 |    0 |         0 |
| SepPair               |  65.35 |  214 |   48 |         2 |
| All/ThreeTags         |  72.86 |  590 |  112 |         3 |
| All/FiveTags          | 141.00 |  305 |  240 |         4 |

### Modifier Combinators

| Benchmark    |   ns/op |  MB/s |  B/op | allocs/op |
| ------------ | ------: | ----: | ----: | --------: |
| Opt/Match    |    5.55 |  7742 |     0 |         0 |
| Opt/NoMatch  |   76.27 |   603 |   117 |         3 |
| Map          |   17.15 |  2799 |     0 |         0 |
| Many/Small   |  193.60 |    72 |   356 |         7 |
| Many/Medium  |  762.40 |   136 |  2277 |        10 |
| Many/Large   | 5212.00 |   193 | 18927 |        13 |
| Peek/Ascii   |    9.47 |  4541 |     0 |         0 |
| Peek/Unicode |   12.66 |  4739 |     0 |         0 |
| Flatten      |   97.14 |   443 |   128 |         4 |

`Opt/NoMatch` and `Many/*` construct and discard at least one `CombinatorParseError` internally (the failed match `Opt` swallows, or the zero-width guard that stops a `Many` loop) - non-zero, but bounded and independent of input size.

### Control Flow Combinators

| Benchmark         |  ns/op |  MB/s | B/op | allocs/op |
| ----------------- | -----: | ----: | ---: | --------: |
| First/FirstMatch  |   5.45 |  7893 |    0 |         0 |
| First/LastMatch   | 171.80 |    82 |  240 |         6 |
| Verify/Pass       |  20.65 |  2325 |    0 |         0 |
| Recognize/Ascii   |  49.53 |   868 |   48 |         2 |
| Recognize/Unicode |  53.97 |  1112 |   48 |         2 |
| Consumed          |  84.68 |   508 |  112 |         4 |
| Eof               |   2.09 |     - |    0 |         0 |
| AllConsuming      |   5.89 |  3226 |    0 |         0 |
| Rest/Ascii        |   2.64 | 16300 |    0 |         0 |
| Rest/Unicode      |   2.62 | 22937 |    0 |         0 |
| Value             |   5.44 |  8633 |    0 |         0 |
| Cond/True         |   5.07 |  8486 |    0 |         0 |
| Cond/False        |   1.74 | 24719 |    0 |         0 |
| Cut               |   5.28 |  8148 |    0 |         0 |

`First/LastMatch` discards two failed alternatives (`the`, `quick`) before matching `fox`, each constructing its own `CombinatorParseError`.

### Parser Combinators

| Benchmark   |  ns/op |  MB/s | B/op | allocs/op |
| ----------- | -----: | ----: | ---: | --------: |
| Crlf        |   2.49 | 18112 |    0 |         0 |
| Eol/Ascii   |  59.51 |   739 |    0 |         0 |
| Eol/Unicode |  58.25 |  1047 |    0 |         0 |

### Scaling Benchmarks

| Benchmark           |    ns/op |  MB/s | B/op | allocs/op |
| ------------------- | -------: | ----: | ---: | --------: |
| UntilScaling/Small  |     7.17 |  1255 |    0 |         0 |
| UntilScaling/Medium |     8.82 | 11798 |    0 |         0 |
| UntilScaling/Large  |   116.00 | 86205 |    0 |         0 |
| WhileScaling/Small  |    11.29 |   531 |    0 |         0 |
| WhileScaling/Medium |   290.30 |   355 |    0 |         0 |
| WhileScaling/Large  | 26046.00 |   384 |    0 |         0 |
