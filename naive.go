package main

import (
	"fmt"
	"sort"
	"strings"
)

func naiveSuggest(cities []City, prefix string, k int) {
	var matches []City
	lowerCasedPrefix := strings.ToLower(strings.TrimSpace(prefix))

	for i := 0; i < len(cities); i++ {
		if strings.HasPrefix(strings.ToLower(cities[i].Name), lowerCasedPrefix) {
			matches = append(matches, cities[i])
		}
	}

	if k >= 0 { // only positive k value allowed
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

	fmt.Printf("Matched entries are: %v\n", matches)
}
