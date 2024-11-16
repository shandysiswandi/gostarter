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
func (*GobCodec) Encode(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decode decodes Gob data from a byte slice into a value.
func (*GobCodec) Decode(data []byte, v any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	return dec.Decode(v)
}
