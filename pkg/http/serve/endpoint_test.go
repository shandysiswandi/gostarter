package serve

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEndpoint_ServeHTTP(t *testing.T) {
	type fields struct {
		h           func(context.Context, *http.Request) (any, error)
		mws         []func(http.Handler) http.Handler
		resultCodec func(context.Context, http.ResponseWriter, any) error
		errorCodec  func(context.Context, http.ResponseWriter, error)
	}
	tests := []struct {
		name       string
		fields     fields
		statusCode int
	}{
		{
			name: "ErrorFromHandler",
			fields: fields{
				h: func(context.Context, *http.Request) (any, error) {
					return nil, http.ErrBodyNotAllowed
				},
				mws: []func(http.Handler) http.Handler{func(h http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						h.ServeHTTP(w, r)
					})
				}},
				resultCodec: func(context.Context, http.ResponseWriter, any) error {
					return nil
				},
				errorCodec: func(_ context.Context, w http.ResponseWriter, _ error) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"error":"internal server error"}`))
				},
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "ErrorFromResultCoder",
			fields: fields{
				h: func(context.Context, *http.Request) (any, error) {
					return nil, nil
				},
				resultCodec: func(context.Context, http.ResponseWriter, any) error {
					return http.ErrBodyNotAllowed
				},
				errorCodec: func(_ context.Context, w http.ResponseWriter, _ error) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"error":"internal server error"}`))
				},
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Success",
			fields: fields{
				h:           func(context.Context, *http.Request) (any, error) { return nil, nil },
				resultCodec: func(context.Context, http.ResponseWriter, any) error { return nil },
			},
			statusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &Endpoint{
				h:           tt.fields.h,
				mws:         tt.fields.mws,
				resultCodec: tt.fields.resultCodec,
				errorCodec:  tt.fields.errorCodec,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()
			e.ServeHTTP(rr, req)

			if tt.statusCode != rr.Code {
				t.Fatalf("expected status code:%d, got %d", tt.statusCode, rr.Code)
			}
		})
	}
}
