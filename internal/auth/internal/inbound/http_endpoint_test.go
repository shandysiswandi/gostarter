package inbound

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
)

func TestRegisterRESTEndpoint(t *testing.T) {
	type args struct {
		router *httprouter.Router
		h      *Endpoint
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterAuthServiceServer(tt.args.router, tt.args.h)
		})
	}
}

func TestNewEndpoint(t *testing.T) {
	type args struct {
		loginUC          domain.Login
		registerUC       domain.Register
		refreshTokenUC   domain.RefreshToken
		forgotPasswordUC domain.ForgotPassword
		resetPasswordUC  domain.ResetPassword
	}
	tests := []struct {
		name string
		args args
		want *Endpoint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHTTPEndpoint(tt.args.loginUC, tt.args.registerUC, tt.args.refreshTokenUC, tt.args.forgotPasswordUC, tt.args.resetPasswordUC); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoint_Login(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.Login(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Endpoint.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoint_Register(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.Register(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Endpoint.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoint_RefreshToken(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.RefreshToken(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Endpoint.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint.RefreshToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoint_ForgotPassword(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.ForgotPassword(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Endpoint.ForgotPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint.ForgotPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoint_ResetPassword(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.ResetPassword(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Endpoint.ResetPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint.ResetPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
