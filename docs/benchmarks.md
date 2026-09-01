# Benchmarks

Benchmarks are run on an Apple M4 Pro (darwin/arm64) using Go's built-in testing framework.

```
go test -bench=. -benchmem ./...
```

## Results

### Basic Combinators

| Benchmark         |  ns/op |  MB/s | B/op | allocs/op |
| ----------------- | -----: | ----: | ---: | --------: |
| Tag/Ascii         |   4.47 |  9615 |    0 |         0 |
| Tag/Unicode       |   4.47 | 13424 |    0 |         0 |
| TagNoCase/Ascii   |  49.98 |   860 |    0 |         0 |
| Char/Ascii        |   3.47 | 12400 |    0 |         0 |
| Char/Unicode      |   3.87 | 15524 |    0 |         0 |
| AnyChar/Ascii     |   3.35 | 12828 |    0 |         0 |
| AnyChar/Unicode   |   3.75 | 15984 |    0 |         0 |
| Take/Ascii        |  17.89 |  2403 |    0 |         0 |
| Take/Unicode      |  12.37 |  4851 |    0 |         0 |
| Until/Ascii       |   7.74 |  5558 |    0 |         0 |
| Until/Unicode     |  13.91 |  4313 |    0 |         0 |
| Any/Small/Ascii   |  10.19 |  4220 |    0 |         0 |
| Any/Large/Ascii   |  77.73 |   553 |    0 |         0 |
| Any/Small/Unicode |  30.71 |  1954 |    0 |         0 |
| Any/Large/Unicode |  64.36 |   932 |    0 |         0 |
| Not/Small/Ascii   |  50.89 |   845 |    0 |         0 |
| Not/Large/Ascii   | 118.90 |   362 |    0 |         0 |
| Not/Small/Unicode |  97.94 |   613 |    0 |         0 |
| Not/Large/Unicode | 115.30 |   521 |    0 |         0 |
| OneOf/Ascii       |  13.52 |  3181 |    0 |         0 |
| OneOf/Unicode     |   5.74 | 10456 |    0 |         0 |
| NoneOf/Ascii      |   5.55 |  7743 |    0 |         0 |
| NoneOf/Unicode    |   7.86 |  7637 |    0 |         0 |

### Predicate Combinators

| Benchmark              |  ns/op |  MB/s | B/op | allocs/op |
| ---------------------- | -----: | ----: | ---: | --------: |
| While/Digit            |  31.78 |  1668 |    0 |         0 |
| While/Letter/Ascii     |  12.68 |  3392 |    0 |         0 |
| While/Letter/Unicode   | 200.50 |   299 |    0 |         0 |
| While/Alphanumeric     |  56.15 |   285 |    0 |         0 |
| While/Space            |  21.23 |  2261 |    0 |         0 |
| WhileNot/Digit/Ascii   | 128.50 |   358 |    0 |         0 |
| WhileNot/Digit/Unicode | 163.50 |   385 |    0 |         0 |
| Satisfy/Ascii          |   3.75 | 11474 |    0 |         0 |
| Satisfy/Unicode        |   5.18 | 11588 |    0 |         0 |

### Sequence Combinators

| Benchmark             |  ns/op | MB/s | B/op | allocs/op |
| --------------------- | -----: | ---: | ---: | --------: |
| Pair/Ascii            |  44.89 |  958 |   48 |         2 |
| Pair/Unicode          |  44.51 | 1348 |   48 |         2 |
| Delimited/Parentheses |  12.50 | 1680 |    0 |         0 |
| Delimited/Quotes      |  12.31 | 1705 |    0 |         0 |
| SepPair               |  66.83 |  209 |   48 |         2 |
| All/ThreeTags         |  73.78 |  583 |  112 |         3 |
| All/FiveTags          | 114.50 |  375 |  240 |         4 |

### Modifier Combinators

| Benchmark    |   ns/op |  MB/s |  B/op | allocs/op |
| ------------ | ------: | ----: | ----: | --------: |
| Opt/Match    |    5.24 |  8209 |     0 |         0 |
| Opt/NoMatch  |   22.14 |  2078 |    64 |         1 |
| Map          |   17.72 |  2708 |     0 |         0 |
| Many/Small   |  145.70 |    96 |   304 |         5 |
| Many/Medium  |  719.10 |   145 |  2224 |         8 |
| Many/Large   | 5234.00 |   192 | 18864 |        11 |
| Peek/Ascii   |    9.49 |  4529 |     0 |         0 |
| Peek/Unicode |   12.63 |  4749 |     0 |         0 |
| Flatten      |   98.29 |   438 |   128 |         4 |

### Control Flow Combinators

| Benchmark         |  ns/op |  MB/s | B/op | allocs/op |
| ----------------- | -----: | ----: | ---: | --------: |
| First/FirstMatch  |   5.49 |  7832 |    0 |         0 |
| First/LastMatch   | 129.90 |   108 |  160 |         4 |
| Verify/Pass       |  17.21 |  2789 |    0 |         0 |
| Recognize/Ascii   |  50.21 |   856 |   48 |         2 |
| Recognize/Unicode |  54.28 |  1105 |   48 |         2 |
| Consumed          |  88.09 |   488 |  112 |         4 |
| Eof               |   2.07 |     - |    0 |         0 |
| AllConsuming      |   5.99 |  3170 |    0 |         0 |
| Rest/Ascii        |   2.65 | 16244 |    0 |         0 |
| Rest/Unicode      |   2.65 | 22668 |    0 |         0 |
| Value             |   5.26 |  8942 |    0 |         0 |
| Cond/True         |   4.99 |  8611 |    0 |         0 |
| Cond/False        |   1.75 | 24600 |    0 |         0 |
| Cut               |   5.25 |  8184 |    0 |         0 |

### Parser Combinators

| Benchmark   |  ns/op |  MB/s | B/op | allocs/op |
| ----------- | -----: | ----: | ---: | --------: |
| Crlf        |   2.50 | 17977 |    0 |         0 |
| Eol/Ascii   |  61.52 |   715 |    0 |         0 |
| Eol/Unicode |  60.73 |  1004 |    0 |         0 |

### Scaling Benchmarks

| Benchmark           |    ns/op |  MB/s | B/op | allocs/op |
| ------------------- | -------: | ----: | ---: | --------: |
| UntilScaling/Small  |     7.24 |  1243 |    0 |         0 |
| UntilScaling/Medium |     8.70 | 11948 |    0 |         0 |
| UntilScaling/Large  |   118.00 | 84760 |    0 |         0 |
| WhileScaling/Small  |    11.40 |   526 |    0 |         0 |
| WhileScaling/Medium |   277.30 |   371 |    0 |         0 |
| WhileScaling/Large  | 27089.00 |   369 |    0 |         0 |

### Real-World Patterns

| Benchmark     |  ns/op | MB/s | B/op | allocs/op |
| ------------- | -----: | ---: | ---: | --------: |
| KeyValuePair  |  57.18 |  560 |   48 |         2 |
| GitDiffHeader | 109.20 |  174 |  176 |         5 |
| CSVField      |  27.17 | 1104 |    0 |         0 |
