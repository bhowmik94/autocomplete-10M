### Naive Version (Level 0)

Machine: Ryzen 5 5500U, Windows, Go 1.26.4
Dataset: cities15000 (~34k rows)

| Level | Prefix | ns/op     | B/op   | allocs/op |
|-------|--------|-----------|--------|-----------|
| 0     | a      | 3,995,848 | 777536 | 34167     |
| 0     | dha    | 3,543,568 | 462144 | 34160     |
| 0     | dhaka  | 3,534,650 | 458688 | 34153     |
| 0     | xyzq   | 3,585,912 | 458616 | 34151     |

Load (cities15000, 34k rows): ~18 ms
Baseline heap (empty program): 0.3 MB
Heap after load, v1 (substrings of file text): 10.2 MB
Heap after load, v2 (strings.Clone):            2.5 MB