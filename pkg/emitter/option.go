// Package emitter provides a simple event emitter implementation
// that allows for registering listeners, emitting events, and
// managing listeners in a thread-safe manner.
package emitter

import "time"

// TimeProvider is an interface that allows the EventEmitter
// to obtain the current time. This is useful for testing and
// other scenarios where the standard time.Now() function needs
// to be overridden.
type TimeProvider interface {
	Now() time.Time
}

// Option defines a function that applies a configuration option
// to an EventEmitter instance.
type Option func(*EventEmitter)

// WithTimeProvider is an option that allows a custom TimeProvider
// to be used by the EventEmitter. This is useful for injecting
// a mock or custom time provider for testing purposes.
//
// Usage:
//
//	emitter := NewEventEmitter(WithTimeProvider(myTimeProvider))
func WithTimeProvider(tp TimeProvider) Option {
	return func(ee *EventEmitter) {
		ee.timeProvider = tp
	}
}

// WithTopicValidation is an option that allows a custom validation
// function to be used for validating topics before they are used
// in the EventEmitter. If the validation function returns an error,
// the Emit method will fail with that error.
//
// Usage:
//
//	emitter := NewEventEmitter(WithTopicValidation(myValidator))
func WithTopicValidation(tv func(string) error) Option {
	return func(ee *EventEmitter) {
		ee.topicValidation = tv
	}
}

// WithArgumentsValidation is an option that allows a custom validation
// function to be used for validating the arguments passed to events.
// If the validation function returns an error, the Emit method will
// fail with that error.
//
// Usage:
//
//	emitter := NewEventEmitter(WithArgumentsValidation(myArgValidator))
func WithArgumentsValidation(av func([]any) error) Option {
	return func(ee *EventEmitter) {
		ee.argsValidation = av
	}
}

// WithBufferSize is an option that allows a custom buffer size to be
// specified for the channels created by the EventEmitter. This buffer
// size determines how many events can be queued in a listener's channel
// before blocking occurs.
//
// Usage:
//
//	emitter := NewEventEmitter(WithBufferSize(10))
func WithBufferSize(size int) Option {
	return func(ee *EventEmitter) {
		ee.bufferSize = size
	}
}
