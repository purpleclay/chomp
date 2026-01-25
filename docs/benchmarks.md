# Benchmarks

Benchmarks are run on an Apple M4 Pro (darwin/arm64) using Go's built-in testing framework.

```
go test -bench=. -benchmem ./...
```

## Results

### Basic Combinators

| Benchmark       | ns/op |  MB/s | B/op | allocs/op |
| --------------- | ----: | ----: | ---: | --------: |
| Tag/Ascii       |  1.38 | 31168 |    0 |         0 |
| Tag/Unicode     |  1.42 | 42406 |    0 |         0 |
| TagNoCase/Ascii |  6.71 |  6404 |    0 |         0 |
| Char/Ascii      | 58.50 |   735 |  180 |         2 |
| Char/Unicode    | 80.27 |   747 |    4 |         1 |
| AnyChar/Ascii   | 52.12 |   825 |  176 |         1 |
| AnyChar/Unicode | 74.61 |   804 |    0 |         0 |
| Take/Ascii      | 72.04 |   597 |  176 |         1 |
| Take/Unicode    | 93.07 |   645 |    0 |         0 |
| Until/Ascii     |  6.40 |  6724 |    0 |         0 |
| Until/Unicode   | 11.67 |  5142 |    0 |         0 |
| Any/Ascii       | 38.42 |  1119 |    0 |         0 |
| Any/Unicode     | 60.88 |   986 |    0 |         0 |
| Not/Ascii       | 64.76 |   664 |    0 |         0 |
| Not/Unicode     | 114.6 |   523 |    0 |         0 |
| OneOf/Ascii     | 63.77 |   674 |  176 |         1 |
| OneOf/Unicode   | 75.27 |   797 |    0 |         0 |
| NoneOf/Ascii    | 53.70 |   801 |  176 |         1 |
| NoneOf/Unicode  | 78.33 |   766 |    0 |         0 |

### Predicate Combinators

| Benchmark              | ns/op | MB/s | B/op | allocs/op |
| ---------------------- | ----: | ---: | ---: | --------: |
| While/Digit            | 27.72 | 1912 |    0 |         0 |
| While/Letter/Ascii     | 10.16 | 4231 |    0 |         0 |
| While/Letter/Unicode   | 243.2 |  247 |    0 |         0 |
| While/Alphanumeric     | 58.19 |  275 |    0 |         0 |
| While/Space            | 21.10 | 2275 |    0 |         0 |
| WhileNot/Digit/Ascii   | 121.5 |  378 |    0 |         0 |
| WhileNot/Digit/Unicode | 197.8 |  318 |    0 |         0 |
| Satisfy/Ascii          | 58.74 |  732 |  180 |         2 |
| Satisfy/Unicode        | 81.26 |  738 |    4 |         1 |

### Sequence Combinators

| Benchmark             | ns/op | MB/s | B/op | allocs/op |
| --------------------- | ----: | ---: | ---: | --------: |
| Pair/Ascii            | 40.04 | 1074 |   48 |         2 |
| Pair/Unicode          | 39.42 | 1522 |   48 |         2 |
| Delimited/Parentheses |  8.44 | 2487 |    0 |         0 |
| Delimited/Quotes      |  8.63 | 2432 |    0 |         0 |
| SepPair               | 64.19 |  218 |   48 |         2 |
| All/ThreeTags         | 64.86 |  663 |  112 |         3 |
| All/FiveTags          | 102.2 |  421 |  240 |         4 |

### Modifier Combinators

| Benchmark    | ns/op |  MB/s |  B/op | allocs/op |
| ------------ | ----: | ----: | ----: | --------: |
| Opt/Match    |  1.81 | 23783 |     0 |         0 |
| Opt/NoMatch  |  1.81 | 25486 |     0 |         0 |
| Map          | 16.47 |  2914 |     0 |         0 |
| Many/Small   | 127.7 |   110 |   288 |         5 |
| Many/Medium  | 605.1 |   172 |  2208 |         8 |
| Many/Large   |  4324 |   232 | 18848 |        11 |
| Peek/Ascii   | 63.05 |   682 |   176 |         1 |
| Peek/Unicode | 93.80 |   640 |     0 |         0 |
| Flatten      | 87.92 |   489 |   128 |         4 |

### Control Flow Combinators

| Benchmark         | ns/op |   MB/s | B/op | allocs/op |
| ----------------- | ----: | -----: | ---: | --------: |
| First/FirstMatch  |  3.96 |  10868 |    0 |         0 |
| First/LastMatch   | 118.2 |    118 |  128 |         4 |
| Verify/Pass       | 16.66 |   2881 |    0 |         0 |
| Recognize/Ascii   | 163.4 |    263 |  400 |         6 |
| Recognize/Unicode | 213.5 |    281 |   80 |         4 |
| Consumed          | 77.71 |    553 |  112 |         4 |
| Eof               |  0.23 |      - |    0 |         0 |
| AllConsuming      |  3.71 |   5125 |    0 |         0 |
| Rest/Ascii        |  0.23 | 190713 |    0 |         0 |
| Rest/Unicode      |  0.23 | 265547 |    0 |         0 |
| Value             |  3.71 |  12677 |    0 |         0 |
| Cond/True         |  3.48 |  12340 |    0 |         0 |
| Cond/False        |  1.24 |  34555 |    0 |         0 |
| Cut               |  3.67 |  11720 |    0 |         0 |

### Parser Combinators

| Benchmark   | ns/op |  MB/s | B/op | allocs/op |
| ----------- | ----: | ----: | ---: | --------: |
| Crlf        |  3.44 | 13096 |    0 |         0 |
| Eol/Ascii   | 157.9 |   279 |   80 |         3 |
| Eol/Unicode | 143.0 |   427 |   80 |         3 |

### Scaling Benchmarks

| Benchmark           | ns/op |  MB/s | B/op | allocs/op |
| ------------------- | ----: | ----: | ---: | --------: |
| UntilScaling/Small  |  5.46 |  1648 |    0 |         0 |
| UntilScaling/Medium |  6.75 | 15400 |    0 |         0 |
| UntilScaling/Large  | 109.3 | 91541 |    0 |         0 |
| WhileScaling/Small  | 11.15 |   538 |    0 |         0 |
| WhileScaling/Medium | 278.3 |   370 |    0 |         0 |
| WhileScaling/Large  | 29807 |   336 |    0 |         0 |

### Real-World Patterns

| Benchmark     | ns/op | MB/s | B/op | allocs/op |
| ------------- | ----: | ---: | ---: | --------: |
| KeyValuePair  | 48.10 |  665 |   48 |         2 |
| GitDiffHeader | 102.1 |  186 |  160 |         5 |
| CSVField      | 32.89 |  912 |    0 |         0 |
