# Autocomplete Suggestion

This is a practice project designed to learn and experiment *autocomplete suggestion by keyword*, on a huge geolocation dataset (10 million rows). The main goal is to develop a lightning-fast Go API, that can search and return suggestion results in milliseconds.

## Setup

At the project root directory, run the following command in your terminal:

```bash
go run main.go
```

## Development Plan

To kickstart the project, I have started out with a naive version which employs simple searching and sorting logic to search for matches in the city lists. In future versions, I will apply more advanced methods (e.g. trie, top-K per node, etc.) and measure the differences in performance. The final goal is to compare those performance metrics for a more detailed breakdown.