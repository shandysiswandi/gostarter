// Package codec provides an XMLCodec implementation of the Codec interface.
//
// XMLCodec is used for encoding and decoding data in XML format.
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
//
// v: The value to encode. It can be of any type that is supported by XML encoding.
// Returns: A byte slice containing the XML-encoded data, and an error if encoding fails.
func (*XMLCodec) Encode(v any) ([]byte, error) {
	return xml.Marshal(v)
}

// Decode decodes XML data from a byte slice into a value.
//
// data: The byte slice containing XML-encoded data.
// v: A pointer to the value where the decoded data will be stored.
// Returns: An error if decoding fails.
func (*XMLCodec) Decode(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}
