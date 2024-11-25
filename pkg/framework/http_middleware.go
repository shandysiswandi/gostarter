package framework

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/shandysiswandi/gostarter/pkg/jwt"
)

// Middleware defines the type for middleware functions.
type Middleware func(http.Handler) http.Handler

// Chain applies a sequence of Middleware functions to an http.Handler.
//
// The middleware functions are applied in the order they are provided, so the first middleware
// in the list will be the outermost wrapper around the handler, and the last middleware will be
// the innermost wrapper.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}

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
				log.Printf("panic because: %v\n", err)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")

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

func JWT(jwte jwt.JWT, audience string, skipPaths ...string) Middleware {
	mj := &middlewareJWT{
		jwte:      jwte,
		audience:  audience,
		skipPaths: skipPaths,
	}

	return mj.handle
}

type middlewareJWT struct {
	jwte      jwt.JWT
	audience  string
	skipPaths []string
}

func (mj *middlewareJWT) handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mj.shouldSkipPath(r.URL.Path) {
			h.ServeHTTP(w, r)

			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSON(w, map[string]string{"error": "authorization header missing"}, http.StatusUnauthorized)

			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeJSON(w, map[string]string{"error": "invalid format"}, http.StatusUnauthorized)

			return
		}

		clm, err := mj.jwte.Verify(strings.TrimPrefix(authHeader, "Bearer "))
		if errors.Is(err, jwt.ErrTokenExpired) {
			writeJSON(w, map[string]string{"error": "expired token"}, http.StatusUnauthorized)

			return
		}

		if err != nil {
			writeJSON(w, map[string]string{"error": "invalid token"}, http.StatusUnauthorized)

			return
		}

		if !clm.VerifyAudience(mj.audience, true) {
			writeJSON(w, map[string]string{"error": "invalid token audience"}, http.StatusUnauthorized)

			return
		}

		h.ServeHTTP(w, r.WithContext(jwt.SetClaim(r.Context(), clm)))
	})
}

func (mj *middlewareJWT) shouldSkipPath(path string) bool {
	for _, prefix := range mj.skipPaths {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}
