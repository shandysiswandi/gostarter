package inbound

import (
	"context"
	"reflect"
	"testing"

	pb "github.com/shandysiswandi/gostarter/api/gen-proto/auth"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/stretchr/testify/assert"
)

func TestNewGrpcEndpoint(t *testing.T) {
	type args struct {
		telemetry        *telemetry.Telemetry
		loginUC          domain.Login
		registerUC       domain.Register
		refreshTokenUC   domain.RefreshToken
		forgotPasswordUC domain.ForgotPassword
		resetPasswordUC  domain.ResetPassword
	}
	tests := []struct {
		name string
		args args
		want *GrpcEndpoint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewGrpcEndpoint(
				tt.args.telemetry,
				tt.args.loginUC,
				tt.args.registerUC,
				tt.args.refreshTokenUC,
				tt.args.forgotPasswordUC,
				tt.args.resetPasswordUC,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

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
		mockFn  func(a args) *GrpcEndpoint
	}{
		// TODO: Add test cases.
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
		g       *GrpcEndpoint
		args    args
		want    *pb.RegisterResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.Register(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcEndpoint.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcEndpoint.Register() = %v, want %v", got, tt.want)
			}
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
		g       *GrpcEndpoint
		args    args
		want    *pb.RefreshTokenResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.RefreshToken(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcEndpoint.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcEndpoint.RefreshToken() = %v, want %v", got, tt.want)
			}
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
		g       *GrpcEndpoint
		args    args
		want    *pb.ForgotPasswordResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.ForgotPassword(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcEndpoint.ForgotPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcEndpoint.ForgotPassword() = %v, want %v", got, tt.want)
			}
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
		g       *GrpcEndpoint
		args    args
		want    *pb.ResetPasswordResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.ResetPassword(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcEndpoint.ResetPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcEndpoint.ResetPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
