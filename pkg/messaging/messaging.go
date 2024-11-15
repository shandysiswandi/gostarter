// Package messaging provides an interface for publishing and subscribing
// to messages within a messaging system. It defines a Client interface
// for interacting with the messaging system, as well as types for handling
// message processing and subscription management.
package messaging

import (
	"context"
	"io"
)

// SubscriberHandlerFunc is a function type that processes a message received on a specific topic
// and subscription. It takes a context for managing deadlines and cancellations,
// a topic string, a subscription ID string, and the message content as a byte slice.
// It returns an error if the message could not be processed.
type SubscriberHandlerFunc func(ctx context.Context, topic string, subscriptionID string, msg []byte) error

// SubscriptionHandler is an interface for managing a subscription to a topic.
// It embeds the io.Closer interface, requiring an implementation of the Close method,
// which should clean up any resources associated with the subscription.
type SubscriptionHandler interface {
	io.Closer
}

// Client is an interface that defines methods for interacting with a messaging system.
// It embeds the io.Closer interface, requiring an implementation of the Close method,
// which should be used to clean up resources when the client is no longer needed.
type Client interface {
	io.Closer

	// Publish sends a single message to the specified topic. The context parameter
	// allows the operation to be cancelled or have a timeout applied. It returns
	// an error if the message could not be published.
	Publish(ctx context.Context, topic string, message []byte) error

	// BulkPublish sends multiple messages to the specified topic. The context
	// parameter allows the operation to be cancelled or have a timeout applied.
	// It returns an error if any of the messages could not be published.
	BulkPublish(ctx context.Context, topic string, messages [][]byte) error

	// Subscribe subscribes to a specified topic with a given subscription ID and
	// handler function. The handler is called whenever a message is received on the
	// topic. The method returns a SubscriptionHandler, which can be used to manage
	// the subscription, and an error if the subscription could not be created.
	Subscribe(ctx context.Context, topic, subscriptionID string, handler SubscriberHandlerFunc) (
		SubscriptionHandler, error)
}
