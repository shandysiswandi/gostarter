package http

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

func recovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic from: %v\n", err)

				if err == http.ErrAbortHandler { //nolint:err113,errorlint // ignore for this
					return
				}

				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
				}

				debug.PrintStack()
				//nolint:errcheck // never error, maybe
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": http.StatusText(http.StatusInternalServerError),
				})
			}
		}()

		h.ServeHTTP(w, r)
	})
}
