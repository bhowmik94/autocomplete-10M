package main

import (
	"strconv"
	"strings"
)

type Place struct {
	ID          int
	Name        string
	CountryCode string
	Population  int
}

func loadPlaces(input string) []Place {

	// Convert into array of strings by new line
	lines := strings.Split(input, "\n")
	var cityData []Place

	for i := 0; i < len(lines); i++ {
		lineContents := strings.Split(lines[i], "\t")

		// Check till Population field (Column 15)
		if len(lineContents) >= 15 { // Check for omitting the last empty line in file
			var cityItem Place

			// String to number conversion for ID field
			id, err := strconv.Atoi(lineContents[0])
			if err == nil {
				cityItem.ID = id
			}
			// String to number conversion for Population field
			population, err := strconv.Atoi(lineContents[14])
			if err == nil {
				cityItem.Population = population
			}
			cityItem.CountryCode = strings.Clone(lineContents[8])
			cityItem.Name = strings.Clone(lineContents[1])

			cityData = append(cityData, cityItem)
		}
	}

	return cityData
}
