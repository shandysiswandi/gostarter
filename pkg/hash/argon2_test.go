package hash

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArgon2Hash(t *testing.T) {
	type args struct {
		time    uint32
		memory  uint32
		threads uint8
		keyLen  uint32
	}
	tests := []struct {
		name string
		args args
		want *Argon2Hash
	}{
		{
			name: "Success",
			args: args{},
			want: &Argon2Hash{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewArgon2Hash(tt.args.time, tt.args.memory, tt.args.threads, tt.args.keyLen)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArgon2Hash_Hash(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr error
		h       *Argon2Hash
	}{
		{
			name:    "Success",
			args:    args{str: "hash"},
			want:    []byte("34c7f834f2c91274a0bb457a4f21877e:df8dff20d173d80d"),
			wantErr: nil,
			h: &Argon2Hash{
				time:    1,
				threads: 1,
				memory:  1,
				keyLen:  8,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.Hash(tt.args.str)
			assert.Equal(t, tt.wantErr, err)
			log.Println(string(got))
			assert.Equal(t, len(tt.want), len(got))
		})
	}
}

func TestArgon2Hash_Verify(t *testing.T) {
	type args struct {
		hashed string
		str    string
	}
	tests := []struct {
		name string
		args args
		want bool
		h    *Argon2Hash
	}{
		{
			name: "Success",
			args: args{hashed: "34c7f834f2c91274a0bb457a4f21877e:df8dff20d173d80d", str: "hash"},
			want: true,
			h: &Argon2Hash{
				time:    1,
				threads: 1,
				memory:  1,
				keyLen:  8,
			},
		},
		{
			name: "ErrorFormatColon",
			args: args{hashed: "34c7f834f2c91274a0bb457a4f21877e", str: "hash"},
			want: false,
			h: &Argon2Hash{
				time:    1,
				threads: 1,
				memory:  1,
				keyLen:  8,
			},
		},
		{
			name: "ErrorFormatSalt",
			args: args{hashed: "1:1", str: "hash"},
			want: false,
			h: &Argon2Hash{
				time:    1,
				threads: 1,
				memory:  1,
				keyLen:  8,
			},
		},
		{
			name: "ErrorFormatHashHex",
			args: args{hashed: "34c7f834f2c91274a0bb457a4f21877e:1", str: "hash"},
			want: false,
			h: &Argon2Hash{
				time:    1,
				threads: 1,
				memory:  1,
				keyLen:  8,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.h.Verify(tt.args.hashed, tt.args.str)
			assert.Equal(t, tt.want, got)
		})
	}
}
