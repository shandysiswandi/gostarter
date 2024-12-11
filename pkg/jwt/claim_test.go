package jwt

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewClaim(t *testing.T) {
	type args struct {
		authID uint64
		email  string
		exp    time.Time
		aud    []string
	}
	tests := []struct {
		name string
		args args
		want *Claim
	}{
		{
			name: "Success",
			args: args{
				authID: 101,
				email:  "email@email.com",
				exp:    time.Time{},
				aud:    []string{"aud"},
			},
			want: &Claim{
				AuthID: 101,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "GO_STARTER",
					Subject:   "email@email.com",
					Audience:  []string{"aud"},
					ExpiresAt: jwt.NewNumericDate(time.Time{}),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewClaim(tt.args.authID, tt.args.email, tt.args.exp, tt.args.aud)
			assert.Equal(t, tt.want.AuthID, got.AuthID)
			assert.Equal(t, tt.want.RegisteredClaims, got.RegisteredClaims)
		})
	}
}

func TestExtractClaimFromToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want *Claim
	}{
		{
			name: "Success",
			args: args{token: "a.eyJzdWIiOiJ0ZXN0IiwiYXV0aF9pZCI6IjEwMSJ9.a"},
			want: &Claim{
				AuthID: 101,
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "test",
				},
			},
		},
		{
			name: "ErrorFormat",
			args: args{token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"},
			want: nil,
		},
		{
			name: "ErrorPartDecode",
			args: args{token: "1.1.1"},
			want: nil,
		},
		{
			name: "ErrorPartUnmarshal",
			args: args{token: "1.eyJlbWFpbCI6dHJ1ZSwiaWF0Ijp0cnVlfQ.1"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ExtractClaimFromToken(tt.args.token)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSetClaim(t *testing.T) {
	type args struct {
		ctx context.Context
		clm *Claim
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(a args) context.Context
	}{
		{
			name: "Success",
			args: args{ctx: context.Background(), clm: &Claim{}},
			mockFn: func(a args) context.Context {
				return context.WithValue(a.ctx, contextKey{}, a.clm)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := tt.mockFn(tt.args)
			got := SetClaim(tt.args.ctx, tt.args.clm)
			assert.Equal(t, ctx, got)
		})
	}
}

func TestGetClaim(t *testing.T) {
	tests := []struct {
		name string
		ctx  func() context.Context
		want *Claim
	}{
		{
			name: "NoClaim",
			ctx: func() context.Context {
				return context.Background()
			},
			want: nil,
		},
		{
			name: "Success",
			ctx: func() context.Context {
				return SetClaim(context.Background(), &Claim{})
			},
			want: &Claim{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := GetClaim(tt.ctx())
			assert.Equal(t, tt.want, got)
		})
	}
}
