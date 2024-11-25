package framework

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkNewRouter(b *testing.B) {
	router := NewRouter()

	// Set up routes
	router.HandleFunc("GET", "/home", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("POST", "/user/:id", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("GET", "/files/*", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("DELETE", "/user/:id", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("GET", "/deep/nested/path", func(w http.ResponseWriter, r *http.Request) {})

	// Define test cases for route lookups
	testCases := []struct {
		method string
		path   string
	}{
		{"GET", "/home"},              // Static route
		{"POST", "/user/123"},         // Dynamic route
		{"GET", "/files/images/img1"}, // Wildcard route
		{"GET", "/deep/nested/path"},  // Deep static route
		{"DELETE", "/user/456"},       // Dynamic with different method
		//
		{"GET", "/home"},              // Static route
		{"POST", "/user/123"},         // Dynamic route
		{"GET", "/files/images/img1"}, // Wildcard route
		{"GET", "/deep/nested/path"},  // Deep static route
		{"DELETE", "/user/456"},       // Dynamic with different method
		//
		{"GET", "/user/456"},       // method not allowed
		{"POST", "/not-found/456"}, // path not found
	}

	for _, tc := range testCases {
		b.Run(tc.method+" "+tc.path, func(b *testing.B) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()

			for i := 0; i < b.N; i++ {
				router.ServeHTTP(rec, req)
			}
		})
	}
}
