// Package emitter provides a simple event emitter implementation
// that allows for registering listeners, emitting events, and
// managing listeners in a thread-safe manner.
package emitter

import (
	"sync"
	"time"
)

// EventEmitter is a thread-safe event emitter that allows for
// registering listeners, emitting events, and managing listeners.
// It supports custom validation for topics and arguments, and
// allows for configurable buffer sizes for listener channels.
type EventEmitter struct {
	mu              sync.RWMutex
	listeners       map[string][]chan Event
	timeProvider    TimeProvider
	topicValidation func(string) error
	argsValidation  func([]any) error
	bufferSize      int
}

// NewEventEmitter creates a new instance of EventEmitter with
// optional configuration options. The returned EventEmitter
// is ready to be used for adding listeners and emitting events.
//
// Options can include:
//   - WithTimeProvider: To provide a custom time provider.
//   - WithTopicValidation: To specify a custom topic validation function.
//   - WithArgumentsValidation: To specify a custom arguments validation function.
//   - WithBufferSize: To set a custom buffer size for listener channels.
//
// Usage example:
//
//	emitter := NewEventEmitter(WithBufferSize(10), WithTimeProvider(myTimeProvider))
func NewEventEmitter(opts ...Option) *EventEmitter {
	ee := &EventEmitter{
		listeners:       make(map[string][]chan Event),
		argsValidation:  func([]any) error { return nil },
		topicValidation: func(string) error { return nil },
	}

	for _, opt := range opts {
		opt(ee)
	}

	return ee
}

// Emit sends an event to all listeners registered under the
// specified topic. The event includes the topic, arguments, and
// a timestamp. This method does not wait for listeners to consume
// the event and is fire-and-forget. It returns an error if the
// topic or arguments are invalid.
//
// Parameters:
//   - topic: A string representing the topic of the event.
//   - args: Optional arguments associated with the event.
//
// Returns:
//   - error: An error if the event could not be emitted due to
//     validation failure.
func (e *EventEmitter) Emit(topic string, args ...any) error {
	if e.topicValidation == nil || e.argsValidation == nil {
		return nil
	}

	if err := e.topicValidation(topic); err != nil {
		return err
	}

	if err := e.argsValidation(args); err != nil {
		return err
	}

	now := time.Now()
	if e.timeProvider != nil {
		now = e.timeProvider.Now()
	}

	event := Event{topic: topic, args: args, timestamp: now}

	e.mu.RLock()
	listeners, ok := e.listeners[topic]
	e.mu.RUnlock()
	if !ok {
		return nil
	}

	// Avoid holding the lock while sending to channels
	for _, ch := range listeners {
		select {
		case ch <- event:
		default:
			// Optionally handle channels that are full (e.g., log, drop, etc.)
			// This could be a place to log a warning if necessary.
		}
	}

	return nil
}

// AddListener registers a new listener for the specified topic.
// It returns a channel that receives events for that topic.
// If no buffer size was set, a default buffer size of 3 is used.
//
// Parameters:
//   - topic: A string representing the topic to which the listener
//     will subscribe.
//
// Returns:
//   - <-chan Event: A channel from which events for the specified
//     topic can be received.
func (e *EventEmitter) AddListener(topic string) <-chan Event {
	if e.listeners == nil {
		e.listeners = make(map[string][]chan Event)
	}

	if e.bufferSize <= 0 {
		e.bufferSize = 3 // Default buffer size if not set
	}

	ch := make(chan Event, e.bufferSize)

	e.mu.Lock()
	e.listeners[topic] = append(e.listeners[topic], ch)
	e.mu.Unlock()

	return ch
}

// RemoveListener unregisters listeners from the specified topic.
// If no specific listeners are provided, all listeners under
// the topic are unregistered and their channels are closed. If
// specific listeners are provided, only those are unregistered.
//
// Parameters:
//   - topic: A string representing the topic from which listeners
//     should be removed.
//   - listeners: Optional channels to be removed. If none are
//     provided, all listeners for the topic will be removed.
func (e *EventEmitter) RemoveListener(topic string, listeners ...<-chan Event) {
	e.mu.Lock()
	defer e.mu.Unlock()

	chans, found := e.listeners[topic]
	if !found {
		return
	}

	if len(listeners) == 0 {
		// No specific listener provided, remove all listeners for this topic
		for _, ch := range chans {
			close(ch)
		}
		delete(e.listeners, topic)

		return
	}

	// Remove only the specified listeners
	remainingChans := chans[:0]
	for _, ch := range chans {
		shouldClose := false
		for _, listener := range listeners {
			if ch == listener {
				close(ch)
				shouldClose = true

				break
			}
		}
		if !shouldClose {
			remainingChans = append(remainingChans, ch)
		}
	}

	if len(remainingChans) == 0 {
		delete(e.listeners, topic)
	} else {
		e.listeners[topic] = remainingChans
	}
}

// Close stops the EventEmitter, closing all listener channels
// and clearing the internal state. This method ensures that all
// resources are properly cleaned up.
//
// Returns:
//   - error: An error if any issues occur during the cleanup process.
func (e *EventEmitter) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, chans := range e.listeners {
		for _, ch := range chans {
			close(ch)
		}
	}
	e.listeners = nil

	return nil
}
