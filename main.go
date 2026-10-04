package main

import (
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
	contentBytes, err := os.ReadFile("data/cities15000.txt")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Convert bytes to a string for manipulation
	content := string(contentBytes)

	// Perform operations
	start := time.Now()
	cityData := loadCities(content)
	fmt.Println("load time:", time.Since(start))
	printMem("after load")

	// Search by prefix
	matchedCities := naiveSuggest(cityData, "Abu", 4)

	fmt.Printf("Matched entries are: %v\n", matchedCities)
}
