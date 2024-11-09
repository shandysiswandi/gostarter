package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHMACSHA256Hash(t *testing.T) {
	type args struct {
		secret string
	}
	tests := []struct {
		name string
		args args
		want *HMACSHA256Hash
	}{
		{
			name: "Success",
			args: args{},
			want: &HMACSHA256Hash{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewHMACSHA256Hash(tt.args.secret)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHMACSHA256Hash_Hash(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr error
		h       *HMACSHA256Hash
	}{
		{
			name:    "Success",
			args:    args{str: "hash"},
			want:    []byte("c835893d96769822abb85d90f16c4b1e54a88e573d11a660030a282a02daa3b4"),
			wantErr: nil,
			h:       &HMACSHA256Hash{secret: "hash"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.h.Hash(tt.args.str)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, len(tt.want), len(got))
		})
	}
}

func TestHMACSHA256Hash_Verify(t *testing.T) {
	type args struct {
		hashedHex string
		str       string
	}
	tests := []struct {
		name string
		args args
		want bool
		h    *HMACSHA256Hash
	}{
		{
			name: "Success",
			args: args{
				hashedHex: "c835893d96769822abb85d90f16c4b1e54a88e573d11a660030a282a02daa3b4",
				str:       "hash",
			},
			want: true,
			h:    &HMACSHA256Hash{secret: "hash"},
		},
		{
			name: "Error",
			args: args{
				hashedHex: "hash",
				str:       "hash",
			},
			want: false,
			h:    &HMACSHA256Hash{secret: "hash"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.h.Verify(tt.args.hashedHex, tt.args.str)
			assert.Equal(t, tt.want, got)
		})
	}
}
