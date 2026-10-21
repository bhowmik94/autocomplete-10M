package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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
	addr := flag.String("addr", "localhost:8080", "HTTP server address")
	flag.Parse()

	printMem("before load")
	places, err := loadPlacesFromFile(*dataPath)

	if err != nil {
		log.Fatalf("failed to load %s: %v", *dataPath, err)
	}
	printMem("after load")

	s := &server{places: places, suggest: suggestTopK}
	srv := &http.Server{
		Addr:              *addr,
		Handler:           s.routes(), // HTTP request handler in separate server.go file
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

	// fmt.Printf("Matched entries are: %v\n", naiveSuggest(places, "Ber", 4))
	// fmt.Printf("Matched entries are: %v\n", suggestSorted(places, "Ber", 4))
	// fmt.Printf("Matched entries are: %v\n", suggestTopK(places, "Ber", 4))
}
