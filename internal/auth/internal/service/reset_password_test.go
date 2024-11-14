package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

func TestNewResetPassword(t *testing.T) {
	type args struct {
		t *telemetry.Telemetry
		v validation.Validator
		h hash.Hash
		s ResetPasswordStore
	}
	tests := []struct {
		name string
		args args
		want *ResetPassword
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResetPassword(tt.args.t, tt.args.v, tt.args.h, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResetPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResetPassword_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.ResetPasswordInput
	}
	tests := []struct {
		name    string
		s       *ResetPassword
		args    args
		want    *domain.ResetPasswordOutput
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Call(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResetPassword.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResetPassword.Call() = %v, want %v", got, tt.want)
			}
		})
	}
}
