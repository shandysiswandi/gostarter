package framework

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	libjwt "github.com/golang-jwt/jwt/v4"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/jwt/mocker"
	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	tests := []struct {
		name            string
		mws             []Middleware
		handlerFunc     http.HandlerFunc
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "Success",
			mws:  []Middleware{Recovery},
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

func TestRecovery(t *testing.T) {
	tests := []struct {
		name            string
		handlerFunc     http.HandlerFunc
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "No panic",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("OK"))
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "OK",
		},
		{
			name: "Panic occurs",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				panic("something went wrong")
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: `{"error":"Internal Server Error"}`,
		},
		{
			name: "Abort handler panic",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				panic(http.ErrAbortHandler)
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Wrap the handler with the Recovery middleware
			handler := Recovery(tt.handlerFunc)

			// Create a new HTTP request
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Serve the HTTP request
			handler.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Check the response body
			if !strings.Contains(rr.Body.String(), tt.expectedMessage) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedMessage, rr.Body.String())
			}

			// If testing the panic case, check that the response is JSON formatted
			if tt.expectedStatus == http.StatusInternalServerError {
				var responseBody map[string]string
				if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
					t.Errorf("response body is not valid JSON: %v", err)
				}

				if responseBody["error"] != "Internal Server Error" {
					t.Errorf("expected error message %q, got %q", "Internal Server Error", responseBody["error"])
				}
			}
		})
	}
}

func TestJWT(t *testing.T) {
	tests := []struct {
		name            string
		opts            []string
		path            string
		jwte            func() jwt.JWT
		handlerFunc     http.HandlerFunc
		expectedStatus  int
		expectedMessage string
		mockFn          func(r *http.Request)
	}{
		{
			name: "SkipPath",
			path: "/skip",
			opts: []string{"/skip"},
			jwte: func() jwt.JWT {
				return &mocker.MockJWT{}
			},
			mockFn: nil,
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "OK",
		},
		{
			name: "ErrorNoHeaderAuthorization",
			path: "/continue",
			jwte: func() jwt.JWT {
				return &mocker.MockJWT{}
			},
			mockFn:          nil,
			handlerFunc:     nil,
			expectedStatus:  http.StatusUnauthorized,
			expectedMessage: "{\"error\":\"authorization header missing\"}\n",
		},
		{
			name: "ErrorNoPrefixBearer",
			path: "/continue",
			jwte: func() jwt.JWT {
				return &mocker.MockJWT{}
			},
			mockFn: func(r *http.Request) {
				r.Header.Set("Authorization", "zzz")
			},
			handlerFunc:     nil,
			expectedStatus:  http.StatusUnauthorized,
			expectedMessage: "{\"error\":\"invalid format\"}\n",
		},
		{
			name: "ErrorTokenExpired",
			path: "/continue",
			jwte: func() jwt.JWT {
				mjwt := &mocker.MockJWT{}

				mjwt.EXPECT().Verify("ay").Return(nil, jwt.ErrTokenExpired)

				return mjwt
			},
			mockFn: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer ay")
			},
			handlerFunc:     nil,
			expectedStatus:  http.StatusUnauthorized,
			expectedMessage: "{\"error\":\"expired token\"}\n",
		},
		{
			name: "ErrorJWTVerify",
			path: "/continue",
			jwte: func() jwt.JWT {
				mjwt := &mocker.MockJWT{}

				mjwt.EXPECT().Verify("ay").Return(nil, assert.AnError)

				return mjwt
			},
			mockFn: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer ay")
			},
			handlerFunc:     nil,
			expectedStatus:  http.StatusUnauthorized,
			expectedMessage: "{\"error\":\"invalid token\"}\n",
		},
		{
			name: "ErrorInvalidAudience",
			path: "/continue",
			jwte: func() jwt.JWT {
				mjwt := &mocker.MockJWT{}

				mjwt.EXPECT().Verify("ay").Return(&jwt.Claim{}, nil)

				return mjwt
			},
			mockFn: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer ay")
			},
			handlerFunc:     nil,
			expectedStatus:  http.StatusUnauthorized,
			expectedMessage: "{\"error\":\"invalid token audience\"}\n",
		},
		{
			name: "Success",
			path: "/continue",
			jwte: func() jwt.JWT {
				mjwt := &mocker.MockJWT{}

				mjwt.EXPECT().Verify("ay").Return(&jwt.Claim{
					RegisteredClaims: libjwt.RegisteredClaims{
						Audience: []string{"gostarter.access.token"},
					},
				}, nil)

				return mjwt
			},
			mockFn: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer ay")
			},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := JWT(tt.jwte(), "gostarter.access.token", tt.opts...)(tt.handlerFunc)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)

			if tt.mockFn != nil {
				tt.mockFn(req)
			}

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedMessage, rr.Body.String())
		})
	}
}
