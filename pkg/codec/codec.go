// Package codec provides a set of interfaces and implementations for encoding
// and decoding data in various formats. It defines the Codec interface and provides
// implementations for JSON, Gob, MessagePack, and XML formats.
package codec

// Codec is an interface for encoding and decoding data.
type Codec interface {
	// Encode encodes a value into a byte slice.
	Encode(v any) ([]byte, error)

	// Decode decodes data from a byte slice into a value.
	Decode(data []byte, v any) error
}
