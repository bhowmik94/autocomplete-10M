package main

import (
	"sort"
	"strings"
	"unicode"
)

func naiveSuggest(cities []Place, prefix string, k int) []Place {
	var matches []Place

	if k >= 0 { // only positive k value allowed
		lowerCasedPrefix := strings.ToLower(strings.TrimLeftFunc(prefix, unicode.IsSpace)) // Trim only the strting whitespaces

		for i := 0; i < len(cities); i++ {
			if strings.HasPrefix(strings.ToLower(cities[i].Name), lowerCasedPrefix) {
				matches = append(matches, cities[i])
			}
		}
		sort.Slice(matches, func(i, j int) bool {

			// If two cities have the same population, sort asc by city name
			if matches[i].Population == matches[j].Population {
				return matches[i].Name < matches[j].Name
			}

			return matches[i].Population > matches[j].Population
		})

		if len(matches) > k {
			matches = matches[:k]
		}
	}

	return matches
}
