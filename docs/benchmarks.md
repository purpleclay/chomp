# Benchmarks

Benchmarks are run on an Apple M4 Pro (darwin/arm64) using Go's built-in testing framework.

```
go test -bench=. -benchmem ./...
```

## Results

### Basic Combinators

| Benchmark            |  ns/op |  MB/s | B/op | allocs/op |
| -------------------- | -----: | ----: | ---: | --------: |
| Tag/Ascii            |   4.38 |  9824 |    0 |         0 |
| Tag/Unicode          |   4.40 | 13630 |    0 |         0 |
| TagNoCase/Ascii      |  49.43 |   870 |    0 |         0 |
| Char/Ascii           |   3.35 | 12820 |    0 |         0 |
| Char/Unicode         |   3.89 | 15429 |    0 |         0 |
| AnyChar/Ascii        |   3.31 | 12988 |    0 |         0 |
| AnyChar/Unicode      |   3.76 | 15944 |    0 |         0 |
| Take/Ascii           |  17.86 |  2407 |    0 |         0 |
| Take/Unicode         |  12.35 |  4858 |    0 |         0 |
| Until/Ascii          |   8.19 |  5250 |    0 |         0 |
| Until/Unicode        |  13.78 |  4354 |    0 |         0 |
| IsA/Small/Ascii      |  10.57 |  4067 |    0 |         0 |
| IsA/Large/Ascii      |  75.41 |   570 |    0 |         0 |
| IsA/Small/Unicode    |  28.91 |  2076 |    0 |         0 |
| IsA/Large/Unicode    |  59.17 |  1014 |    0 |         0 |
| IsNot/Small/Ascii    |  55.16 |   780 |    0 |         0 |
| IsNot/Large/Ascii    | 117.40 |   366 |    0 |         0 |
| IsNot/Small/Unicode  | 103.90 |   578 |    0 |         0 |
| IsNot/Large/Unicode  | 115.20 |   521 |    0 |         0 |
| OneOf/Ascii          |  13.78 |  3120 |    0 |         0 |
| OneOf/Unicode        |   5.82 | 10313 |    0 |         0 |
| NoneOf/Ascii         |   5.45 |  7889 |    0 |         0 |
| NoneOf/Unicode       |   8.19 |  7330 |    0 |         0 |

### Predicate Combinators

| Benchmark              |  ns/op |  MB/s | B/op | allocs/op |
| ---------------------- | -----: | ----: | ---: | --------: |
| While/Digit            |  31.28 |  1694 |    0 |         0 |
| While/Letter/Ascii     |  13.04 |  3297 |    0 |         0 |
| While/Letter/Unicode   | 204.40 |   294 |    0 |         0 |
| While/Alphanumeric     |  56.32 |   284 |    0 |         0 |
| While/Space            |  21.68 |  2214 |    0 |         0 |
| WhileNot/Digit/Ascii   | 131.50 |   350 |    0 |         0 |
| WhileNot/Digit/Unicode | 186.50 |   338 |    0 |         0 |
| Satisfy/Ascii          |   3.73 | 11542 |    0 |         0 |
| Satisfy/Unicode        |   4.46 | 13444 |    0 |         0 |

### Sequence Combinators

| Benchmark             |  ns/op | MB/s | B/op | allocs/op |
| ---------------------- | -----: | ---: | ---: | --------: |
| Pair/Ascii             |   8.71 | 4936 |    0 |         0 |
| Pair/Unicode           |   8.69 | 6907 |    0 |         0 |
| Delimited/Parentheses  |  12.33 | 1703 |    0 |         0 |
| Delimited/Quotes       |  12.43 | 1690 |    0 |         0 |
| SepPair                |  29.36 |  477 |    0 |         0 |
| All/ThreeTags          |  67.52 |  637 |  112 |         3 |
| All/FiveTags           | 102.70 |  419 |  240 |         4 |

### Modifier Combinators

| Benchmark    |   ns/op |  MB/s |  B/op | allocs/op |
| ------------ | ------: | ----: | ----: | --------: |
| Opt/Match    |    5.27 |  8156 |     0 |         0 |
| Opt/NoMatch  |   75.80 |   607 |   117 |         3 |
| Map          |   17.44 |  2752 |     0 |         0 |
| Many/Small   |  193.40 |    72 |   356 |         7 |
| Many/Medium  |  717.00 |   145 |  2277 |        10 |
| Many/Large   | 4868.00 |   206 | 18927 |        13 |
| Peek/Ascii   |    9.50 |  4526 |     0 |         0 |
| Peek/Unicode |   13.17 |  4555 |     0 |         0 |
| Flatten      |   91.88 |   468 |   128 |         4 |

### Control Flow Combinators

| Benchmark         |  ns/op |  MB/s | B/op | allocs/op |
| ----------------- | -----: | ----: | ---: | --------: |
| First/FirstMatch  |   5.73 |  7511 |    0 |         0 |
| First/LastMatch   | 167.40 |    84 |  240 |         6 |
| Verify/Pass       |  18.47 |  2599 |    0 |         0 |
| Recognize/Ascii   |  14.50 |  2965 |    0 |         0 |
| Recognize/Unicode |  19.18 |  3128 |    0 |         0 |
| Consumed          |  17.79 |  2417 |    0 |         0 |
| Eof               |   2.09 |     - |    0 |         0 |
| AllConsuming      |   7.41 |  2564 |    0 |         0 |
| Rest/Ascii        |   2.67 | 16084 |    0 |         0 |
| Rest/Unicode      |   2.65 | 22613 |    0 |         0 |
| Value             |   5.48 |  8573 |    0 |         0 |
| Cond/True         |   5.75 |  7476 |    0 |         0 |
| Cond/False        |   1.75 | 24647 |    0 |         0 |
| Cut               |   5.93 |  7248 |    0 |         0 |

### Parser Combinators

| Benchmark   |  ns/op |  MB/s | B/op | allocs/op |
| ----------- | -----: | ----: | ---: | --------: |
| Crlf        |   2.49 | 18042 |    0 |         0 |
| Eol/Ascii   |  61.45 |   716 |    0 |         0 |
| Eol/Unicode |  58.89 |  1036 |    0 |         0 |

### Scaling Benchmarks

| Benchmark           |    ns/op |  MB/s | B/op | allocs/op |
| ------------------- | -------: | ----: | ---: | --------: |
| UntilScaling/Small  |     7.24 |  1244 |    0 |         0 |
| UntilScaling/Medium |     9.33 | 11146 |    0 |         0 |
| UntilScaling/Large  |   119.10 | 84008 |    0 |         0 |
| WhileScaling/Small  |    11.31 |   531 |    0 |         0 |
| WhileScaling/Medium |   294.30 |   350 |    0 |         0 |
| WhileScaling/Large  | 27141.00 |   369 |    0 |         0 |
