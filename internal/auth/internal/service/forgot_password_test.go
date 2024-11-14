package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

func TestNewForgotPassword(t *testing.T) {
	type args struct {
		t       *telemetry.Telemetry
		v       validation.Validator
		idnum   uid.NumberID
		secHash hash.Hash
		s       ForgotPasswordStore
	}
	tests := []struct {
		name string
		args args
		want *ForgotPassword
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewForgotPassword(tt.args.t, tt.args.v, tt.args.idnum, tt.args.secHash, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewForgotPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForgotPassword_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.ForgotPasswordInput
	}
	tests := []struct {
		name    string
		s       *ForgotPassword
		args    args
		want    *domain.ForgotPasswordOutput
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Call(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ForgotPassword.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ForgotPassword.Call() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForgotPassword_doBest(t *testing.T) {
	type args struct {
		ctx  context.Context
		in   domain.ForgotPasswordInput
		user *domain.User
		ps   *domain.PasswordReset
	}
	tests := []struct {
		name    string
		s       *ForgotPassword
		args    args
		want    *domain.ForgotPasswordOutput
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.doBest(tt.args.ctx, tt.args.in, tt.args.user, tt.args.ps)
			if (err != nil) != tt.wantErr {
				t.Errorf("ForgotPassword.doBest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ForgotPassword.doBest() = %v, want %v", got, tt.want)
			}
		})
	}
}
