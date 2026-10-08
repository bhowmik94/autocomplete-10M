## Level 0: Naive Version

Machine: Ryzen 5 5500U, Windows, Go 1.26.4\
Dataset: cities15000 (~34k rows)

| Level | Prefix | ns/op     | B/op   | allocs/op |
|-------|--------|-----------|--------|-----------|
| 0     | a      | 3,995,848 | 777536 | 34167     |
| 0     | dha    | 3,543,568 | 462144 | 34160     |
| 0     | dhaka  | 3,534,650 | 458688 | 34153     |
| 0     | xyzq   | 3,585,912 | 458616 | 34151     |

Load (cities15000, 34k rows): ~18 ms\
Baseline heap (empty program): 0.3 MB\
Heap after load, v1 (substrings of file text): 10.2 MB\
Heap after load, v2 (strings.Clone):            2.5 MB\

Dataset: allCountries, 13.47M rows | Load: 9.8 s | Heap after load: 993.5 MB

### Query benchmarks (k = 10)

| Level | Prefix | ns/op         | B/op        | allocs/op  |
|-------|--------|---------------|-------------|------------|
| 0     | a      | 1,757,419,860 | 413,568,708 | 13,469,051 |
| 0     | dha    | 1,413,401,000 | 243,977,918 | 13,469,033 |
| 0     | dhaka  | 1,472,474,380 | 241,684,158 | 13,469,023 |
| 0     | xyzq   | 1,452,689,980 | 241,651,961 | 13,469,011 |

Load (allCountries, 13.47M rows): ~9.7 s\
Baseline heap (empty program):       0.3 MB\
Heap after load, v2 (strings.Clone): 993.5 MB\

## Level 1: sorted array + binary search

Machine: AMD Ryzen 5 5500U, Windows, Go 1.26.4\
Dataset: allCountries (13,472,243 rows)

### Load

| Metric                         | Level 0 (v2)  | Level 1       |
|--------------------------------|---------------|---------------|
| Baseline heap (empty program)  | 0.3 MB        | 0.3 MB        |
| Parse time                     | ~9.7 s        | 10.5 s        |
| Sort time                      | -             | 10.0 s        |
| Total load time                | ~9.7 s        | ~20.5 s       |
| Heap after load                | 993.5 MB      | 1403.6 MB     |
| GC runs during load            | 62            | 58            |

### Query benchmarks (k = 10)

| Level | Prefix | ns/op         | B/op        | allocs/op  |
|-------|--------|---------------|-------------|------------|
| 0     | a      | 1,757,419,860 | 413,568,708 | 13,469,051 |
| 0     | dha    | 1,413,401,000 | 243,977,918 | 13,469,033 |
| 0     | dhaka  | 1,472,474,380 | 241,684,158 | 13,469,023 |
| 0     | xyzq   | 1,452,689,980 | 241,651,961 | 13,469,011 |
| 1     | a      | 144,189,800   | 37,692,238  | 6          |
| 1     | dha    | 2,134,666     | 770,896     | 6          |
| 1     | dhaka  | 19,322        | 14,416      | 6          |
| 1     | xyzq   | 572           | 32          | 2          |

Naive re-measured on the sorted, Level 1 data layout:

| Prefix | ns/op         | B/op        | allocs/op  |
|--------|---------------|-------------|------------|
| a      | 2,635,491,220 | 454,659,704 | 13,469,052 |
| dha    | 2,341,362,720 | 245,427,838 | 13,469,035 |
| dhaka  | 2,392,271,300 | 241,692,286 | 13,469,024 |
| xyzq   | 2,326,834,160 | 241,651,952 | 13,469,011 |

### Speedup vs naive (same run)

| Prefix | Naive   | Sorted  | Speedup     |
|--------|---------|---------|-------------|
| xyzq   | 2.33 s  | 572 ns  | ~4,000,000x |
| dhaka  | 2.39 s  | 19.3 µs | ~124,000x   |
| dha    | 2.34 s  | 2.13 ms | ~1,100x     |
| a      | 2.64 s  | 144 ms  | ~18x        |

### Findings

- Query cost now depends on the size of the matching block, not on the dataset size.
  Block sizes (estimated from B/op / ~64-byte Place): `a` ~590k places, `dha` ~12k.
- Time per place in the block is ~180-250 ns, which is the cost of copying and sorting it.
- `xyzq` (no match) costs two binary searches only: 572 ns.
- The one-letter prefix is still slow (144 ms, 37.7 MB per query). This motivates Level 2.
- Trade-off: +410 MB heap (the `LowerName` field and string data) and ~2x load time
  (10 s sort) in exchange for query speedups of 18x up to millions.
- Naive got slower after sorting (~1.45 s -> ~2.35 s). Untested hypothesis: cache misses,
  because the sorted structs point to name strings still laid out in file order.

### Reproduce

```bash
git checkout level1/binary-search
$env:DATA_FILE = "data/allCountries.txt"
go test -bench "Suggest/naive" -benchmem -run "^$" -benchtime=5x
go test -bench "Suggest/sorted" -benchmem -run "^$"
```
