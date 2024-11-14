package service

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

func TestNewRefreshToken(t *testing.T) {
	type args struct {
		t       *telemetry.Telemetry
		v       validation.Validator
		idnum   uid.NumberID
		secHash hash.Hash
		j       jwt.JWT
		s       RefreshTokenStore
	}
	tests := []struct {
		name string
		args args
		want *RefreshToken
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRefreshToken(tt.args.t, tt.args.v, tt.args.idnum, tt.args.secHash, tt.args.j, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRefreshToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRefreshToken_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.RefreshTokenInput
	}
	tests := []struct {
		name    string
		s       *RefreshToken
		args    args
		want    *domain.RefreshTokenOutput
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Call(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RefreshToken.Call() = %v, want %v", got, tt.want)
			}
		})
	}
}
