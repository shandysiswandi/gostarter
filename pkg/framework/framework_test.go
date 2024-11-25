package framework

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_defaultNotFound(t *testing.T) {
	tests := []struct {
		name        string
		wantMessage string
		wantCode    int
	}{
		{
			name:        "Success",
			wantMessage: "{\"error\":\"endpoint not found\"}\n",
			wantCode:    404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			http.HandlerFunc(defaultNotFound).ServeHTTP(rr, req)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Equal(t, tt.wantMessage, rr.Body.String())
		})
	}
}

func Test_defaultMethodNotAllowed(t *testing.T) {
	tests := []struct {
		name        string
		wantMessage string
		wantCode    int
	}{
		{
			name:        "Success",
			wantMessage: "{\"error\":\"method not allowed\"}\n",
			wantCode:    405,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			http.HandlerFunc(defaultMethodNotAllowed).ServeHTTP(rr, req)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Equal(t, tt.wantMessage, rr.Body.String())
		})
	}
}

func Test_writeJSON(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		data any
		code int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Error",
			args: args{
				w:    httptest.NewRecorder(),
				data: make(chan string),
				code: 500,
			},
		},
		{
			name: "Success",
			args: args{
				w:    httptest.NewRecorder(),
				data: "success",
				code: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			writeJSON(tt.args.w, tt.args.data, tt.args.code)
		})
	}
}
