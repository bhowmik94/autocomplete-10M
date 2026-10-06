package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

func printMem(label string) {
	runtime.GC() // collect garbage first so the number is meaningful
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%s: HeapAlloc = %.1f MB, NumGC = %d\n",
		label, float64(m.HeapAlloc)/1024/1024, m.NumGC)
}

func main() {
	// Read the file into mamory
	printMem("before load")

	// configurable file path, default
	dataPath := flag.String("data", "data/allCountries.txt", "path to GeoNames file")
	flag.Parse()

	contentBytes, err := os.ReadFile(*dataPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Convert bytes to a string for manipulation
	content := string(contentBytes)

	// Perform operations
	start := time.Now()
	cityData := loadPlaces(content)
	fmt.Println("load time:", time.Since(start))
	printMem("after load")

	// Search by prefix
	matchedCities := naiveSuggest(cityData, "Abu", 4)

	fmt.Printf("Matched entries are: %v\n", matchedCities)
}
