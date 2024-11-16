package codec

import (
	"encoding/xml"
)

// XMLCodec is a Codec implementation for XML encoding/decoding.
type XMLCodec struct{}

// NewXMLCodec creates a new instance of XMLCodec.
func NewXMLCodec() *XMLCodec {
	return &XMLCodec{}
}

// Encode encodes a value into XML format.
func (*XMLCodec) Encode(v any) ([]byte, error) {
	return xml.Marshal(v)
}

// Decode decodes XML data from a byte slice into a value.
func (*XMLCodec) Decode(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}
