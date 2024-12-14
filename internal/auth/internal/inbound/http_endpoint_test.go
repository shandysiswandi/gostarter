package inbound

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/stretchr/testify/assert"
)

func Test_httpEndpoint_Login(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorDecodeBody",
			c: func() framework.Context {
				body := bytes.NewBufferString("fake request")
				c := framework.NewTestContext(http.MethodPost, "/auth/login", body)
				return c.Build()
			},
			want:    nil,
			wantErr: goerror.NewInvalidFormat("invalid request body"),
			mockFn: func(ctx context.Context) *httpEndpoint {
				return &httpEndpoint{}
			},
		},
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"email":"email","password":"password"}`)
				c := framework.NewTestContext(http.MethodPost, "/auth/login", body)
				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				loginMock := new(mockz.MockLogin)

				in := domain.LoginInput{Email: "email", Password: "password"}
				loginMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					loginUC: loginMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"email":"email","password":"password"}`)
				c := framework.NewTestContext(http.MethodPost, "/auth/login", body)
				return c.Build()
			},
			want: LoginResponse{
				AccessToken:      "access_token",
				RefreshToken:     "refresh_token",
				AccessExpiresIn:  10,
				RefreshExpiresIn: 20,
			},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				loginMock := new(mockz.MockLogin)

				in := domain.LoginInput{Email: "email", Password: "password"}
				loginMock.EXPECT().
					Call(ctx, in).
					Return(&domain.LoginOutput{
						AccessToken:      "access_token",
						RefreshToken:     "refresh_token",
						AccessExpiresIn:  10,
						RefreshExpiresIn: 20,
					}, nil)

				return &httpEndpoint{
					loginUC: loginMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.Login(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_httpEndpoint_Register(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorDecodeBody",
			c: func() framework.Context {
				body := bytes.NewBufferString("fake request")
				c := framework.NewTestContext(http.MethodPost, "/auth/register", body)
				return c.Build()
			},
			want:    nil,
			wantErr: goerror.NewInvalidFormat("invalid request body"),
			mockFn: func(ctx context.Context) *httpEndpoint {

				return &httpEndpoint{}
			},
		},
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"name":"fullname","email":"email","password":"password"}`)
				c := framework.NewTestContext(http.MethodPost, "/auth/register", body)
				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				registerMock := new(mockz.MockRegister)

				in := domain.RegisterInput{Name: "fullname", Email: "email", Password: "password"}
				registerMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					registerUC: registerMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"name":"fullname","email":"email","password":"password"}`)
				c := framework.NewTestContext(http.MethodPost, "/auth/register", body)
				return c.Build()
			},
			want:    RegisterResponse{Email: "email"},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				registerMock := new(mockz.MockRegister)

				in := domain.RegisterInput{Name: "fullname", Email: "email", Password: "password"}
				registerMock.EXPECT().
					Call(ctx, in).
					Return(&domain.RegisterOutput{Email: "email"}, nil)

				return &httpEndpoint{
					registerUC: registerMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.Register(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_httpEndpoint_RefreshToken(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorDecodeBody",
			c: func() framework.Context {
				body := bytes.NewBufferString("fake request")
				c := framework.NewTestContext(http.MethodPost, "/auth/refresh-token", body)
				return c.Build()
			},
			want:    nil,
			wantErr: goerror.NewInvalidFormat("invalid request body"),
			mockFn: func(ctx context.Context) *httpEndpoint {

				return &httpEndpoint{}
			},
		},
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"refresh_token":"token"}`)
				c := framework.NewTestContext(http.MethodPost, "/auth/refresh-token", body)
				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				rtMock := new(mockz.MockRefreshToken)

				in := domain.RefreshTokenInput{RefreshToken: "token"}
				rtMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					refreshTokenUC: rtMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"refresh_token":"token"}`)
				c := framework.NewTestContext(http.MethodPost, "/auth/refresh-token", body)
				return c.Build()
			},
			want: RefreshTokenResponse{
				AccessToken:      "access_token",
				RefreshToken:     "refresh_token",
				AccessExpiresIn:  10,
				RefreshExpiresIn: 20,
			},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				rtMock := new(mockz.MockRefreshToken)

				in := domain.RefreshTokenInput{RefreshToken: "token"}
				rtMock.EXPECT().
					Call(ctx, in).
					Return(&domain.RefreshTokenOutput{
						AccessToken:      "access_token",
						RefreshToken:     "refresh_token",
						AccessExpiresIn:  10,
						RefreshExpiresIn: 20,
					}, nil)

				return &httpEndpoint{
					refreshTokenUC: rtMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.RefreshToken(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_httpEndpoint_ForgotPassword(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorDecodeBody",
			c: func() framework.Context {
				body := bytes.NewBufferString("fake request")
				c := framework.NewTestContext(http.MethodPost, "/auth/forgot-password", body)
				return c.Build()
			},
			want:    nil,
			wantErr: goerror.NewInvalidFormat("invalid request body"),
			mockFn: func(ctx context.Context) *httpEndpoint {
				return &httpEndpoint{}
			},
		},
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"email":"email"}`)
				c := framework.NewTestContext(http.MethodPost, "/auth/forgot-password", body)
				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				fpMock := new(mockz.MockForgotPassword)

				in := domain.ForgotPasswordInput{Email: "email"}
				fpMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					forgotPasswordUC: fpMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"email":"email"}`)
				c := framework.NewTestContext(http.MethodPost, "/auth/forgot-password", body)
				return c.Build()
			},
			want: ForgotPasswordResponse{
				Email:   "email",
				Message: "message",
			},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				fpMock := new(mockz.MockForgotPassword)

				in := domain.ForgotPasswordInput{Email: "email"}
				fpMock.EXPECT().
					Call(ctx, in).
					Return(&domain.ForgotPasswordOutput{
						Email:   "email",
						Message: "message",
					}, nil)

				return &httpEndpoint{
					forgotPasswordUC: fpMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.ForgotPassword(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_httpEndpoint_ResetPassword(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorDecodeBody",
			c: func() framework.Context {
				body := bytes.NewBufferString("fake request")
				c := framework.NewTestContext(http.MethodPost, "/auth/reset-password", body)
				return c.Build()
			},
			want:    nil,
			wantErr: goerror.NewInvalidFormat("invalid request body"),
			mockFn: func(ctx context.Context) *httpEndpoint {

				return &httpEndpoint{}
			},
		},
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"token":"token","password":"password"}`)
				c := framework.NewTestContext(http.MethodPost, "/auth/reset-password", body)
				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				rpMock := new(mockz.MockResetPassword)

				in := domain.ResetPasswordInput{Token: "token", Password: "password"}
				rpMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					resetPasswordUC: rpMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"token":"token","password":"password"}`)
				c := framework.NewTestContext(http.MethodPost, "/auth/reset-password", body)
				return c.Build()
			},
			want:    ResetPasswordResponse{Message: "message"},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				rpMock := new(mockz.MockResetPassword)

				in := domain.ResetPasswordInput{Token: "token", Password: "password"}
				rpMock.EXPECT().
					Call(ctx, in).
					Return(&domain.ResetPasswordOutput{Message: "message"}, nil)

				return &httpEndpoint{
					resetPasswordUC: rpMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.ResetPassword(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
