package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecovery(t *testing.T) {
	tests := []struct {
		name            string
		handlerFunc     http.HandlerFunc
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "No panic",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("OK"))
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "OK",
		},
		{
			name: "Panic occurs",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				panic("something went wrong")
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: `{"error":"Internal Server Error"}`,
		},
		{
			name: "Abort handler panic",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				panic(http.ErrAbortHandler)
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Wrap the handler with the Recovery middleware
			handler := Recovery(tt.handlerFunc)

			// Create a new HTTP request
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Serve the HTTP request
			handler.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Check the response body
			if !strings.Contains(rr.Body.String(), tt.expectedMessage) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedMessage, rr.Body.String())
			}

			// If testing the panic case, check that the response is JSON formatted
			if tt.expectedStatus == http.StatusInternalServerError {
				var responseBody map[string]string
				if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
					t.Errorf("response body is not valid JSON: %v", err)
				}

				if responseBody["error"] != "Internal Server Error" {
					t.Errorf("expected error message %q, got %q", "Internal Server Error", responseBody["error"])
				}
			}
		})
	}
}
