package inbound

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/framework"
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
			wantErr: goerror.NewInvalidFormat("Request payload malformed"),
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "auth.inbound.http.Login")
				defer span.End()

				return &httpEndpoint{
					telemetry: tel,
				}
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
				loginMock := mockz.NewMockLogin(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "auth.inbound.http.Login")
				defer span.End()

				in := domain.LoginInput{
					Email:    "email",
					Password: "password",
				}
				loginMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					telemetry: tel,
					loginUC:   loginMock,
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
				loginMock := mockz.NewMockLogin(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "auth.inbound.http.Login")
				defer span.End()

				in := domain.LoginInput{
					Email:    "email",
					Password: "password",
				}
				out := &domain.LoginOutput{
					AccessToken:      "access_token",
					RefreshToken:     "refresh_token",
					AccessExpiresIn:  10,
					RefreshExpiresIn: 20,
				}
				loginMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &httpEndpoint{
					telemetry: tel,
					loginUC:   loginMock,
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
			wantErr: goerror.NewInvalidFormat("Request payload malformed"),
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "auth.inbound.http.Register")
				defer span.End()

				return &httpEndpoint{
					telemetry: tel,
				}
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
				registerMock := mockz.NewMockRegister(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "auth.inbound.http.Register")
				defer span.End()

				in := domain.RegisterInput{
					Name:     "fullname",
					Email:    "email",
					Password: "password",
				}
				registerMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					telemetry:  tel,
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
				registerMock := mockz.NewMockRegister(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "auth.inbound.http.Register")
				defer span.End()

				in := domain.RegisterInput{
					Name:     "fullname",
					Email:    "email",
					Password: "password",
				}
				out := &domain.RegisterOutput{Email: "email"}
				registerMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &httpEndpoint{
					telemetry:  tel,
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

func Test_httpEndpoint_Verify(t *testing.T) {
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
				c := framework.NewTestContext(http.MethodPost, "/auth/verify", body)
				return c.Build()
			},
			want:    nil,
			wantErr: goerror.NewInvalidFormat("Request payload malformed"),
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "auth.inbound.http.Verify")
				defer span.End()

				return &httpEndpoint{
					telemetry: tel,
				}
			},
		},
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"code":"code22","email":"email"}`)
				c := framework.NewTestContext(http.MethodPost, "/auth/verify", body)
				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				verifyMock := mockz.NewMockVerify(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "auth.inbound.http.Verify")
				defer span.End()

				in := domain.VerifyInput{
					Code:  "code22",
					Email: "email",
				}
				verifyMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					telemetry: tel,
					verifyUC:  verifyMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"code":"code22","email":"email"}`)
				c := framework.NewTestContext(http.MethodPost, "/auth/verify", body)
				return c.Build()
			},
			want: VerifyResponse{
				Email:    "email",
				VerifyAt: "0001-01-01T00:00:00Z",
			},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				verifyMock := mockz.NewMockVerify(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "auth.inbound.http.Verify")
				defer span.End()

				in := domain.VerifyInput{
					Code:  "code22",
					Email: "email",
				}
				out := &domain.VerifyOutput{
					Email:    "email",
					VerifyAt: time.Time{},
				}
				verifyMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &httpEndpoint{
					telemetry: tel,
					verifyUC:  verifyMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.Verify(c)
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
			wantErr: goerror.NewInvalidFormat("Request payload malformed"),
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "auth.inbound.http.RefreshToken")
				defer span.End()

				return &httpEndpoint{
					telemetry: tel,
				}
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
				rtMock := mockz.NewMockRefreshToken(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "auth.inbound.http.RefreshToken")
				defer span.End()

				in := domain.RefreshTokenInput{RefreshToken: "token"}
				rtMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					telemetry:      tel,
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
				rtMock := mockz.NewMockRefreshToken(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "auth.inbound.http.RefreshToken")
				defer span.End()

				in := domain.RefreshTokenInput{RefreshToken: "token"}
				out := &domain.RefreshTokenOutput{
					AccessToken:      "access_token",
					RefreshToken:     "refresh_token",
					AccessExpiresIn:  10,
					RefreshExpiresIn: 20,
				}
				rtMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &httpEndpoint{
					telemetry:      tel,
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
			wantErr: goerror.NewInvalidFormat("Request payload malformed"),
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "auth.inbound.http.ForgotPassword")
				defer span.End()

				return &httpEndpoint{
					telemetry: tel,
				}
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
				fpMock := mockz.NewMockForgotPassword(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "auth.inbound.http.ForgotPassword")
				defer span.End()

				in := domain.ForgotPasswordInput{Email: "email"}
				fpMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					telemetry:        tel,
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
				fpMock := mockz.NewMockForgotPassword(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "auth.inbound.http.ForgotPassword")
				defer span.End()

				in := domain.ForgotPasswordInput{Email: "email"}
				out := &domain.ForgotPasswordOutput{
					Email:   "email",
					Message: "message",
				}
				fpMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &httpEndpoint{
					telemetry:        tel,
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
			wantErr: goerror.NewInvalidFormat("Request payload malformed"),
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "auth.inbound.http.ForgotPassword")
				defer span.End()

				return &httpEndpoint{
					telemetry: tel,
				}
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
				rpMock := mockz.NewMockResetPassword(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "auth.inbound.http.ForgotPassword")
				defer span.End()

				in := domain.ResetPasswordInput{
					Token:    "token",
					Password: "password",
				}
				rpMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					telemetry:       tel,
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
				rpMock := mockz.NewMockResetPassword(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "auth.inbound.http.ForgotPassword")
				defer span.End()

				in := domain.ResetPasswordInput{
					Token:    "token",
					Password: "password",
				}
				out := &domain.ResetPasswordOutput{Message: "message"}
				rpMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &httpEndpoint{
					telemetry:       tel,
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
