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
			mj.jsonResponse(w, "authorization header missing")

			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			mj.jsonResponse(w, "invalid format")

			return
		}

		clm, err := mj.jwte.Verify(strings.TrimPrefix(authHeader, "Bearer "))
		if errors.Is(err, jwt.ErrTokenExpired) {
			mj.jsonResponse(w, "expired token")

			return
		}

		if err != nil {
			mj.jsonResponse(w, "invalid token")

			return
		}

		if !clm.VerifyAudience(mj.audience, true) {
			mj.jsonResponse(w, "invalid token audience")

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

func (mj *middlewareJWT) jsonResponse(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": msg}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
