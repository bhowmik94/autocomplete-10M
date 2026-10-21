package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
)

func get(t *testing.T, s *server, url string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.routes().ServeHTTP(rec, req) // goes through the router, like a real request
	return rec
}

func TestHandleSuggest(t *testing.T) {
	s := &server{places: testPlaces(), suggest: suggestTopK}

	rec := get(t, s, "/suggest?q=ber&limit=2")

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type %q, want application/json", ct)
	}

	var got []Place
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	want := []string{"Berlin", "Bergen"}
	if !slices.Equal(placeNames(got), want) {
		t.Errorf("got %v, want %v", placeNames(got), want)
	}
}

// A fake implementation captures the k that the handler passes down.
func TestLimitHandling(t *testing.T) {
	cases := []struct {
		name, url string
		wantK     int
	}{
		{"missing limit uses default", "/suggest?q=a", 10},
		{"valid limit", "/suggest?q=a&limit=5", 5},
		{"zero is clamped up", "/suggest?q=a&limit=0", 1},
		{"negative is clamped up", "/suggest?q=a&limit=-3", 1},
		{"huge is clamped down", "/suggest?q=a&limit=100000", 50},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotK := -1
			s := &server{suggest: func(_ []Place, _ string, k int) []Place {
				gotK = k
				return nil
			}}
			get(t, s, tc.url)
			if gotK != tc.wantK {
				t.Errorf("k = %d, want %d", gotK, tc.wantK)
			}
		})
	}
}

func TestInvalidLimitReturnsBadRequest(t *testing.T) {
	called := false
	s := &server{suggest: func([]Place, string, int) []Place {
		called = true
		return nil
	}}

	rec := get(t, s, "/suggest?q=a&limit=abc")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if called {
		t.Error("suggest function was called for an invalid limit")
	}
}

func TestEmptyResultIsArrayNotNull(t *testing.T) {
	s := &server{suggest: func([]Place, string, int) []Place { return nil }}
	rec := get(t, s, "/suggest?q=zzz")
	if body := strings.TrimSpace(rec.Body.String()); body != "[]" {
		t.Errorf("body %q, want []", body)
	}
}

func TestResponseHidesInternalFields(t *testing.T) {
	s := &server{places: testPlaces(), suggest: suggestTopK}
	rec := get(t, s, "/suggest?q=ber")
	if strings.Contains(rec.Body.String(), "LowerName") {
		t.Error("response leaks LowerName")
	}
}

func TestHealthzAndUnknownRoute(t *testing.T) {
	s := &server{}
	if rec := get(t, s, "/healthz"); rec.Code != http.StatusOK {
		t.Errorf("/healthz status %d, want 200", rec.Code)
	}
	if rec := get(t, s, "/nope"); rec.Code != http.StatusNotFound {
		t.Errorf("/nope status %d, want 404", rec.Code)
	}
}

func TestConcurrentRequests(t *testing.T) {
	s := &server{places: testPlaces(), suggest: suggestTopK}
	h := s.routes()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodGet, "/suggest?q=a&limit=3", nil)
			rec := httptest.NewRecorder() // one recorder per goroutine
			h.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Errorf("status %d", rec.Code) // Errorf, not Fatalf, inside goroutines
			}
		}()
	}
	wg.Wait()
}
