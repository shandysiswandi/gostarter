package googlepubsub

import (
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"google.golang.org/api/option"
)

// Option defines a configuration function for the Client.
type Option func(*Client)

// WithAutoAck configures the client to automatically acknowledge messages.
func WithAutoAck() Option {
	return func(client *Client) {
		client.autoAck = true
	}
}

// WithSyncPublisher configures the client to publish messages synchronously.
func WithSyncPublisher() Option {
	return func(client *Client) {
		client.syncPublisher = true
	}
}

// WithLogger sets a custom logger for the Client.
func WithLogger(log logger.Logger) Option {
	return func(client *Client) {
		client.log = log
	}
}

// WithLogger sets a custom logger for the Client.
func WithClientOption(co ...option.ClientOption) Option {
	return func(client *Client) {
		client.clientOptions = co
	}
}
