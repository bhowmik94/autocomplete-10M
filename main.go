package main

import (
	"log"
	"os"
)

func main() {
	// Read the file into mamory
	contentBytes, err := os.ReadFile("cities15000.txt")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Convert bytes to a string for manipulation
	content := string(contentBytes)

	// Perform operations
	cityData := loadCities(content)
	// fmt.Printf("First city data: %v\n", cityData[:5])

	// Search by prefix
	naiveSuggest(cityData, "Abu", 20)
}
