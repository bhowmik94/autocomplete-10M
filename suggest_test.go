package main

import (
	"os"
	"slices"
	"strings"
	"sync"
	"testing"
)

type implementation struct {
	name string
	fn   func([]Place, string, int) []Place
}

// A slice, not a map, so the order is stable in test and benchmark output.
var implementations = []implementation{
	{"naive", naiveSuggest},
	{"sorted", suggestSorted},
}

func placeNames(places []Place) []string {
	names := make([]string, 0, len(places))
	for _, p := range places {
		names = append(names, p.Name)
	}
	return names
}

func testPlaces() []Place {
	return buildIndex([]Place{
		// "ber" prefix family, distinct populations
		{ID: 1, Name: "Berlin", CountryCode: "DE", Population: 3600000},
		{ID: 2, Name: "Bergen", CountryCode: "NO", Population: 285000},
		{ID: 3, Name: "Bern", CountryCode: "CH", Population: 134000},
		{ID: 4, Name: "Bergamo", CountryCode: "IT", Population: 120000},

		// contain "ber" but do NOT start with it (catches Contains vs HasPrefix)
		{ID: 5, Name: "Oberhausen", CountryCode: "DE", Population: 210000},
		{ID: 6, Name: "Amberg", CountryCode: "DE", Population: 42000},

		// ties: same population, so the tie-break rule decides the order
		{ID: 7, Name: "Bayreuth", CountryCode: "DE", Population: 77000},
		{ID: 8, Name: "Bamberg", CountryCode: "DE", Population: 77000},

		// same name, different cities
		{ID: 9, Name: "Springfield", CountryCode: "US", Population: 114000},
		{ID: 10, Name: "Springfield", CountryCode: "US", Population: 60000},

		// non-ASCII
		{ID: 11, Name: "München", CountryCode: "DE", Population: 1500000},
		{ID: 12, Name: "Münster", CountryCode: "DE", Population: 315000},
		{ID: 13, Name: "Mülheim", CountryCode: "DE", Population: 170000},

		// names with spaces
		{ID: 14, Name: "New York", CountryCode: "US", Population: 8400000},
		{ID: 15, Name: "New Delhi", CountryCode: "IN", Population: 250000},
		{ID: 16, Name: "Newark", CountryCode: "US", Population: 300000},

		// zero population
		{ID: 17, Name: "Tinyville", CountryCode: "US", Population: 0},

		// many matches for a one-letter prefix (tests the k cutoff)
		{ID: 18, Name: "Amsterdam", CountryCode: "NL", Population: 870000},
		{ID: 19, Name: "Athens", CountryCode: "GR", Population: 660000},
		{ID: 20, Name: "Augsburg", CountryCode: "DE", Population: 300000},
		{ID: 21, Name: "Aachen", CountryCode: "DE", Population: 250000},

		// pure non-matches
		{ID: 22, Name: "Paris", CountryCode: "FR", Population: 2100000},
		{ID: 23, Name: "Bremen", CountryCode: "DE", Population: 560000},
	})
}

func TestSuggest(t *testing.T) {
	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			places := testPlaces()            // fresh copy for each implementation
			for _, tc := range suggestCases { // possible test cases on testCases.go file
				t.Run(tc.name, func(t *testing.T) {
					got := placeNames(impl.fn(places, tc.prefix, tc.k))
					if !slices.Equal(got, tc.want) {
						t.Errorf("prefix %q, k=%d: got %v, want %v", tc.prefix, tc.k, got, tc.want)
					}
				})
			}
		})
	}
}

// Equivalence test
func TestSortedMatchesNaive(t *testing.T) {
	places := testPlaces()
	prefixes := []string{"", "zzz", " ber"}
	for _, p := range places {
		runes := []rune(strings.ToLower(p.Name))
		for n := 1; n <= len(runes); n++ {
			prefixes = append(prefixes, string(runes[:n]))
		}
	}
	for _, prefix := range prefixes {
		for _, k := range []int{-1, 0, 1, 3, 10, 100} {
			a := naiveSuggest(places, prefix, k)
			b := suggestSorted(places, prefix, k)
			if !slices.Equal(a, b) {
				t.Errorf("prefix %q, k=%d: naive %v, sorted %v", prefix, k, placeNames(a), placeNames(b))
			}
		}
	}
}

// ---- benchmarks ----

var (
	benchOnce   sync.Once
	benchPlaces []Place
	benchErr    error
)

// Loads the big file once for all benchmarks
func loadBenchPlaces(b *testing.B) []Place {
	b.Helper()
	benchOnce.Do(func() {
		path := os.Getenv("DATA_FILE")
		if path == "" {
			path = "data/allCountries.txt"
		}
		benchPlaces, benchErr = loadPlacesFromFile(path)
	})
	if benchErr != nil {
		b.Skipf("could not load data: %v", benchErr)
	}
	return benchPlaces
}

var sink []Place // keeps the compiler from optimizing the call away

func BenchmarkSuggest(b *testing.B) {
	places := loadBenchPlaces(b)
	prefixes := []string{"a", "dha", "dhaka", "xyzq"}
	for _, impl := range implementations {
		for _, p := range prefixes {
			b.Run(impl.name+"/"+p, func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					sink = impl.fn(places, p, 10)
				}
			})
		}
	}
}
