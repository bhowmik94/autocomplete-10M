package main

import (
	"slices"
	"strings"
)

func topK(matches []Place, k int) []Place {
	if k <= 0 || len(matches) == 0 {
		return []Place{}
	}
	podium := make([]Place, 0, min(k, len(matches)))

	for _, matchItem := range matches {
		var i int
		if len(podium) < k {
			podium = append(podium, matchItem) // podium not full yet: always take it
			i = len(podium) - 1                // why this logic here?
		} else if isBetter(matchItem, podium[len(podium)-1]) {
			i = len(podium) - 1 // full: p replaces the worst entry
		} else {
			continue // not good enough, skip
		}

		// bubble p left past every entry it beats
		for i > 0 && isBetter(matchItem, podium[i-1]) { // why checking i > 0?
			podium[i] = podium[i-1]
			i--
		}
		podium[i] = matchItem
	}
	return podium
}

func suggestTopK(places []Place, prefix string, k int) []Place {
	lowerCasedPrefix := lowerTrim(prefix)
	if k <= 0 || lowerCasedPrefix == "" {
		return []Place{}
	}

	cmp := func(p Place, target string) int {
		return strings.Compare(p.LowerName, target)
	}

	startIndex, _ := slices.BinarySearchFunc(places, lowerCasedPrefix, cmp)

	shiftedPrefix := shiftLastCharacter(lowerCasedPrefix)
	relativeIndex, _ := slices.BinarySearchFunc(places[startIndex:], shiftedPrefix, cmp)
	endIndex := startIndex + relativeIndex

	return topK(places[startIndex:endIndex], k)
}
