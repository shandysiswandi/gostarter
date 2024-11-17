// Package messaging provides an interface for publishing and subscribing
// to messages within a messaging system. It defines a Client interface
// for interacting with the messaging system, as well as types for handling
// message processing and subscription management.
package messaging

import (
	"context"
	"errors"
	"io"
)

var (
	ErrDataNil    = errors.New("data nil")
	ErrMessageNil = errors.New("data message nil")
)

// SubscriberFunc defines a callback function to process messages received from a subscription.
type SubscriberFunc func(ctx context.Context, data *Data) error

// Client defines methods for interacting with a messaging system.
type Client interface {
	io.Closer

	// Publish sends a single message to the specified topic. The context parameter
	// allows cancellation or timeouts. Returns an error if the operation fails.
	Publish(ctx context.Context, topic string, data *Data) error

	// BulkPublish sends multiple messages to the specified topic in a single operation.
	// If any message fails, an aggregated error is returned. Atomicity is implementation-defined.
	BulkPublish(ctx context.Context, topic string, data []*Data) error

	// Subscribe sets up a subscription to a topic, invoking the provided SubscriberFunc
	// for each received message. Returns an error if subscription setup fails.
	Subscribe(ctx context.Context, topic, subscriptionID string, handler SubscriberFunc) error
}

// Data represents a message with its payload and attributes.
type Data struct {
	Msg        []byte
	Attributes map[string]string
}

func (d *Data) Validate() error {
	if d == nil {
		return ErrDataNil
	}

	if d.Msg == nil {
		return ErrMessageNil
	}

	return nil
}
