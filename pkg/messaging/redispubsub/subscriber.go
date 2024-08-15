// Package redispubsub provides a client for interacting with Redis Pub/Sub.
// It includes functionality for publishing and subscribing to messages with various configuration options.
package redispubsub

import "context"

// SubscriberHandler manages a subscription and its associated resources.
// It provides methods to gracefully close the subscription and release resources.
type SubscriberHandler struct {
	cancelFunc context.CancelFunc // Function to cancel the subscription context.
	client     *Client            // The client associated with this handler.
	name       string             // The name of the subscription.
}

// Close terminates the subscription and releases associated resources.
// It cancels the subscription context by calling the cancelFunc and removes
// the subscription handler from the client.
//
// Returns nil if the handler was successfully closed.
func (sh *SubscriberHandler) Close() error {
	if sh.cancelFunc != nil {
		// Remove the subscription handler from the client.
		sh.client.removeHandler(sh.name)
		// Cancel the subscription context.
		sh.cancelFunc()
		// Set the cancelFunc to nil to prevent double cancellation.
		sh.cancelFunc = nil
	}

	// Return nil as no error handling is implemented in this method.
	return nil
}
