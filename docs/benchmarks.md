# Benchmarks

Benchmarks are run on an Apple M4 Pro (darwin/arm64) using Go's built-in testing framework.

```
go test -bench=. -benchmem ./...
```

## Results

### Basic Combinators

| Benchmark         | ns/op |  MB/s | B/op | allocs/op |
| ----------------- | ----: | ----: | ---: | --------: |
| Tag/Ascii         |  1.49 | 28874 |    0 |         0 |
| Tag/Unicode       |  1.54 | 38849 |    0 |         0 |
| TagNoCase/Ascii   |  6.92 |  6217 |    0 |         0 |
| Char/Ascii        |  1.77 | 24247 |    0 |         0 |
| Char/Unicode      |  2.48 | 24239 |    0 |         0 |
| AnyChar/Ascii     |  1.78 | 24203 |    0 |         0 |
| AnyChar/Unicode   |  2.24 | 26819 |    0 |         0 |
| Take/Ascii        | 16.04 |  2681 |    0 |         0 |
| Take/Unicode      | 10.25 |  5856 |    0 |         0 |
| Until/Ascii       |  6.45 |  6666 |    0 |         0 |
| Until/Unicode     | 12.66 |  4740 |    0 |         0 |
| Any/Small/Ascii   |  9.32 |  4614 |    0 |         0 |
| Any/Large/Ascii   | 70.29 |   612 |    0 |         0 |
| Any/Small/Unicode | 28.23 |  2126 |    0 |         0 |
| Any/Large/Unicode | 60.30 |   995 |    0 |         0 |
| Not/Small/Ascii   | 52.11 |   825 |    0 |         0 |
| Not/Large/Ascii   | 105.6 |   407 |    0 |         0 |
| Not/Small/Unicode | 93.54 |   641 |    0 |         0 |
| Not/Large/Unicode | 106.9 |   562 |    0 |         0 |
| OneOf/Ascii       | 11.38 |  3779 |    0 |         0 |
| OneOf/Unicode     |  4.46 | 13452 |    0 |         0 |
| NoneOf/Ascii      |  3.67 | 11709 |    0 |         0 |
| NoneOf/Unicode    |  6.44 |  9324 |    0 |         0 |

### Predicate Combinators

| Benchmark              | ns/op |  MB/s | B/op | allocs/op |
| ---------------------- | ----: | ----: | ---: | --------: |
| While/Digit            | 22.92 |  2313 |    0 |         0 |
| While/Letter/Ascii     |  8.69 |  4947 |    0 |         0 |
| While/Letter/Unicode   | 216.4 |   277 |    0 |         0 |
| While/Alphanumeric     | 49.99 |   320 |    0 |         0 |
| While/Space            | 15.74 |  3050 |    0 |         0 |
| WhileNot/Digit/Ascii   | 100.8 |   456 |    0 |         0 |
| WhileNot/Digit/Unicode | 176.0 |   358 |    0 |         0 |
| Satisfy/Ascii          |  2.39 | 17980 |    0 |         0 |
| Satisfy/Unicode        |  2.97 | 20203 |    0 |         0 |

### Sequence Combinators

| Benchmark             | ns/op | MB/s | B/op | allocs/op |
| --------------------- | ----: | ---: | ---: | --------: |
| Pair/Ascii            | 41.16 | 1045 |   48 |         2 |
| Pair/Unicode          | 40.89 | 1468 |   48 |         2 |
| Delimited/Parentheses |  8.24 | 2550 |    0 |         0 |
| Delimited/Quotes      |  8.48 | 2476 |    0 |         0 |
| SepPair               | 60.00 |  233 |   48 |         2 |
| All/ThreeTags         | 68.62 |  627 |  112 |         3 |
| All/FiveTags          | 105.0 |  409 |  240 |         4 |

### Modifier Combinators

| Benchmark    | ns/op |  MB/s |  B/op | allocs/op |
| ------------ | ----: | ----: | ----: | --------: |
| Opt/Match    |  2.15 | 20017 |     0 |         0 |
| Opt/NoMatch  |  2.20 | 20937 |     0 |         0 |
| Map          | 14.52 |  3307 |     0 |         0 |
| Many/Small   | 130.4 |   107 |   288 |         5 |
| Many/Medium  | 615.4 |   169 |  2208 |         8 |
| Many/Large   |  4270 |   235 | 18848 |        11 |
| Peek/Ascii   |  7.29 |  5895 |     0 |         0 |
| Peek/Unicode | 11.05 |  5429 |     0 |         0 |
| Flatten      | 91.81 |   468 |   128 |         4 |

### Control Flow Combinators

| Benchmark         | ns/op |   MB/s | B/op | allocs/op |
| ----------------- | ----: | -----: | ---: | --------: |
| First/FirstMatch  |  3.97 |  10833 |    0 |         0 |
| First/LastMatch   | 124.9 |    112 |  128 |         4 |
| Verify/Pass       | 14.37 |   3341 |    0 |         0 |
| Recognize/Ascii   | 45.99 |    935 |   48 |         2 |
| Recognize/Unicode | 49.37 |   1215 |   48 |         2 |
| Consumed          | 81.93 |    525 |  112 |         4 |
| Eof               |  0.23 |      - |    0 |         0 |
| AllConsuming      |  3.73 |   5100 |    0 |         0 |
| Rest/Ascii        |  0.23 | 190050 |    0 |         0 |
| Rest/Unicode      |  0.23 | 265677 |    0 |         0 |
| Value             |  3.51 |  13393 |    0 |         0 |
| Cond/True         |  3.48 |  12358 |    0 |         0 |
| Cond/False        |  1.34 |  32030 |    0 |         0 |
| Cut               |  3.72 |  11564 |    0 |         0 |

### Parser Combinators

| Benchmark   | ns/op |  MB/s | B/op | allocs/op |
| ----------- | ----: | ----: | ---: | --------: |
| Crlf        |  3.24 | 13871 |    0 |         0 |
| Eol/Ascii   | 58.23 |   756 |    0 |         0 |
| Eol/Unicode | 56.85 |  1073 |    0 |         0 |

### Scaling Benchmarks

| Benchmark           | ns/op |  MB/s | B/op | allocs/op |
| ------------------- | ----: | ----: | ---: | --------: |
| UntilScaling/Small  |  5.66 |  1589 |    0 |         0 |
| UntilScaling/Medium |  6.98 | 14909 |    0 |         0 |
| UntilScaling/Large  | 115.9 | 86300 |    0 |         0 |
| WhileScaling/Small  |  8.70 |   690 |    0 |         0 |
| WhileScaling/Medium | 205.6 |   501 |    0 |         0 |
| WhileScaling/Large  | 19816 |   505 |    0 |         0 |

### Real-World Patterns

| Benchmark     | ns/op | MB/s | B/op | allocs/op |
| ------------- | ----: | ---: | ---: | --------: |
| KeyValuePair  | 51.34 |  623 |   48 |         2 |
| GitDiffHeader | 102.4 |  185 |  160 |         5 |
| CSVField      | 26.64 | 1126 |    0 |         0 |
