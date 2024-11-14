package serve

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type responseTest struct{}

func (responseTest) StatusCode() int { return http.StatusCreated }

type errorTest struct{}

func (errorTest) Error() string   { return "error" }
func (errorTest) StatusCode() int { return http.StatusNotFound }

func TestServe(t *testing.T) {
	tests := []struct {
		name       string
		handler    func(ctx context.Context, r *http.Request) (any, error)
		statusCode int
	}{
		{
			name: "SuccessWithCustomStatusCode",
			handler: func(ctx context.Context, r *http.Request) (any, error) {
				return responseTest{}, nil
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "Success",
			handler: func(ctx context.Context, r *http.Request) (any, error) {
				return "ok", nil
			},
			statusCode: http.StatusOK,
		},
		{
			name: "SuccessWithoutContent",
			handler: func(ctx context.Context, r *http.Request) (any, error) {
				return nil, nil
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "ErrorFromDecodeResult",
			handler: func(ctx context.Context, r *http.Request) (any, error) {
				return make(chan string), nil
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "ErrorFromHandler",
			handler: func(ctx context.Context, r *http.Request) (any, error) {
				return nil, errorTest{}
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := New(
				WithErrorCodec(defaultErrorCodec),
				WithResultCodec(defaultResultCodec),
				WithMiddlewares(func(h http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						h.ServeHTTP(w, r)
					})
				}),
			)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()
			endpoint := s.Endpoint(tt.handler)
			endpoint.ServeHTTP(rr, req)

			if tt.statusCode != rr.Code {
				t.Errorf("expected status %d, got %d", tt.statusCode, rr.Code)
			}
		})
	}
}
