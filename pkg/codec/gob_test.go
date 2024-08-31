package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGobCodec(t *testing.T) {
	tests := []struct {
		name string
		want *GobCodec
	}{
		{
			name: "CreateNewGobCodec",
			want: &GobCodec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewGobCodec()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGobCodec_Encode(t *testing.T) {
	tests := []struct {
		name    string
		v       any
		wantErr bool
	}{
		{
			name:    "EncodeStringSuccess",
			v:       "test string",
			wantErr: false,
		},
		{
			name:    "EncodeStructSuccess",
			v:       struct{ Field1 string }{Field1: "value"},
			wantErr: false,
		},
		{
			name:    "EncodeMapSuccess",
			v:       map[string]int{"one": 1, "two": 2},
			wantErr: false,
		},
		{
			name:    "EncodeNilValue",
			v:       nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &GobCodec{}
			_, err := c.Encode(tt.v)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGobCodec_Decode(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		v       any
		want    any
		wantErr bool
	}{
		{
			name: "DecodeStringSuccess",
			data: func() []byte {
				c := &GobCodec{}
				data, _ := c.Encode("test string")
				return data
			}(),
			v:       new(string),
			want:    "test string",
			wantErr: false,
		},
		{
			name: "DecodeStructSuccess",
			data: func() []byte {
				c := &GobCodec{}
				data, _ := c.Encode(struct{ Field1 string }{Field1: "value"})
				return data
			}(),
			v:       &struct{ Field1 string }{},
			want:    struct{ Field1 string }{Field1: "value"},
			wantErr: false,
		},
		{
			name: "DecodeMapSuccess",
			data: func() []byte {
				c := &GobCodec{}
				data, _ := c.Encode(map[string]int{"one": 1, "two": 2})
				return data
			}(),
			v:       &map[string]int{},
			want:    map[string]int{"one": 1, "two": 2},
			wantErr: false,
		},
		{
			name:    "DecodeInvalidData",
			data:    []byte("invalid data"),
			v:       new(string),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &GobCodec{}
			err := c.Decode(tt.data, tt.v)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, func(v any) any {
					switch v := v.(type) {
					case *string:
						return *v
					case *struct{ Field1 string }:
						return *v
					case *map[string]int:
						return *v
					default:
						return v
					}
				}(tt.v))
			}
		})
	}
}
