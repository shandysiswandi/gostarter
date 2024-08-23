// Package emitter provides a simple event emitter implementation
// that allows for registering listeners, emitting events, and
// managing listeners in a thread-safe manner.
package emitter

import "time"

// TimeProvider is an interface that allows the EventEmitter
// to obtain the current time. This is useful for testing and
// other scenarios where the standard time.Now() function needs
// to be overridden.
//
// Implementing this interface allows you to control the time source
// in a predictable manner, which is particularly valuable in testing
// environments.
//
// Example:
//
//	type MockTimeProvider struct {}
//
//	func (m MockTimeProvider) Now() time.Time {
//	    return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
//	}
//
//	emitter := NewEventEmitter(WithTimeProvider(MockTimeProvider{}))
type TimeProvider interface {
	Now() time.Time
}

// Options holds the configuration options for the EventEmitter.
//
// You can customize the behavior of the EventEmitter by providing
// these options, such as a custom time provider, topic validation,
// arguments validation, and buffer size.
type Options struct {
	TimeProvider        TimeProvider
	TopicValidation     func(string) error
	ArgumentsValidation func([]any) error
	BufferSize          int
}

// Option defines a function that applies a configuration option
// to an EventEmitter instance.
//
// Options allow you to customize the EventEmitter's behavior when
// it is created. Each option is a function that modifies the Options
// struct.
type Option func(*Options)

// WithTimeProvider is an option that allows a custom TimeProvider
// to be used by the EventEmitter. This is useful for injecting
// a mock or custom time provider for testing purposes.
//
// Usage:
//
//	emitter := NewEventEmitter(WithTimeProvider(myTimeProvider))
func WithTimeProvider(tp TimeProvider) Option {
	return func(o *Options) {
		o.TimeProvider = tp
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
	return func(o *Options) {
		o.TopicValidation = tv
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
	return func(o *Options) {
		o.ArgumentsValidation = av
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
	return func(o *Options) {
		o.BufferSize = size
	}
}
