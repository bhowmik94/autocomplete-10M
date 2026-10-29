# Autocomplete Suggestion

A fast prefix-search API in Go, standard library only. Type `dha`, get the 10 most populous places whose names start with `dha`, out of 13.5 million.

It is a learning project: I built the simplest version first, measured it, and improved it step by step. On the full GeoNames dataset, a query for the prefix `a` went from **2.6 s to 2.9 ms**, and a query for `dhaka` from **2.4 s to 1.9 µs**.

Read the write-up: [article link](https://dev.to/sourav_bhowmik_73d35592ab/from-full-table-scans-to-millisecond-search-in-go-optimizing-prefix-search-for-135-million-places-3i2h)

## Quick start

Requirements: Go 1.21 or newer (see `go.mod`).

1. Download `cities15000.zip` from <https://download.geonames.org/export/dump/>. Create a `data` folder in the project root and unzip the file into it, so you have `data/cities15000.txt`. (About 34,000 places, loads instantly.)
2. Start the server:
   ```bash
   go run . -data data/cities15000.txt
   ```
3. Open <http://localhost:8080> for the search page, or call the API:
   ```bash
   curl "http://localhost:8080/suggest?q=ber&limit=5"
   ```

### Full dataset (13.5 million places)

Download `allCountries.zip` from the same page, unzip it into `data/` (the file is over 1 GB), and run:

```bash
go run .
```

Startup takes about 20 seconds (parsing and sorting) and the index needs about 1.4 GB of RAM.

| Flag | Default | Meaning |
| --- | --- | --- |
| `-data` | `data/allCountries.txt` | Path to the GeoNames file |
| `-addr` | `localhost:8080` | Address to listen on. Use `:8080` to accept connections from other machines or containers |

The `data/` folder is not committed to git.

## API

`GET /suggest?q=<prefix>&limit=<n>`

| Parameter | Rules |
| --- | --- |
| `q` | Prefix to search for. Case-insensitive, leading spaces ignored. An empty `q` returns `[]` |
| `limit` | Number of results. Default 10, values outside 1 to 50 are clamped, a non-integer returns `400` |

The response is a JSON array, best match first. Each item has `id`, `name`, `country_code` and `population`. Matches are ranked by population (highest first); on a tie the higher `id` wins. No match returns `[]`.

Other routes: `GET /healthz` returns `200`, and `GET /` serves the search page.

## Three versions

Each version is tagged, so you can check out any of them:

| Level | Idea | Git tag |
| --- | --- | --- |
| 0. Naive | Scan every place, lowercase its name, collect the matches, sort them | `level-0-naive` |
| 1. Binary search | Lowercase names and sort once at startup. Two binary searches find the *block* of matching places. Copy the block and sort it | `level-1-binary-search` |
| 2. Top-k | Same as Level 1, but walk the block once and keep only the best k. No copy, no sort | `level-2-topk` |

```bash
git checkout level-0-naive
```

**What is a block?** After sorting by name, every place that starts with the same prefix sits in one unbroken stretch of the array. The first binary search finds where it starts. The second finds where it ends by searching for the prefix with its last character increased by one (`dha` becomes `dhb`).

## Results

13,472,243 places, k = 10, AMD Ryzen 5 5500U, Windows, Go 1.26.4. Time per query, one run per level:

| Prefix | Matching places | Naive | Level 1 | Level 2 |
| --- | --- | --- | --- | --- |
| `a` | about 590,000 | 2.64 s | 144 ms | 2.9 ms |
| `dha` | about 12,000 | 2.34 s | 2.1 ms | 47 µs |
| `dhaka` | about 200 | 2.39 s | 19 µs | 1.9 µs |
| `xyzq` | 0 | 2.33 s | 572 ns | 454 ns |

Matching-place counts are estimates. Full numbers, memory per query, load times and notes are in [RESULTS.md](RESULTS.md).

## Tests and benchmarks

Run the tests (fast, no data file needed):

```bash
go test ./...
```

They use a small built-in sample and cover: a shared table of cases that runs against every implementation, a check that each level returns the same results as the naive one, and the HTTP handler (including 50 concurrent requests).

Benchmarks load the real data file. Set `DATA_FILE` (default `data/allCountries.txt`) and pick an implementation (`naive`, `sorted` or `topk`).

PowerShell:

```powershell
$env:DATA_FILE = "data/allCountries.txt"; go test -bench "Suggest/topk" -benchmem -run "^$"
```

bash:

```bash
DATA_FILE=data/allCountries.txt go test -bench 'Suggest/topk' -benchmem -run '^$'
```

The naive version takes seconds per query, so add `-benchtime=5x` when you benchmark it. If the data file is missing, the benchmarks are skipped.

## Limitations

- Only the name is matched: there is no accent folding, so `munchen` does not find München, and alternate names are not indexed.
- Places with the same name are listed separately, and ranking uses population only.
- Startup takes about 20 seconds and the index needs about 1.4 GB of RAM on the full dataset.
- The benchmarks call the search functions directly, so they leave out the network, HTTP parsing and JSON encoding. There is no end-to-end load test yet.

## Ideas for next steps

- Store populations and IDs in separate compact arrays so the scan reads less memory.
- A range-maximum structure (sparse table or segment tree) so the cost depends on k, not on the block size.
- Precomputed top 10 for the shortest prefixes.
- A Postgres version with an index and a cache, for comparison.

## Credits

Data: [GeoNames](https://www.geonames.org), licensed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).