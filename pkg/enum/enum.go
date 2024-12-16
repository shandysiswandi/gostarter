// Package enum provides a generic way to define enums.
package enum

import (
	"database/sql/driver"
	"errors"
)

var ErrInvalidEnum = errors.New("invalid enum value")

const UnknownValue = "UNKNOWN"

// Enumerate defines the interface that must be implemented by any enum type.
type Enumerate interface {
	// Values returns a map where the key is the enum value and the value is its string representation.
	Values() map[Enumerate]string
}

// DefaultValue is an interface that can be implemented by enums that want to provide a custom default value.
type DefaultValue interface {
	// Default returns a custom default string if the enum value is not valid.
	Default() string
}

// Enum is a generic type that holds an enum value of type T.
type Enum[T Enumerate] struct {
	// enum holds the actual enum value of type T.
	enum T
}

// New creates a new Enum instance with the specified enum value.
func New[T Enumerate](e T) Enum[T] {
	return Enum[T]{enum: e}
}

// Enum get actual value of T.
func (e Enum[T]) Enum() T {
	return e.enum
}

// String returns the string representation of the enum value.
// If the value doesn't exist in the enum, it falls back to a custom default value (if any),
// or to "UNKNOWN" if no custom default is provided.
func (e Enum[T]) String() string {
	// Attempt to find the string representation of the enum value in the Values() map.
	if str, ok := e.enum.Values()[e.enum]; ok {
		return str
	}

	// If the enum implements the DefaultValue interface, return its custom default string.
	if value, ok := any(e.enum).(DefaultValue); ok {
		return value.Default()
	}

	// If no match is found and no custom default is implemented, return "UNKNOWN".
	return UnknownValue
}

// Scan scans a database value into an enum.
func (e *Enum[T]) Scan(value any) error {
	var zero T
	switch col := value.(type) {
	case []byte:
		zero = Parse[T](string(col))
	case string:
		zero = Parse[T](col)
	default:
		return ErrInvalidEnum
	}

	e.enum = zero

	return nil
}

// Value returns the string representation of an enum for database storage.
func (e Enum[T]) Value() (driver.Value, error) {
	return e.String(), nil
}

// Parse converts a string to an enum value of type T.
// It looks up the string in the Values() map for the enum and returns the corresponding enum value.
// If the string doesn't match any valid enum value, it returns the zero value of type T.
func Parse[T Enumerate](s string) T {
	var zeroValue T // Zero value for T, returned on failure.

	// Get the map of values from the enum type.
	for key, val := range zeroValue.Values() {
		if val == s {
			// Attempt to assert the key to type T.
			if enumValue, ok := key.(T); ok {
				return enumValue
			}
		}
	}

	// Return the zero value of T if no match is found.
	return zeroValue
}
