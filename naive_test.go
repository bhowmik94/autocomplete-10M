package main

import (
	"os"
	"slices"
	"testing"
)

func cityNames(cities []City) []string {
	var names []string
	for i := 0; i < len(cities); i++ {
		names = append(names, cities[i].Name)
	}
	return names
}

func TestNaiveSuggest(t *testing.T) {
	cities := []City{
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
	}

	tests := []struct {
		name   string
		prefix string
		k      int
		want   []string // want city names in order, not whole structs
	}{
		// basics
		{"basic prefix", "ber", 10, []string{"Berlin", "Bergen", "Bern", "Bergamo"}},
		{"uppercase prefix", "BER", 10, []string{"Berlin", "Bergen", "Bern", "Bergamo"}},
		{"mixed case prefix", "bErLiN", 10, []string{"Berlin"}},
		{"leading whitespace trim", "  ber", 10, []string{"Berlin", "Bergen", "Bern", "Bergamo"}},
		{"trailing space is kept", "new ", 10, []string{"New York", "New Delhi"}},
		{"full name as prefix", "berlin", 10, []string{"Berlin"}},

		// no results
		{"no match", "zzz", 10, []string{}},
		{"prefix longer than any name", "berlinnn", 10, []string{}},
		{"substring is not a prefix", "erg", 10, []string{}},

		// k handling
		{"k zero", "ber", 0, []string{}},
		{"k negative", "ber", -5, []string{}},
		{"k one", "ber", 1, []string{"Berlin"}},
		{"k smaller than matches", "a", 3, []string{"Amsterdam", "Athens", "Augsburg"}},
		{"k equals matches", "ber", 4, []string{"Berlin", "Bergen", "Bern", "Bergamo"}},
		{"k larger than matches", "a", 100, []string{"Amsterdam", "Athens", "Augsburg", "Aachen", "Amberg"}},

		// ordering and ties
		{"tie broken by name", "ba", 10, []string{"Bamberg", "Bayreuth"}},
		{"several names sharing a prefix", "new", 10, []string{"New York", "Newark", "New Delhi"}},
		{"zero population still returned", "tiny", 10, []string{"Tinyville"}},

		// non-ASCII
		{"umlaut prefix", "mü", 10, []string{"München", "Münster", "Mülheim"}},
		{"uppercase umlaut prefix", "MÜ", 10, []string{"München", "Münster", "Mülheim"}},

		// space handling
		{"prefix with inner space", "new y", 10, []string{"New York"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := naiveSuggest(cities, tc.prefix, tc.k)
			gotNames := cityNames(got)

			if !slices.Equal(gotNames, tc.want) {
				t.Errorf("prefix %v, got %v, want %v: ", tc.prefix, got, tc.want)
			}
		})
		t.Skip()
	}
}

var sink []City // keeps the compiler from optimizing the call away

func BenchmarkNaiveSuggest(b *testing.B) {
	content, err := os.ReadFile("data/cities15000.txt") // adjust path to where your file is
	if err != nil {
		b.Skip("data file not found, skipping benchmark")
	}
	cities := loadCities(string(content)) // your loader, called ONCE

	prefixes := []string{"a", "dha", "dhaka", "xyzq"}
	for _, p := range prefixes {
		b.Run(p, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sink = naiveSuggest(cities, p, 10)
			}
		})
	}
}
