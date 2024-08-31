package codec

import (
	"encoding/xml"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewXMLCodec(t *testing.T) {
	tests := []struct {
		name string
		want *XMLCodec
	}{
		{
			name: "CreateNewXMLCodec",
			want: &XMLCodec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewXMLCodec()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestXMLCodec_Encode(t *testing.T) {
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
				data, _ := xml.Marshal(v)
				return data
			},
			wantErr: false,
		},
		{
			name: "EncodeStructSuccess",
			v: struct {
				XMLName xml.Name `xml:"anonym"`
				Field1  string   `xml:"field1"`
			}{Field1: "value"},
			want: func(v any) []byte {
				data, _ := xml.Marshal(v)
				return data
			},
			wantErr: false,
		},
		{
			name: "EncodeNilValue",
			v:    nil,
			want: func(v any) []byte {
				data, _ := xml.Marshal(v)
				return data
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &XMLCodec{}
			got, err := c.Encode(tt.v)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want(tt.v), got)
			}
		})
	}
}

func TestXMLCodec_Decode(t *testing.T) {
	tests := []struct {
		name    string
		data    func() []byte
		v       any
		want    any
		wantErr bool
	}{
		{
			name: "DecodeStringSuccess",
			data: func() []byte {
				data, _ := xml.Marshal("test string")
				return data
			},
			v:       new(string),
			want:    "test string",
			wantErr: false,
		},
		{
			name: "DecodeStructSuccess",
			data: func() []byte {
				data, _ := xml.Marshal(struct {
					XMLName xml.Name `xml:"item"`
					Field1  string   `xml:"field1"`
				}{Field1: "value"})
				return data
			},
			v: &struct {
				XMLName xml.Name `xml:"item"`
				Field1  string   `xml:"field1"`
			}{},
			want: struct {
				XMLName xml.Name `xml:"item"`
				Field1  string   `xml:"field1"`
			}{Field1: "value", XMLName: xml.Name{Local: "item"}},
			wantErr: false,
		},
		{
			name:    "DecodeInvalidData",
			data:    func() []byte { return []byte("invalid data") },
			v:       new(string),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &XMLCodec{}
			err := c.Decode(tt.data(), tt.v)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				log.Println(tt.want, "--", tt.v)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, func(v any) any {
					switch v := v.(type) {
					case *string:
						return *v
					case *struct {
						XMLName xml.Name `xml:"item"`
						Field1  string   `xml:"field1"`
					}:
						return *v
					default:
						return v
					}
				}(tt.v))
			}
		})
	}
}
