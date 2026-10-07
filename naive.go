package main

import (
	"strings"
)

func naiveSuggest(cities []Place, prefix string, k int) []Place {
	var matches []Place

	if k >= 0 { // only positive k value allowed
		lowerCasedPrefix := lowerTrim(prefix) // Trim only the strting whitespaces

		for i := 0; i < len(cities); i++ {
			if strings.HasPrefix(strings.ToLower(cities[i].Name), lowerCasedPrefix) {
				matches = append(matches, cities[i])
			}
		}

		matches = sortAndCutMatches(matches, k)
	}

	return matches
}
