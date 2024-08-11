// Package uid provides interfaces and implementations for generating unique identifiers (UIDs).
package uid

import "errors"

var ErrInvalidLength = errors.New("uid: invalid length, must be between 1 and 255")

// String defines an interface for generating unique string-based UIDs.
type String interface {
	// Generate generates a unique identifier as a string.
	Generate() string
}

// Number defines an interface for generating unique numeric UIDs.
type Number interface {
	// Generate generates a unique identifier as a uint64 number.
	Generate() uint64
}
