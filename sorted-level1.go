package main

import (
	"slices"
	"sort"
	"strings"
)

func sortAndCutMatches(matches []Place, k int) []Place {
	sort.Slice(matches, func(i, j int) bool {

		// If two cities have the same population, sort desc by place ID
		if matches[i].Population == matches[j].Population {
			return matches[i].ID > matches[j].ID
		}

		return matches[i].Population > matches[j].Population
	})

	if len(matches) > k {
		matches = matches[:k]
		matches = slices.Clone(matches)
	}
	return matches
}

func suggestSorted(places []Place, prefix string, k int) []Place {
	var matches []Place

	if k > 0 && len(prefix) > 0 {
		lowerCasedPrefix := lowerTrim(prefix) // Trim only the strting whitespaces

		startIndex, _ := slices.BinarySearchFunc(places, lowerCasedPrefix, func(p Place, target string) int {
			return strings.Compare(p.LowerName, target)
		})

		shiftedPrefix := shiftLastCharacter(lowerCasedPrefix)
		relativeIndex, _ := slices.BinarySearchFunc(places[startIndex:], shiftedPrefix, func(p Place, target string) int {
			return strings.Compare(p.LowerName, target)
		})
		endIndex := startIndex + relativeIndex

		matches = places[startIndex:endIndex]
		block := slices.Clone(matches) // Cloning to protect original matches slice

		matches = sortAndCutMatches(block, k)
	}

	return matches
}
