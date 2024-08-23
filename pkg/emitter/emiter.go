// Package emitter provides a simple event emitter implementation
// that allows for registering listeners, emitting events, and
// managing listeners in a thread-safe manner.
package emitter

import (
	"io"
)

// Emitter represents an event emitter that allows firing events
// and listening to them. It provides methods for emitting events
// to registered listeners, adding listeners to specific topics,
// and removing listeners from topics.
type Emitter interface {
	io.Closer

	// Emit sends an event to all listeners registered under the
	// specified topic. This method is fire-and-forget; it does not
	// wait for listeners to consume the event. If the provided topic
	// or arguments are invalid, it will return an error.
	//
	// Parameters:
	//   - topic: A string representing the topic to which the event is sent.
	//   - args: Optional arguments associated with the event.
	//
	// Returns:
	//   - error: An error if the event could not be emitted due to invalid
	//     topic or arguments.
	Emit(topic string, args ...any) error

	// AddListener registers a listener to receive events for the
	// specified topic. It returns a channel from which the events
	// can be received. Listeners will receive events that are emitted
	// with the corresponding topic.
	//
	// Parameters:
	//   - topic: A string representing the topic to which the listener
	//     will be subscribed.
	//
	// Returns:
	//   - <-chan Event: A channel through which events for the specified
	//     topic can be received.
	AddListener(topic string) <-chan Event

	// RemoveListener unregisters listeners from the specified topic.
	// If no specific listeners are provided, all listeners under the
	// topic are unregistered. If specific listeners are provided,
	// only those listeners are unregistered and their channels are closed.
	//
	// Parameters:
	//   - topic: A string representing the topic from which listeners
	//     should be removed.
	//   - listeners: Optional channels to be removed. If no channels are
	//     provided, all listeners for the topic will be removed.
	RemoveListener(topic string, listeners ...<-chan Event)
}
