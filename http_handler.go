package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

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
			limitInInt = 10

			// If strictly reject bad inputs instead of falling back:
			// http.Error(w, "Invalid limit parameter, must be an integer", http.StatusBadRequest)
			// return
		}
	}

	if limitInInt < minLimit {
		limitInInt = minLimit
	} else if limitInInt > maxLimit {
		limitInInt = maxLimit
	}

	results := s.suggest(s.places, queryStr, limitInInt)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
