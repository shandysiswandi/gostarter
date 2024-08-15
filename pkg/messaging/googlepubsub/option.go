// Package googlepubsub provides a client for interacting with Google Cloud Pub/Sub.
// It includes functionality for publishing and subscribing to messages with various configuration options.
package googlepubsub

// ClientOption defines a function type used to configure the Pub/Sub client.
type ClientOption func(*Client)

// WithAutoAck configures the client to automatically acknowledge messages.
//
// isAuto: If true, messages will be automatically acknowledged after being received.
func WithAutoAck(isAuto bool) ClientOption {
	return func(client *Client) {
		client.autoAck = isAuto
	}
}

// WithAutoCreateTopic configures the client to automatically create a topic if it does not exist.
//
// isAuto: If true, the client will create the topic when attempting to publish or subscribe.
func WithAutoCreateTopic(isAuto bool) ClientOption {
	return func(client *Client) {
		client.autoCreateTopic = isAuto
	}
}

// WithAutoCreateSubscriber configures the client to automatically create a subscription
// and associated topic if they do not exist.
//
// isAuto: If true, the client will create both the subscription and the topic as necessary.
func WithAutoCreateSubscriber(isAuto bool) ClientOption {
	return func(client *Client) {
		if isAuto {
			client.autoCreateTopic = true
		}
		client.autoCreateSubscriber = isAuto
	}
}

// WithSyncPublisher configures the client to publish messages synchronously.
//
// isSyncPublisher: If true, the Publish call will wait for the message to be acknowledged before returning.
func WithSyncPublisher(isSyncPublisher bool) ClientOption {
	return func(client *Client) {
		client.syncPublisher = isSyncPublisher
	}
}
