package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewJWTSymetric(t *testing.T) {
	type args struct {
		secret []byte
	}
	tests := []struct {
		name string
		args args
		want *Symetric
	}{
		{
			name: "Test with non-empty secret",
			args: args{secret: []byte("test")},
			want: &Symetric{secret: []byte("test")},
		},
		{
			name: "Test with empty secret",
			args: args{secret: []byte("")},
			want: &Symetric{secret: []byte("")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewJWTSymetric(tt.args.secret)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJWTSymetric_Generate(t *testing.T) {
	tests := []struct {
		name    string
		arg     *Claim
		want    string
		wantErr error
		mockFn  func() *Symetric
	}{
		{
			name: "Success",
			arg: &Claim{
				AuthID: 101,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "test",
					Subject:   "test",
					Audience:  []string{"test"},
					ExpiresAt: jwt.NewNumericDate(time.Date(2034, time.December, 1, 0, 0, 0, 0, time.Local)),
					NotBefore: jwt.NewNumericDate(time.Date(2024, time.December, 1, 0, 0, 0, 0, time.Local)),
					IssuedAt:  jwt.NewNumericDate(time.Date(2024, time.December, 1, 0, 0, 0, 0, time.Local)),
				},
			},
			wantErr: nil,
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX2lkIjoiMTAxIiwiaXNzIjoidGVzdCIsInN1YiI6InRlc3QiLCJhdWQiOlsidGVzdCJdLCJleHAiOjIwNDg1MTg4MDAsIm5iZiI6MTczMjk4NjAwMCwiaWF0IjoxNzMyOTg2MDAwfQ.X6oa_41wmqWyjoT8ckg7Psj-jzJXEutW3luh--leVoc",
			mockFn: func() *Symetric {
				return NewJWTSymetric([]byte("test"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			js := tt.mockFn()
			got, err := js.Generate(tt.arg)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJWTSymetric_Verify(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		want    *Claim
		wantErr bool
		mockFn  func(t string) *Symetric
	}{
		{
			name:  "Success",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX2lkIjoiMTAxIiwiaXNzIjoidGVzdCIsInN1YiI6InRlc3QiLCJhdWQiOlsidGVzdCJdLCJleHAiOjIwNDg1MTg4MDAsIm5iZiI6MTczMjk4NjAwMCwiaWF0IjoxNzMyOTg2MDAwfQ.X6oa_41wmqWyjoT8ckg7Psj-jzJXEutW3luh--leVoc",
			want: &Claim{
				AuthID: 101,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "test",
					Subject:   "test",
					Audience:  []string{"test"},
					ExpiresAt: jwt.NewNumericDate(time.Date(2034, time.December, 1, 0, 0, 0, 0, time.Local)),
					NotBefore: jwt.NewNumericDate(time.Date(2024, time.December, 1, 0, 0, 0, 0, time.Local)),
					IssuedAt:  jwt.NewNumericDate(time.Date(2024, time.December, 1, 0, 0, 0, 0, time.Local)),
				},
			},
			wantErr: false,
			mockFn: func(t string) *Symetric {
				return NewJWTSymetric([]byte("test"))
			},
		},
		{
			name:    "ErrorExpired",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX2lkIjoiMTAxIiwiaXNzIjoidGVzdCIsInN1YiI6InRlc3QiLCJhdWQiOlsidGVzdCJdLCJleHAiOjExMDE4MzQwMDAsIm5iZiI6MTczMjk4NjAwMCwiaWF0IjoxNzMyOTg2MDAwfQ.UmgDfeLb-d_L7ZKq-33inhqoLR2jfXnmh3_jPaf9LoQ",
			want:    nil,
			wantErr: true,
			mockFn: func(t string) *Symetric {
				return NewJWTSymetric([]byte("test"))
			},
		},
		{
			name:    "ErrorVerify",
			token:   "",
			want:    nil,
			wantErr: true,
			mockFn: func(t string) *Symetric {
				return NewJWTSymetric([]byte("test"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			js := tt.mockFn(tt.token)
			got, err := js.Verify(tt.token)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
