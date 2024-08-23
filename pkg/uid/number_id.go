// Package uid provides interfaces and implementations for generating unique identifiers (UIDs).
package uid

// NumberID defines an interface for generating unique numeric UIDs.
type NumberID interface {
	// Generate generates a unique identifier as a uint64 number.
	Generate() uint64
}
