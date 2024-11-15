package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/shandysiswandi/gostarter/pkg/framework/httpserver"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
)

func JWT(jwte jwt.JWT, audience string, skipPaths ...string) httpserver.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(skipPaths) > 0 {
				for _, prefix := range skipPaths {
					if strings.HasPrefix(r.URL.Path, prefix) {
						h.ServeHTTP(w, r)

						return
					}
				}
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				jsonResponse(w, "authorization header missing")

				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				jsonResponse(w, "invalid format")

				return
			}

			clm, err := jwte.Verify(strings.TrimPrefix(authHeader, "Bearer "))
			if errors.Is(err, jwt.ErrTokenExpired) {
				jsonResponse(w, "expired token")

				return
			}

			if err != nil {
				jsonResponse(w, "invalid token")

				return
			}

			if !clm.VerifyAudience(audience, true) {
				jsonResponse(w, "invalid token audience")

				return
			}

			h.ServeHTTP(w, r.WithContext(jwt.SetClaim(r.Context(), clm)))
		})
	}
}

func jsonResponse(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": msg}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
