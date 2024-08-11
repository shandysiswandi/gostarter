// Package codec provides a set of interfaces and implementations for encoding
// and decoding data in various formats.
//
// It defines the Codec interface and provides implementations for JSON, Gob,
// MessagePack, and XML formats.
package codec

// Codec is an interface for encoding and decoding data.
//
// The Encode method takes a value of any type and returns a byte slice
// containing the encoded data, or an error if encoding fails.
//
// The Decode method takes a byte slice containing encoded data and a pointer
// to a value where the decoded data will be stored, and returns an error
// if decoding fails.
type Codec interface {
	// Encode encodes a value into a byte slice.
	// v: The value to encode. It can be of any type.
	// Returns: A byte slice containing the encoded data, and an error if encoding fails.
	Encode(v any) ([]byte, error)

	// Decode decodes data from a byte slice into a value.
	// data: The byte slice containing the encoded data.
	// v: A pointer to the value where the decoded data will be stored.
	// Returns: An error if decoding fails.
	Decode(data []byte, v any) error
}
