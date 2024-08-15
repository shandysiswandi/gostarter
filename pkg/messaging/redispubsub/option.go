// Package redispubsub provides a client for interacting with Redis Pub/Sub.
// It includes functionality for publishing and subscribing to messages with various configuration options.
package redispubsub

import "github.com/redis/go-redis/v9"

// Option defines a function type used to configure the Pub/Sub client.
type Option func(*Client)

// WithSyncPublisher configures the client to publish messages synchronously.
//
// isSyncPublisher: If true, the Publish call will wait for the message to be acknowledged before returning.
func WithSyncPublisher(isSyncPublisher bool) Option {
	return func(client *Client) {
		client.syncPublisher = isSyncPublisher
	}
}

// WithExistsClient uses an existing Redis client.
//
// ruc: An instance of redis.UniversalClient. This option allows you to provide
// your own Redis client to the Pub/Sub client, useful for custom configurations
// or sharing the client across different parts of your application.
func WithExistsClient(ruc redis.UniversalClient) Option {
	return func(client *Client) {
		client.client = ruc
	}
}
