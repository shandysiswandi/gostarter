// Package uid provides interfaces and implementations for generating unique identifiers (UIDs).
package uid

import (
	nanoid "github.com/matoous/go-nanoid/v2"
)

// NanoIDString is an implementation of UIDString that uses go-nanoid for generating unique string-based UIDs.
//
// By default, NanoIDString generates UIDs of length 21. You can configure the length of the generated UIDs
// by using the NewNanoIDString function with a specific length.
type NanoIDString struct {
	length int
}

// NewNanoIDString creates a new NanoIDString instance with the specified length.
//
// The length parameter must be between 1 and 255. If the length is outside this range, it returns an error.
func NewNanoIDString(length int) (*NanoIDString, error) {
	if length < 1 || length > 255 {
		return nil, ErrInvalidLength
	}

	return &NanoIDString{length: length}, nil
}

// Generate generates a unique identifier as a string using NanoID with the specified length.
//
// It uses the length specified at the time of creation. The length must be between 1 and 255.
func (n *NanoIDString) Generate() string {
	return nanoid.Must(n.length)
}
