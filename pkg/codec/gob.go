// Package codec provides a GobCodec implementation of the Codec interface.
//
// GobCodec is used for encoding and decoding data in Go's Gob format.
package codec

import (
	"bytes"
	"encoding/gob"
)

// GobCodec is a Codec implementation for gob encoding/decoding.
type GobCodec struct{}

// NewGobCodec creates a new instance of GobCodec.
func NewGobCodec() *GobCodec {
	return &GobCodec{}
}

// Encode encodes a value into Gob format.
//
// v: The value to encode. It can be of any type that is supported by Gob encoding.
// Returns: A byte slice containing the Gob-encoded data, and an error if encoding fails.
func (*GobCodec) Encode(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decode decodes Gob data from a byte slice into a value.
//
// data: The byte slice containing Gob-encoded data.
// v: A pointer to the value where the decoded data will be stored.
// Returns: An error if decoding fails.
func (*GobCodec) Decode(data []byte, v any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	return dec.Decode(v)
}
