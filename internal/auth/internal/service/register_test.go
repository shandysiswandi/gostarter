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

func TestNewRegister(t *testing.T) {
	type args struct {
		t     *telemetry.Telemetry
		v     validation.Validator
		idnum uid.NumberID
		hash  hash.Hash
		s     RegisterStore
	}
	tests := []struct {
		name string
		args args
		want *Register
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRegister(tt.args.t, tt.args.v, tt.args.idnum, tt.args.hash, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRegister() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegister_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.RegisterInput
	}
	tests := []struct {
		name    string
		s       *Register
		args    args
		want    *domain.RegisterOutput
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Call(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register.Call() = %v, want %v", got, tt.want)
			}
		})
	}
}
