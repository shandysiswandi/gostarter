// Package middleware provides HTTP middleware utilities that can be used to enhance
// the functionality of HTTP handlers in a web application. This package includes
// middleware for recovering from panics in HTTP handlers.
package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

// Recovery is a middleware that recovers from any panics that occur during the execution
// of an HTTP handler and writes a 500 Internal Server Error response to the client.
//
// If a panic occurs, the middleware logs the panic message and stack trace,
// and sends a JSON-encoded error message with a status code of 500.
//
// The `Recovery` middleware is particularly useful for preventing server crashes
// due to unexpected errors, ensuring that the server continues to run and handle
// subsequent requests.
//
// Usage:
//
//	http.Handle("/some-endpoint", middleware.Recovery(yourHandler))
//
// The middleware should be the first in the chain to ensure it wraps the entire
// request handling process.
func Recovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			//nolint:err113,errorlint,errcheck // ignore for this
			// Recover from any panics and handle them appropriately.
			if err := recover(); err != nil && err != http.ErrAbortHandler {
				// Log the panic message.
				log.Printf("panic from: %v\n", err)

				// If the connection is not being upgraded, write a 500 status code.
				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
				}

				// Print the stack trace for debugging purposes.
				debug.PrintStack()

				// Send a default fallback response to the client.
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": http.StatusText(http.StatusInternalServerError),
				})
			}
		}()

		// Continue to the next handler in the chain.
		h.ServeHTTP(w, r)
	})
}
