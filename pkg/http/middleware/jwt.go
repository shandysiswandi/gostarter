package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/shandysiswandi/gostarter/pkg/jwt"
)

func JWT(jwte jwt.JWT) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			if !clm.VerifyAudience("gostarter.access.token", true) {
				jsonResponse(w, "invalid token audience")

				return
			}

			h.ServeHTTP(w, r.WithContext(jwt.SetClaimToContext(r.Context(), clm)))
		})
	}
}

func jsonResponse(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	//nolint:errcheck,errchkjson // ignore for this, it never error
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
