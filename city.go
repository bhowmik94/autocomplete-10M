package main

import (
	"strconv"
	"strings"
)

type City struct {
	ID          int
	Name        string
	CountryCode string
	Population  int
}

func loadCities(input string) []City {

	// Convert into array of strings by new line
	lines := strings.Split(input, "\n")
	var cityData []City

	for i := 0; i < len(lines); i++ {
		lineContents := strings.Split(lines[i], "\t")

		// Check till Population field (Column 15)
		if len(lineContents) >= 15 { // Check for omitting the last empty line in file
			var cityItem City

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
			cityItem.CountryCode = lineContents[8]
			cityItem.Name = lineContents[1]

			cityData = append(cityData, cityItem)
		}
	}

	return cityData
}
