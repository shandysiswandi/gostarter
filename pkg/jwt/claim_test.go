package jwt

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestClaim_Now(t *testing.T) {
	tests := []struct {
		name string
		c    *Claim
		want time.Time
	}{
		{
			name: "Success",
			c:    &Claim{},
			want: time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.c.Now()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewClaim(t *testing.T) {
	type args struct {
		email string
		exp   time.Duration
		aud   []string
	}
	tests := []struct {
		name string
		args args
		want *Claim
	}{
		{
			name: "Success",
			args: args{
				email: "email@email.com",
				exp:   1,
				aud:   []string{"aud"},
			},
			want: &Claim{
				Email: "email@email.com",
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "gostarter",
					Subject:   "email@email.com",
					Audience:  []string{"aud"},
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(1)),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewClaim(tt.args.email, tt.args.exp, tt.args.aud)
			assert.Equal(t, tt.want.Email, got.Email)
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
			args: args{token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImVtYWlsQGVtYWlsLmNvbSIsImlhdCI6MTUxNjIzOTAyMn0.UOQFRvx2JwT1PcDKqbfj9f_WN66Gs_giUMGv3bgVcE8"},
			want: &Claim{
				now:   time.Time{},
				Email: "email@email.com",
				RegisteredClaims: jwt.RegisteredClaims{
					IssuedAt: jwt.NewNumericDate(time.Date(2018, time.January, 18, 8, 30, 22, 0, time.Local)),
					ID:       "",
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
