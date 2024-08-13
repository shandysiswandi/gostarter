package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func TestNewMsgpackCodec(t *testing.T) {
	tests := []struct {
		name string
		want *MsgpackCodec
	}{
		{
			name: "CreateNewMsgpackCodec",
			want: &MsgpackCodec{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewMsgpackCodec()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMsgpackCodec_Encode(t *testing.T) {
	tests := []struct {
		name    string
		v       any
		want    func(v any) []byte
		wantErr bool
	}{
		{
			name: "EncodeStringSuccess",
			v:    "test string",
			want: func(v any) []byte {
				data, _ := msgpack.Marshal(v)
				return data
			},
			wantErr: false,
		},
		{
			name: "EncodeStructSuccess",
			v: struct {
				Field1 string `msgpack:"field1"`
			}{Field1: "value"},
			want: func(v any) []byte {
				data, _ := msgpack.Marshal(v)
				return data
			},
			wantErr: false,
		},
		{
			name: "EncodeMapSuccess",
			v:    map[string]int{"one": 1, "two": 2},
			want: func(v any) []byte {
				data, _ := msgpack.Marshal(v)
				return data
			},
			wantErr: false,
		},
		{
			name: "EncodeNilValue",
			v:    nil,
			want: func(v any) []byte {
				data, _ := msgpack.Marshal(v)
				return data
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &MsgpackCodec{}
			got, err := c.Encode(tt.v)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				var decodedExpected, decodedGot interface{}
				_ = msgpack.Unmarshal(tt.want(tt.v), &decodedExpected)
				_ = msgpack.Unmarshal(got, &decodedGot)
				assert.Equal(t, decodedExpected, decodedGot)
			}
		})
	}
}

func TestMsgpackCodec_Decode(t *testing.T) {
	tests := []struct {
		name    string
		v       any
		data    func() []byte
		want    any
		wantErr bool
	}{
		{
			name: "DecodeStringSuccess",
			data: func() []byte {
				data, _ := msgpack.Marshal("test string")
				return data
			},
			v:       new(string),
			want:    "test string",
			wantErr: false,
		},
		{
			name: "DecodeStructSuccess",
			data: func() []byte {
				data, _ := msgpack.Marshal(struct {
					Field1 string `msgpack:"field1"`
				}{Field1: "value"})
				return data
			},
			v: &struct {
				Field1 string `msgpack:"field1"`
			}{},
			want: struct {
				Field1 string `msgpack:"field1"`
			}{Field1: "value"},
			wantErr: false,
		},
		{
			name: "DecodeMapSuccess",
			data: func() []byte {
				data, _ := msgpack.Marshal(map[string]int{"one": 1, "two": 2})
				return data
			},
			v:       &map[string]int{},
			want:    map[string]int{"one": 1, "two": 2},
			wantErr: false,
		},
		{
			name: "DecodeInvalidData",
			data: func() []byte {
				return []byte("invalid data")
			},
			v:       new(string),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &MsgpackCodec{}
			err := c.Decode(tt.data(), tt.v)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, func(v any) any {
					switch v := v.(type) {
					case *string:
						return *v
					case *struct {
						Field1 string `msgpack:"field1"`
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
