package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Place struct {
	ID          int
	Name        string
	LowerName   string
	CountryCode string
	Population  int
}

func buildIndex(places []Place) []Place {
	for i := range places {
		places[i].LowerName = strings.ToLower(places[i].Name)
	}
	t1 := time.Now()
	slices.SortFunc(places, func(a, b Place) int {
		return strings.Compare(a.LowerName, b.LowerName)
	})
	sortTime := time.Since(t1)
	fmt.Println("sort time:", sortTime)
	return places
}

func loadPlacesFromFile(path string) ([]Place, error) {
	t0 := time.Now()
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
	parseTime := time.Since(t0)
	fmt.Println("parse time:", parseTime)

	return buildIndex(places), nil
}

func loadPlaces(scanner *bufio.Scanner) []Place {
	var placeData []Place

	for scanner.Scan() {
		// Gets the current line text
		line := scanner.Text()

		lineContents := strings.Split(line, "\t")

		if len(lineContents) >= 15 { // Check for omitting the last empty line in file
			var placeItem Place

			id, err := strconv.Atoi(lineContents[0]) // String to number conversion
			if err == nil {
				placeItem.ID = id
			}

			population, err := strconv.Atoi(lineContents[14])
			if err == nil {
				placeItem.Population = population
			}
			placeItem.CountryCode = strings.Clone(lineContents[8])
			placeItem.Name = strings.Clone(lineContents[1])

			placeData = append(placeData, placeItem)
		}
	}

	return placeData
}
