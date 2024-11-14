package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	libjwt "github.com/golang-jwt/jwt/v4"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/jwt/mocker"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	tests := []struct {
		name            string
		jwte            func() jwt.JWT
		handlerFunc     http.HandlerFunc
		expectedStatus  int
		expectedMessage string
		mockFn          func(r *http.Request)
	}{
		{
			name: "ErrorNoHeaderAuthorization",
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

			handler := JWT(tt.jwte(), "gostarter.access.token")(tt.handlerFunc)

			req := httptest.NewRequest(http.MethodGet, "/", nil)

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
