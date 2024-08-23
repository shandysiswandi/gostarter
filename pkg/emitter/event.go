// Package emitter provides a simple event emitter implementation
// that allows for registering listeners, emitting events, and
// managing listeners in a thread-safe manner.
package emitter

import "time"

// Event represents a topic-based event with associated arguments.
// It includes the topic to which the event belongs, the arguments
// associated with the event, and the timestamp when the event was created.
type Event struct {
	topic     string
	args      []any
	timestamp time.Time
}

// GetTopic returns the topic of the event. The topic is a string
// that identifies the category or type of the event.
//
// Returns:
//   - string: The topic of the event.
func (e *Event) GetTopic() string {
	return e.topic
}

// GetArgs returns the arguments associated with the event. These
// arguments are the data passed along with the event and can be
// of any type.
//
// Returns:
//   - []any: A slice containing the arguments associated with the event.
func (e *Event) GetArgs() []any {
	return e.args
}

// GetTimestamp returns the time when the event was created. This
// timestamp indicates when the event was emitted and can be used
// to track or sort events based on their occurrence.
//
// Returns:
//   - time.Time: The timestamp of the event creation.
func (e *Event) GetTimestamp() time.Time {
	return e.timestamp
}
