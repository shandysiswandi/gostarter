package framework

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/stretchr/testify/assert"
)

type testResult struct{}

func (testResult) StatusCode() int { return 200 }

func TestNewRouter(t *testing.T) {
	tests := []struct {
		name string
		want *Router
	}{
		{
			name: "Success",
			want: NewRouter(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewRouter()
			assert.NotNil(t, got.hr)
			assert.NotNil(t, got.resultCodec)
			assert.NotNil(t, got.errorCodec)
		})
	}
}

func TestRouter_Endpoint(t *testing.T) {
	type args struct {
		method string
		path   string
		h      Handler
		body   io.Reader
		mws    []Middleware
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		mockFn         func(a args) *Router
	}{
		{
			name: "Error",
			args: args{
				method: "GET",
				path:   "/",
				h: func(Context) (any, error) {
					return nil, goerror.NewServerInternal(assert.AnError)
				},
				body: nil,
				mws:  nil,
			},
			wantStatusCode: 500,
			mockFn: func(a args) *Router {
				return NewRouter()
			},
		},
		{
			name: "NoContent",
			args: args{
				method: "GET",
				path:   "/users/:id",
				h: func(Context) (any, error) {
					return nil, nil
				},
				body: nil,
				mws:  nil,
			},
			wantStatusCode: 204,
			mockFn: func(a args) *Router {
				return NewRouter()
			},
		},
		{
			name: "Success",
			args: args{
				method: "GET",
				path:   "/users/:id",
				h: func(Context) (any, error) {
					return testResult{}, nil
				},
				body: nil,
				mws:  nil,
			},
			wantStatusCode: 200,
			mockFn: func(a args) *Router {
				return NewRouter()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := tt.mockFn(tt.args)
			r.Endpoint(tt.args.method, tt.args.path, tt.args.h, tt.args.mws...)
			req := httptest.NewRequest(tt.args.method, tt.args.path, tt.args.body)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, tt.wantStatusCode, resp.Code)
		})
	}
}

func TestRouter_HandleFunc(t *testing.T) {
	type args struct {
		method string
		path   string
		h      http.HandlerFunc
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		mockFn         func() *Router
	}{
		{
			name: "Error",
			args: args{
				method: "GET",
				path:   "/",
				h:      defaultNotFound,
			},
			wantStatusCode: 404,
			mockFn: func() *Router {
				return NewRouter()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := tt.mockFn()

			r.HandleFunc(tt.args.method, tt.args.path, tt.args.h)

			req := httptest.NewRequest(tt.args.method, tt.args.path, nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, tt.wantStatusCode, resp.Code)
		})
	}
}

func TestRouter_Handler(t *testing.T) {
	type args struct {
		method string
		path   string
		h      http.Handler
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		mockFn         func() *Router
	}{
		{
			name: "Error",
			args: args{
				method: "GET",
				path:   "/",
				h:      http.HandlerFunc(defaultMethodNotAllowed),
			},
			wantStatusCode: 405,
			mockFn: func() *Router {
				return NewRouter()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := tt.mockFn()

			r.Handler(tt.args.method, tt.args.path, tt.args.h)

			req := httptest.NewRequest(tt.args.method, tt.args.path, nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, tt.wantStatusCode, resp.Code)
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
