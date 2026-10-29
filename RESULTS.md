## Benchmark results

All query benchmarks use `k = 10`. Times and allocation figures are reported by Go's benchmark runner. Readable time equivalents are shown alongside the original `ns/op` values; `B/op` and `allocs/op` are unchanged.

## Level 0: naive full scan

### Small dataset: cities15000

Machine: Ryzen 5 5500U, Windows, Go 1.26.4  
Dataset: cities15000, approximately 34k rows.

| Prefix | ns/op (equivalent) | B/op | allocs/op |
|--------|-------------------:|-----:|----------:|
| `a`    | 3,995,848 (4.00 ms) | 777,536 | 34,167 |
| `dha`  | 3,543,568 (3.54 ms) | 462,144 | 34,160 |
| `dhaka`| 3,534,650 (3.53 ms) | 458,688 | 34,153 |
| `xyzq` | 3,585,912 (3.59 ms) | 458,616 | 34,151 |

Load time (cities15000, 34k rows): approximately 18 ms.  
Baseline heap (empty program): 0.3 MB.  
Heap after load, v1 (strings referencing file text): 10.2 MB.  
Heap after load, v2 (`strings.Clone`): 2.5 MB.

### Large dataset: allCountries

Dataset: allCountries, 13.47M rows. Load: approximately 9.7 s. Heap after load: 993.5 MB.

| Prefix | ns/op (equivalent) | B/op | allocs/op |
|--------|-------------------:|-----:|----------:|
| `a`    | 1,757,419,860 (1.76 s) | 413,568,708 | 13,469,051 |
| `dha`  | 1,413,401,000 (1.41 s) | 243,977,918 | 13,469,033 |
| `dhaka`| 1,472,474,380 (1.47 s) | 241,684,158 | 13,469,023 |
| `xyzq` | 1,452,689,980 (1.45 s) | 241,651,961 | 13,469,011 |

Baseline heap (empty program): 0.3 MB.  
Heap after load, v2 (`strings.Clone`): 993.5 MB.

### Findings

- Every query scans all 13.47M rows, even when the prefix has no matches. Runtime is therefore roughly 1.4–1.8 s for each of these prefixes.
- Each city name is lowercased during the scan. This creates roughly 13.47M allocations per query and accounts for most of the 242–414 MB allocated per operation.
- Matching rows are collected, sorted by population and ID, then cut to the requested top 10. This adds work proportional to the match count and stores all matches temporarily.
- The `xyzq` query illustrates the cost of a full scan: it returns nothing but still takes 1.45 s and allocates about 242 MB.

## Level 1: sorted array and binary search

Machine: AMD Ryzen 5 5500U, Windows, Go 1.26.4.  
Dataset: allCountries, 13,472,243 rows.

### Load costs

| Metric | Level 0 (v2) | Level 1 |
|--------|-------------:|--------:|
| Baseline heap (empty program) | 0.3 MB | 0.3 MB |
| Parse time | approximately 9.7 s | 10.5 s |
| Sort time | — | 10.0 s |
| Total load time | approximately 9.7 s | approximately 20.5 s |
| Heap after load | 993.5 MB | 1,403.6 MB |
| GC runs during load | 62 | 58 |

### Query benchmarks

| Level | Prefix | ns/op (equivalent) | B/op | allocs/op |
|------:|--------|------------------:|-----:|----------:|
| 0 | `a` | 1,757,419,860 (1.76 s) | 413,568,708 | 13,469,051 |
| 0 | `dha` | 1,413,401,000 (1.41 s) | 243,977,918 | 13,469,033 |
| 0 | `dhaka` | 1,472,474,380 (1.47 s) | 241,684,158 | 13,469,023 |
| 0 | `xyzq` | 1,452,689,980 (1.45 s) | 241,651,961 | 13,469,011 |
| 1 | `a` | 144,189,800 (144 ms) | 37,692,238 | 6 |
| 1 | `dha` | 2,134,666 (2.13 ms) | 770,896 | 6 |
| 1 | `dhaka` | 19,322 (19.3 us) | 14,416 | 6 |
| 1 | `xyzq` | 572 (572 ns) | 32 | 2 |

### Level 1 speedup over Level 0

Ratios use the Level 0 and Level 1 `ns/op` values shown in the query benchmark table. These values come from separate recorded runs, so treat the ratios as approximate comparisons.

| Prefix | Level 0 | Level 1 | Speedup |
|--------|--------:|--------:|--------:|
| `a` | 1.76 s | 144 ms | approximately 12x |
| `dha` | 1.41 s | 2.13 ms | approximately 663x |
| `dhaka` | 1.47 s | 19.3 us | approximately 76,200x |
| `xyzq` | 1.45 s | 572 ns | approximately 2,540,000x |

### Findings

- Binary search locates the matching name range, so query work depends mainly on the number of matching rows rather than the full dataset size. The approximate ranges inferred from `B/op` and a roughly 64-byte `Place` are 590k rows for `a` and 12k for `dha`.
- Level 1 copies and sorts the whole matching range by population and ID before keeping the top 10. This is why the broad `a` prefix still takes 144 ms and allocates 37.7 MB.
- A query with no matches, such as `xyzq`, needs only the two binary searches and takes 572 ns.
- The index adds about 410 MB of heap and approximately 10 seconds of sorting at load time. Based on the displayed benchmark values, Level 1's query speedups range from about 12x to 2.54 millionx.
- The naive scan measured slower on the sorted data layout (about 1.45 s to 2.35 s across the displayed same-run values). Cache misses are a possible explanation: sorting moves the structs while their name strings remain separately allocated. This is a hypothesis, not a measured cause.

### Reproduce

```powershell
git checkout level1/binary-search
$env:DATA_FILE = "data/allCountries.txt"
go test -bench "Suggest/naive" -benchmem -run "^$" -benchtime=5x
go test -bench "Suggest/sorted" -benchmem -run "^$"
```

## Level 2: binary search and single-pass top-k

Level 2 keeps Level 1's sorted index and binary searches. It scans the matching range once while maintaining an ordered top-10 podium, avoiding both a copy and a full sort of that range. Machine and dataset are the same as Level 1. Load cost is unchanged; the recorded parse time is approximately 9.7 s and sort time approximately 9.4 s.

### Query benchmarks

| Level | Prefix | ns/op (equivalent) | B/op | allocs/op |
|------:|--------|------------------:|-----:|----------:|
| 0 | `a` | 1,757,419,860 (1.76 s) | 413,568,708 | 13,469,052 |
| 0 | `dha` | 1,413,401,000 (1.41 s) | 243,977,918 | 13,469,035 |
| 0 | `dhaka` | 1,472,474,380 (1.47 s) | 241,684,158 | 13,469,024 |
| 0 | `xyzq` | 1,452,689,980 (1.45 s) | 241,651,961 | 13,469,011 |
| 1 | `a` | 144,189,800 (144 ms) | 37,692,238 | 6 |
| 1 | `dha` | 2,134,666 (2.13 ms) | 770,896 | 6 |
| 1 | `dhaka` | 19,322 (19.3 us) | 14,416 | 6 |
| 1 | `xyzq` | 572 (572 ns) | 32 | 2 |
| 2 | `a` | 2,922,128 (2.92 ms) | 708 | 2 |
| 2 | `dha` | 46,908 (46.9 us) | 712 | 2 |
| 2 | `dhaka` | 1,864 (1.86 us) | 712 | 2 |
| 2 | `xyzq` | 454 (454 ns) | 8 | 1 |

### Speedup over Level 0

Ratios use the Level 0 and Level 2 `ns/op` values shown above and are approximate because the values come from separate recorded runs.

| Prefix | Level 0 | Level 2 | Speedup |
|--------|--------:|--------:|--------:|
| `a` | 1.76 s | 2.92 ms | approximately 602x |
| `dha` | 1.41 s | 46.9 us | approximately 30,100x |
| `dhaka` | 1.47 s | 1.86 us | approximately 790,000x |
| `xyzq` | 1.45 s | 454 ns | approximately 3,200,000x |

### Findings

- Level 2 removes Level 1's full-range copy and sort. For prefix `a`, runtime drops from 144 ms to 2.92 ms and allocation drops from 37.7 MB to 708 B. For `dha`, runtime drops from 2.13 ms to 46.9 µs; for `dhaka`, from 19.3 µs to 1.86 µs.
- Memory per query stays around 700 B for matching prefixes because the podium stores at most `k` places. The `xyzq` no-match case uses 8 B and one allocation.
- Runtime still grows with the matching range because Level 2 scans every candidate in that range. For `a` (approximately 590k rows), the scan takes 2.92 ms, or about 5 ns per candidate on average.
- With `k = 10`, the ordered podium is small. Larger values of `k` increase the amount of podium maintenance work.
- Compared with Level 0, Level 2 combines Level 1's indexed range lookup with a top-k scan: it avoids both the full-dataset lowercase scan and Level 1's large temporary match copy and sort.

### Reproduce

```powershell
git checkout level2/top-k
$env:DATA_FILE = "data/allCountries.txt"
go test -bench "Suggest/topk" -benchmem -run "^$"
```
