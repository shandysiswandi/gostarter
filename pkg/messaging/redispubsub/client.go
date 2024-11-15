// Package redispubsub provides a client for interacting with Redis Pub/Sub.
// It includes functionality for publishing and subscribing to messages with various configuration options.
package redispubsub

import (
	"context"
	"errors"
	"log"
	"runtime/debug"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/shandysiswandi/gostarter/pkg/messaging"
)

// ErrInactiveClient is returned when an operation is attempted with a closed or inactive client.
var ErrInactiveClient = errors.New("inactive client")

// Client provides methods for interacting with Redis Pub/Sub, including publishing messages
// and subscribing to topics. It supports both synchronous and asynchronous publishing.
type Client struct {
	syncPublisher bool                          // Whether to use synchronous publishing.
	client        redis.UniversalClient         // The Redis client.
	subscriptions map[string]*SubscriberHandler // A map of active subscriptions.
	mu            sync.RWMutex                  // Mutex for protecting access to subscriptions.
}

// WithSyncPublisher configures the client to publish messages synchronously.
func WithSyncPublisher(isSyncPublisher bool) func(*Client) {
	return func(client *Client) {
		client.syncPublisher = isSyncPublisher
	}
}

// WithExistsClient uses an existing Redis client.
func WithExistsClient(ruc redis.UniversalClient) func(*Client) {
	return func(client *Client) {
		client.client = ruc
	}
}

// NewRedisClient creates a new Client instance configured with the given Redis URL and options.
func NewRedisClient(url string, opts ...func(*Client)) (*Client, error) {
	client := &Client{
		subscriptions: make(map[string]*SubscriberHandler),
	}

	for _, opt := range opts {
		opt(client)
	}

	if client.client == nil {
		cli := redis.NewClient(&redis.Options{Addr: url})
		client.client = cli
	}

	if err := client.client.Ping(context.Background()).Err(); err != nil {
		return nil, ErrInactiveClient
	}

	return client, nil
}

// addHandler adds a new SubscriberHandler to the client's subscriptions map.
func (c *Client) addHandler(key string, handler *SubscriberHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.subscriptions[key] = handler
}

// removeHandler removes a SubscriberHandler from the client's subscriptions map.
func (c *Client) removeHandler(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.subscriptions, key)
}

// goroutine is a helper function that wraps a function in a defer-recover block
// to handle any panics that occur during the execution of the function.
func (c *Client) goroutine(f func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("goroutine recover: %v", r)
			debug.PrintStack()
		}
	}()

	f()
}

// Close terminates all active subscriptions and closes the Redis client connection.
func (c *Client) Close() error {
	if c.client == nil {
		return nil
	}

	for _, sub := range c.subscriptions {
		_ = sub.Close()
	}

	if err := c.client.Close(); err != nil {
		return err
	}

	c.client = nil

	return nil
}

// Publish sends a single message to the specified topic.
// It supports both synchronous and asynchronous publishing based on the client's configuration.
func (c *Client) Publish(ctx context.Context, topic string, msg []byte) error {
	if c.client == nil {
		return ErrInactiveClient
	}

	if c.syncPublisher {
		return c.client.Publish(ctx, topic, msg).Err()
	}

	go c.goroutine(func() { c.client.Publish(ctx, topic, msg) })

	return nil
}

// BulkPublish sends multiple messages to the specified topic using a Redis pipeline.
// It can operate in either synchronous or asynchronous mode based on the syncPublisher setting.
func (c *Client) BulkPublish(ctx context.Context, topic string, messages [][]byte) error {
	if c.client == nil {
		return ErrInactiveClient
	}

	pipe := c.client.Pipeline()
	for _, message := range messages {
		pipe.Publish(ctx, topic, message)
	}

	if c.syncPublisher {
		_, err := pipe.Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	go c.goroutine(func() {
		_, err := pipe.Exec(ctx)
		if err != nil {
			log.Println("bulk publish failed:", err)
		}
	})

	return nil
}

// Subscribe subscribes to a specified topic with a given subscription ID and handler function.
// The handler is called whenever a message is received on the topic.
func (c *Client) Subscribe(ctx context.Context, topic, subID string, h messaging.SubscriberHandlerFunc) (
	messaging.SubscriptionHandler, error,
) {
	if c.client == nil {
		return nil, ErrInactiveClient
	}

	pubsub := c.client.Subscribe(ctx, topic)
	ctx, cancel := context.WithCancel(ctx)
	sh := &SubscriberHandler{cancelFunc: cancel, client: c, name: subID}
	c.addHandler(subID, sh)

	go func(ps *redis.PubSub) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("subscriber recover: %v", r)
				debug.PrintStack()
			}
		}()

		for {
			select {
			case msg := <-ps.Channel():
				err := h(ctx, topic, subID, []byte(msg.Payload))
				if err != nil {
					log.Printf("Failed to handle message topic(%s) subscription(%s): %v", topic, subID, err)
				}

			case <-ctx.Done():
				return
			}
		}
	}(pubsub)

	return sh, nil
}
