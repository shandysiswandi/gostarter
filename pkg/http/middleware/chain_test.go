package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	tests := []struct {
		name            string
		mws             []func(http.Handler) http.Handler
		handlerFunc     http.HandlerFunc
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "Success",
			mws:  []func(http.Handler) http.Handler{Recovery},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chain := Chain(tt.handlerFunc, tt.mws...)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			chain.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedMessage, rr.Body.String())
		})
	}
}
