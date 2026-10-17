package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type server struct {
	places  []Place
	suggest func([]Place, string, int) []Place // which suggestion function to use
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/suggest", s.handleSuggest)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return mux
}

func (s *server) handleSuggest(w http.ResponseWriter, r *http.Request) {
	queryStr := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")

	limitInInt := 10
	const minLimit = 1
	const maxLimit = 50

	if limitStr != "" {
		var err error
		limitInInt, err = strconv.Atoi(limitStr)

		if err != nil {
			// If strictly reject bad inputs instead of falling back:
			http.Error(w, "Invalid limit parameter, must be an integer", http.StatusBadRequest)
			return
		}
	}

	if limitInInt < minLimit {
		limitInInt = minLimit
	} else if limitInInt > maxLimit {
		limitInInt = maxLimit
	}

	results := s.suggest(s.places, queryStr, limitInInt)
	if results == nil {
		results = []Place{} // prevents sending nil
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
