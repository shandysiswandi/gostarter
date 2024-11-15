// Package uid provides interfaces and implementations for generating unique identifiers (UIDs).
package uid

import "github.com/google/uuid"

// UUIDString is an implementation of UIDString that uses UUID for generating unique string-based UIDs.
type UUIDString struct{}

// NewUUIDString creates and returns a new instance of UUIDString.
func NewUUIDString() *UUIDString {
	return &UUIDString{}
}

// Generate generates a unique identifier as a string using google uuid v4.
func (u *UUIDString) Generate() string {
	return uuid.NewString()
}
