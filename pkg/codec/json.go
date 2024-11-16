package codec

import "encoding/json"

// JSONCodec is a Codec implementation for JSON encoding/decoding.
type JSONCodec struct{}

// NewJSONCodec creates a new instance of JSONCodec.
func NewJSONCodec() *JSONCodec {
	return &JSONCodec{}
}

// Encode encodes a value into JSON format.
func (*JSONCodec) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Decode decodes JSON data from a byte slice into a value.
func (*JSONCodec) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
