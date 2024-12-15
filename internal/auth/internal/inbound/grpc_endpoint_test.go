package inbound

import (
	"context"
	"testing"

	pb "github.com/shandysiswandi/gostarter/api/gen-proto/auth"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/stretchr/testify/assert"
)

func TestGrpcEndpoint_Login(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.LoginRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.LoginResponse
		wantErr error
		mockFn  func(a args) *grpcEndpoint
	}{
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				req: &pb.LoginRequest{
					Email:    "email",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *grpcEndpoint {
				loginMock := mockz.NewMockLogin(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "auth.inbound.grpcEndpoint.Login")
				defer span.End()

				in := domain.LoginInput{
					Email:    a.req.Email,
					Password: a.req.Password,
				}
				loginMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
					telemetry: tel,
					loginUC:   loginMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				req: &pb.LoginRequest{
					Email:    "email",
					Password: "password",
				},
			},
			want: &pb.LoginResponse{
				AccessToken:      "access_token",
				RefreshToken:     "refresh_token",
				AccessExpiresIn:  10,
				RefreshExpiresIn: 20,
			},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				loginMock := mockz.NewMockLogin(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "auth.inbound.grpcEndpoint.Login")
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

				return &grpcEndpoint{
					telemetry: tel,
					loginUC:   loginMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := tt.mockFn(tt.args)
			got, err := g.Login(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGrpcEndpoint_Register(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.RegisterRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.RegisterResponse
		wantErr error
		mockFn  func(a args) *grpcEndpoint
	}{
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				req: &pb.RegisterRequest{
					Email:    "email",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *grpcEndpoint {
				registerMock := mockz.NewMockRegister(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "auth.inbound.grpcEndpoint.Register")
				defer span.End()

				in := domain.RegisterInput{
					Name:     a.req.Name,
					Email:    a.req.Email,
					Password: a.req.Password,
				}
				registerMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
					telemetry:  tel,
					registerUC: registerMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				req: &pb.RegisterRequest{
					Email:    "email",
					Password: "password",
				},
			},
			want:    &pb.RegisterResponse{Email: "email"},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				registerMock := mockz.NewMockRegister(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "auth.inbound.grpcEndpoint.Register")
				defer span.End()

				in := domain.RegisterInput{
					Name:     a.req.Name,
					Email:    a.req.Email,
					Password: a.req.Password,
				}
				out := &domain.RegisterOutput{Email: "email"}
				registerMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
					telemetry:  tel,
					registerUC: registerMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := tt.mockFn(tt.args)
			got, err := g.Register(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGrpcEndpoint_RefreshToken(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.RefreshTokenRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.RefreshTokenResponse
		wantErr error
		mockFn  func(a args) *grpcEndpoint
	}{
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				req: &pb.RefreshTokenRequest{RefreshToken: "token"},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *grpcEndpoint {
				refreshTokenMock := mockz.NewMockRefreshToken(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "auth.inbound.grpcEndpoint.RefreshToken")
				defer span.End()

				in := domain.RefreshTokenInput{RefreshToken: a.req.RefreshToken}
				refreshTokenMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
					telemetry:      tel,
					refreshTokenUC: refreshTokenMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				req: &pb.RefreshTokenRequest{RefreshToken: "token"},
			},
			want: &pb.RefreshTokenResponse{
				AccessToken:      "access_token",
				RefreshToken:     "refresh_token",
				AccessExpiresIn:  10,
				RefreshExpiresIn: 20,
			},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				refreshTokenMock := mockz.NewMockRefreshToken(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "auth.inbound.grpcEndpoint.RefreshToken")
				defer span.End()

				in := domain.RefreshTokenInput{RefreshToken: a.req.RefreshToken}
				out := &domain.RefreshTokenOutput{
					AccessToken:      "access_token",
					RefreshToken:     "refresh_token",
					AccessExpiresIn:  10,
					RefreshExpiresIn: 20,
				}
				refreshTokenMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
					telemetry:      tel,
					refreshTokenUC: refreshTokenMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := tt.mockFn(tt.args)
			got, err := g.RefreshToken(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGrpcEndpoint_ForgotPassword(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.ForgotPasswordRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.ForgotPasswordResponse
		wantErr error
		mockFn  func(a args) *grpcEndpoint
	}{
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				req: &pb.ForgotPasswordRequest{Email: "email"},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *grpcEndpoint {
				forgotPasswordMock := mockz.NewMockForgotPassword(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "auth.inbound.grpcEndpoint.ForgotPassword")
				defer span.End()

				in := domain.ForgotPasswordInput{Email: a.req.Email}
				forgotPasswordMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
					telemetry:        tel,
					forgotPasswordUC: forgotPasswordMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				req: &pb.ForgotPasswordRequest{Email: "email"},
			},
			want: &pb.ForgotPasswordResponse{
				Email:   "email",
				Message: "message",
			},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				forgotPasswordMock := mockz.NewMockForgotPassword(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "auth.inbound.grpcEndpoint.ForgotPassword")
				defer span.End()

				in := domain.ForgotPasswordInput{Email: a.req.Email}
				out := &domain.ForgotPasswordOutput{
					Email:   "email",
					Message: "message",
				}
				forgotPasswordMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
					telemetry:        tel,
					forgotPasswordUC: forgotPasswordMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := tt.mockFn(tt.args)
			got, err := g.ForgotPassword(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGrpcEndpoint_ResetPassword(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.ResetPasswordRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.ResetPasswordResponse
		wantErr error
		mockFn  func(a args) *grpcEndpoint
	}{
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				req: &pb.ResetPasswordRequest{
					Token:    "token",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *grpcEndpoint {
				resetPasswordMock := mockz.NewMockResetPassword(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "auth.inbound.grpcEndpoint.ResetPassword")
				defer span.End()

				in := domain.ResetPasswordInput{
					Token:    a.req.Token,
					Password: a.req.Password,
				}
				resetPasswordMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
					telemetry:       tel,
					resetPasswordUC: resetPasswordMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				req: &pb.ResetPasswordRequest{
					Token:    "token",
					Password: "password",
				},
			},
			want:    &pb.ResetPasswordResponse{Message: "message"},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				resetPasswordMock := mockz.NewMockResetPassword(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "auth.inbound.grpcEndpoint.ResetPassword")
				defer span.End()

				in := domain.ResetPasswordInput{
					Token:    a.req.Token,
					Password: a.req.Password,
				}
				out := &domain.ResetPasswordOutput{Message: "message"}
				resetPasswordMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
					telemetry:       tel,
					resetPasswordUC: resetPasswordMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := tt.mockFn(tt.args)
			got, err := g.ResetPassword(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
