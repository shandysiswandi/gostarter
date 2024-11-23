package framework

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/stretchr/testify/assert"
)

type testObject struct{}

func (testObject) StatusCode() int {
	return 200
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Router
	}{
		{
			name: "Success",
			want: &Router{
				router:           httprouter.New(),
				notFound:         http.HandlerFunc(defaultNotFound),
				methodNotAllowed: http.HandlerFunc(defaultMethodNotAllowed),
				resultCodec:      defaultResultCodec,
				errorCodec:       defaultErrorCodec,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := New()
			assert.Equal(t, tt.want.router, got.router)
			assert.NotNil(t, got.notFound)
			assert.NotNil(t, got.methodNotAllowed)
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
		mws    []Middleware
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(a args) *Router
	}{
		{
			name: "Success",
			args: args{
				method: http.MethodGet,
				path:   "/",
				h: func(Context) (any, error) {
					return nil, assert.AnError
				},
				mws: []Middleware{},
			},
			mockFn: func(a args) *Router {
				return &Router{
					router: httprouter.New(),
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := tt.mockFn(tt.args)
			r.Endpoint(tt.args.method, tt.args.path, tt.args.h, tt.args.mws...)
		})
	}
}

func TestRouter_Native(t *testing.T) {
	type args struct {
		method string
		path   string
		h      http.Handler
		mws    []Middleware
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(a args) *Router
	}{
		{
			name: "Success",
			args: args{
				method: http.MethodGet,
				path:   "/",
				h:      http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				mws: []Middleware{
					func(h http.Handler) http.Handler {
						return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							h.ServeHTTP(w, r)
						})
					},
				},
			},
			mockFn: func(a args) *Router {
				return &Router{
					router: httprouter.New(),
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := tt.mockFn(tt.args)
			r.Native(tt.args.method, tt.args.path, tt.args.h, tt.args.mws...)
		})
	}
}

func TestRouter_ServeHTTP(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(a args) *Router
	}{
		{
			name: "Success",
			args: args{
				w:   httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			mockFn: func(a args) *Router {
				return &Router{
					router: httprouter.New(),
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := tt.mockFn(tt.args)
			r.ServeHTTP(tt.args.w, tt.args.req)
		})
	}
}

func TestRouter_wrapHandler(t *testing.T) {
	tests := []struct {
		name            string
		h               Handler
		expectedStatus  int
		expectedMessage string
		mockFn          func() *Router
	}{
		{
			name: "Error",
			h: func(Context) (any, error) {
				return nil, goerror.NewServerInternal(assert.AnError)
			},
			expectedStatus:  500,
			expectedMessage: "{\"error\":\"assert.AnError general error for testing\"}\n",
			mockFn: func() *Router {
				return New()
			},
		},
		{
			name: "NoContent",
			h: func(Context) (any, error) {
				return nil, nil
			},
			expectedStatus:  204,
			expectedMessage: "",
			mockFn: func() *Router {
				return New()
			},
		},
		{
			name: "Success",
			h: func(Context) (any, error) {
				return testObject{}, nil
			},
			expectedStatus:  200,
			expectedMessage: "{}\n",
			mockFn: func() *Router {
				return New()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := tt.mockFn()
			handler := r.wrapHandler(tt.h)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedMessage, rr.Body.String())
		})
	}
}

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
			},
		},
		{
			name: "Success",
			args: args{
				w:    httptest.NewRecorder(),
				data: "success",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			writeJSON(tt.args.w, tt.args.data)
		})
	}
}
