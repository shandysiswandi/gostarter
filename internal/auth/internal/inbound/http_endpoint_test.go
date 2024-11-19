package inbound

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/stretchr/testify/assert"
)

func Test_httpEndpoint_Login(t *testing.T) {
	type args struct {
		ctx context.Context
		r   func() *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr error
		mockFn  func(a args) *httpEndpoint
	}{
		{
			name: "ErrorDecodeBody",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString("fake request")
					return httptest.NewRequest(http.MethodPost, "/auth/login", body)
				},
			},
			want:    nil,
			wantErr: goerror.NewInvalidFormat("invalid request body"),
			mockFn: func(a args) *httpEndpoint {

				return &httpEndpoint{}
			},
		},
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString(`{"email":"email","password":"password"}`)
					return httptest.NewRequest(http.MethodPost, "/auth/login", body)
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *httpEndpoint {
				loginMock := new(mockz.MockLogin)

				in := domain.LoginInput{Email: "email", Password: "password"}
				loginMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					loginUC: loginMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString(`{"email":"email","password":"password"}`)
					return httptest.NewRequest(http.MethodPost, "/auth/login", body)
				},
			},
			want: LoginResponse{
				AccessToken:      "access_token",
				RefreshToken:     "refresh_token",
				AccessExpiresIn:  10,
				RefreshExpiresIn: 20,
			},
			wantErr: nil,
			mockFn: func(a args) *httpEndpoint {
				loginMock := new(mockz.MockLogin)

				in := domain.LoginInput{Email: "email", Password: "password"}
				loginMock.EXPECT().
					Call(a.ctx, in).
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
			e := tt.mockFn(tt.args)
			got, err := e.Login(tt.args.ctx, tt.args.r())
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_httpEndpoint_Register(t *testing.T) {
	type args struct {
		ctx context.Context
		r   func() *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr error
		mockFn  func(a args) *httpEndpoint
	}{
		{
			name: "ErrorDecodeBody",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString("fake request")
					return httptest.NewRequest(http.MethodPost, "/auth/register", body)
				},
			},
			want:    nil,
			wantErr: goerror.NewInvalidFormat("invalid request body"),
			mockFn: func(a args) *httpEndpoint {

				return &httpEndpoint{}
			},
		},
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString(`{"email":"email","password":"password"}`)
					return httptest.NewRequest(http.MethodPost, "/auth/register", body)
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *httpEndpoint {
				registerMock := new(mockz.MockRegister)

				in := domain.RegisterInput{Email: "email", Password: "password"}
				registerMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					registerUC: registerMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString(`{"email":"email","password":"password"}`)
					return httptest.NewRequest(http.MethodPost, "/auth/register", body)
				},
			},
			want:    RegisterResponse{Email: "email"},
			wantErr: nil,
			mockFn: func(a args) *httpEndpoint {
				registerMock := new(mockz.MockRegister)

				in := domain.RegisterInput{Email: "email", Password: "password"}
				registerMock.EXPECT().
					Call(a.ctx, in).
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
			e := tt.mockFn(tt.args)
			got, err := e.Register(tt.args.ctx, tt.args.r())
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_httpEndpoint_RefreshToken(t *testing.T) {
	type args struct {
		ctx context.Context
		r   func() *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr error
		mockFn  func(a args) *httpEndpoint
	}{
		{
			name: "ErrorDecodeBody",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString("fake request")
					return httptest.NewRequest(http.MethodPost, "/auth/refresh-token", body)
				},
			},
			want:    nil,
			wantErr: goerror.NewInvalidFormat("invalid request body"),
			mockFn: func(a args) *httpEndpoint {

				return &httpEndpoint{}
			},
		},
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString(`{"refresh_token":"token"}`)
					return httptest.NewRequest(http.MethodPost, "/auth/refresh-token", body)
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *httpEndpoint {
				rtMock := new(mockz.MockRefreshToken)

				in := domain.RefreshTokenInput{RefreshToken: "token"}
				rtMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					refreshTokenUC: rtMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString(`{"refresh_token":"token"}`)
					return httptest.NewRequest(http.MethodPost, "/auth/refresh-token", body)
				},
			},
			want: RefreshTokenResponse{
				AccessToken:      "access_token",
				RefreshToken:     "refresh_token",
				AccessExpiresIn:  10,
				RefreshExpiresIn: 20,
			},
			wantErr: nil,
			mockFn: func(a args) *httpEndpoint {
				rtMock := new(mockz.MockRefreshToken)

				in := domain.RefreshTokenInput{RefreshToken: "token"}
				rtMock.EXPECT().
					Call(a.ctx, in).
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
			e := tt.mockFn(tt.args)
			got, err := e.RefreshToken(tt.args.ctx, tt.args.r())
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_httpEndpoint_ForgotPassword(t *testing.T) {
	type args struct {
		ctx context.Context
		r   func() *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr error
		mockFn  func(a args) *httpEndpoint
	}{
		{
			name: "ErrorDecodeBody",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString("fake request")
					return httptest.NewRequest(http.MethodPost, "/auth/forgot-password", body)
				},
			},
			want:    nil,
			wantErr: goerror.NewInvalidFormat("invalid request body"),
			mockFn: func(a args) *httpEndpoint {

				return &httpEndpoint{}
			},
		},
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString(`{"email":"email"}`)
					return httptest.NewRequest(http.MethodPost, "/auth/forgot-password", body)
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *httpEndpoint {
				fpMock := new(mockz.MockForgotPassword)

				in := domain.ForgotPasswordInput{Email: "email"}
				fpMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					forgotPasswordUC: fpMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString(`{"email":"email"}`)
					return httptest.NewRequest(http.MethodPost, "/auth/forgot-password", body)
				},
			},
			want: ForgotPasswordResponse{
				Email:   "email",
				Message: "message",
			},
			wantErr: nil,
			mockFn: func(a args) *httpEndpoint {
				fpMock := new(mockz.MockForgotPassword)

				in := domain.ForgotPasswordInput{Email: "email"}
				fpMock.EXPECT().
					Call(a.ctx, in).
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
			e := tt.mockFn(tt.args)
			got, err := e.ForgotPassword(tt.args.ctx, tt.args.r())
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_httpEndpoint_ResetPassword(t *testing.T) {
	type args struct {
		ctx context.Context
		r   func() *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr error
		mockFn  func(a args) *httpEndpoint
	}{
		{
			name: "ErrorDecodeBody",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString("fake request")
					return httptest.NewRequest(http.MethodPost, "/auth/reset-password", body)
				},
			},
			want:    nil,
			wantErr: goerror.NewInvalidFormat("invalid request body"),
			mockFn: func(a args) *httpEndpoint {

				return &httpEndpoint{}
			},
		},
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString(`{"token":"token","password":"password"}`)
					return httptest.NewRequest(http.MethodPost, "/auth/reset-password", body)
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *httpEndpoint {
				rpMock := new(mockz.MockResetPassword)

				in := domain.ResetPasswordInput{Token: "token", Password: "password"}
				rpMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					resetPasswordUC: rpMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				r: func() *http.Request {
					body := bytes.NewBufferString(`{"token":"token","password":"password"}`)
					return httptest.NewRequest(http.MethodPost, "/auth/reset-password", body)
				},
			},
			want:    ResetPasswordResponse{Message: "message"},
			wantErr: nil,
			mockFn: func(a args) *httpEndpoint {
				rpMock := new(mockz.MockResetPassword)

				in := domain.ResetPasswordInput{Token: "token", Password: "password"}
				rpMock.EXPECT().
					Call(a.ctx, in).
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
			e := tt.mockFn(tt.args)
			got, err := e.ResetPassword(tt.args.ctx, tt.args.r())
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
