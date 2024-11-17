package redispubsub

import (
	"github.com/redis/go-redis/v9"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
)

// Option defines a configuration function for the Client.
type Option func(*Client)

// WithLogger sets a custom logger for the Client.
func WithLogger(log logger.Logger) Option {
	return func(client *Client) {
		client.log = log
	}
}

// WithSyncPublisher enables or disables synchronous publishing.
func WithSyncPublisher() Option {
	return func(client *Client) {
		client.syncPublisher = true
	}
}

// WithExistingClient uses an existing Redis client.
func WithExistingClient(redisClient redis.UniversalClient) Option {
	return func(client *Client) {
		client.client = redisClient
	}
}
