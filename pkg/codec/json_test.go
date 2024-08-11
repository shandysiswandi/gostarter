package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJSONCodec(t *testing.T) {
	tests := []struct {
		name string
		want *JSONCodec
	}{
		{
			name: "CreateNewJSONCodec",
			want: &JSONCodec{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewJSONCodec()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJSONCodec_Encode(t *testing.T) {
	tests := []struct {
		name    string
		v       any
		want    string
		wantErr bool
	}{
		{
			name:    "EncodeStringSuccess",
			v:       "test string",
			want:    `"test string"`,
			wantErr: false,
		},
		{
			name: "EncodeStructSuccess",
			v: struct {
				Field1 string `json:"field1"`
			}{Field1: "value"},
			want:    `{"field1":"value"}`,
			wantErr: false,
		},
		{
			name:    "EncodeMapSuccess",
			v:       map[string]int{"one": 1, "two": 2},
			want:    `{"one":1,"two":2}`,
			wantErr: false,
		},
		{
			name:    "EncodeNilValue",
			v:       nil,
			want:    `null`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &JSONCodec{}
			got, err := c.Encode(tt.v)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, tt.want, string(got))
			}
		})
	}
}

func TestJSONCodec_Decode(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		v       any
		want    any
		wantErr bool
	}{
		{
			name:    "DecodeStringSuccess",
			data:    []byte(`"test string"`),
			v:       new(string),
			want:    "test string",
			wantErr: false,
		},
		{
			name: "DecodeStructSuccess",
			data: []byte(`{"field1":"value"}`),
			v: &struct {
				Field1 string `json:"field1"`
			}{},
			want: struct {
				Field1 string `json:"field1"`
			}{Field1: "value"},
			wantErr: false,
		},
		{
			name:    "DecodeMapSuccess",
			data:    []byte(`{"one":1,"two":2}`),
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &JSONCodec{}
			err := c.Decode(tt.data, tt.v)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, func(v any) any {
					switch v := v.(type) {
					case *string:
						return *v
					case *struct {
						Field1 string `json:"field1"`
					}:
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
