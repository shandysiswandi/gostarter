// Package uid provides interfaces and implementations for generating unique identifiers (UIDs).
package uid

// StringID defines an interface for generating unique string-based UIDs.
type StringID interface {
	// Generate generates a unique identifier as a string.
	Generate() string
}
