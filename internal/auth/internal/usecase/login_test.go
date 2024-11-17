package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

func TestNewLogin(t *testing.T) {
	type args struct {
		t       *telemetry.Telemetry
		v       validation.Validator
		idnum   uid.NumberID
		hash    hash.Hash
		secHash hash.Hash
		j       jwt.JWT
		s       LoginStore
	}
	tests := []struct {
		name string
		args args
		want *Login
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogin(tt.args.t, tt.args.v, tt.args.idnum, tt.args.hash, tt.args.secHash, tt.args.j, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogin_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.LoginInput
	}
	tests := []struct {
		name    string
		s       *Login
		args    args
		want    *domain.LoginOutput
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Call(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login.Call() = %v, want %v", got, tt.want)
			}
		})
	}
}
