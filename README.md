# Autocomplete Suggestion

This is a personal project designed to learn and experiment *autocomplete suggestion by keyword*, on a huge GeoNames dataset with 10 million rows. The end goal is to develop a lightning-fast Go API, that can search and return suggestion results in milliseconds.

## Setup

Visit this link: https://download.geonames.org/export/dump/ and search for `allCountries.zip`. Create a folder called `data` in the project root. Download the zip file and unzip it inside the data folder. 

At the project root directory, run the following command in your terminal:

```bash
go run .
```

## Testing

Currently, there is one unit test to test the `naiveSuggest` function. It also comes with a benchmark to measure the performance of the suggestion function.

Run the following command to run tests:

```bash
go test .
```

And run the following command for benchmark performance measurement:

```bash
$env:DATA_FILE = "data/allCountries.txt"; go test -v -bench "." -benchmem -run "^$" -benchtime=5x
```

## Development Plan

To kickstart the project, I have started out with a naive version which employs simple searching and sorting logic to search for matches in the city lists. In future versions, I will apply more advanced methods (e.g. trie, top-K per node, etc.) and measure the differences in performance. The final goal is to compare those performance metrics for a more detailed breakdown.