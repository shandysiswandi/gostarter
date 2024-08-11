// Package clock provides functionality for retrieving the current time.
//
// It defines the Clocker interface and provides a Clock implementation that
// uses the system clock to return the current time.
package clock

import "time"

// Clocker is an interface that defines the ability to get the current time.
//
// Implementations of this interface should provide a method to return the
// current time as a time.Time value.
type Clocker interface {
	// Now returns the current local time.
	//
	// Returns: The current local time as a time.Time value.
	Now() time.Time
}

// Clock is a concrete implementation of the Clocker interface.
//
// It provides the system time using the standard library's time package.
type Clock struct{}

// NewClock returns a new instance of Clock, which implements the Clocker interface.
//
// Returns: A pointer to a new Clock instance.
func NewClock() *Clock {
	return &Clock{}
}

// Now returns the current local time using the system clock.
//
// Returns: The current local time as a time.Time value.
func (*Clock) Now() time.Time {
	return time.Now()
}
