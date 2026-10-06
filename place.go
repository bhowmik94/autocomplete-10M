package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Place struct {
	ID          int
	Name        string
	CountryCode string
	Population  int
}

func loadPlacesFromFile(path string) ([]Place, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	const maxLine = 1024 * 1024 // 1 MB per line
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLine)

	places := loadPlaces(scanner)
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return places, nil
}

func loadPlaces(scanner *bufio.Scanner) []Place {
	var cityData []Place

	for scanner.Scan() {
		// Gets the current line text
		line := scanner.Text()

		lineContents := strings.Split(line, "\t")

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
