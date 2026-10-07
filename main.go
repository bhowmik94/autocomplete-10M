package main

import (
	"flag"
	"fmt"
	"log"
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
	dataPath := flag.String("data", "data/allCountries.txt", "path to GeoNames file")
	flag.Parse()

	printMem("before load")
	t0 := time.Now()
	places, err := loadPlacesFromFile(*dataPath)
	parseTime := time.Since(t0)
	fmt.Println("parse time:", parseTime)

	if err != nil {
		log.Fatalf("failed to load %s: %v", *dataPath, err)
	}
	printMem("after load")

	fmt.Printf("rows loaded: %d\n", len(places))
	// fmt.Printf("Matched entries are: %v\n", naiveSuggest(places, "Ber", 4))

	fmt.Printf("Matched entries are: %v\n", suggestSorted(places, "Ber", 4))
}
