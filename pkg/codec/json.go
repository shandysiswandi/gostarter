// Package codec provides a JSONCodec implementation of the Codec interface.
//
// JSONCodec is used for encoding and decoding data in JSON format.
package codec

import "encoding/json"

// JSONCodec is a Codec implementation for JSON encoding/decoding.
type JSONCodec struct{}

// NewJSONCodec creates a new instance of JSONCodec.
func NewJSONCodec() *JSONCodec {
	return &JSONCodec{}
}

// Encode encodes a value into JSON format.
//
// v: The value to encode. It can be of any type that is supported by JSON encoding.
// Returns: A byte slice containing the JSON-encoded data, and an error if encoding fails.
func (*JSONCodec) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Decode decodes JSON data from a byte slice into a value.
//
// data: The byte slice containing JSON-encoded data.
// v: A pointer to the value where the decoded data will be stored.
// Returns: An error if decoding fails.
func (*JSONCodec) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
