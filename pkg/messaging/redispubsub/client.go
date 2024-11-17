// Package redispubsub provides a client for interacting with Redis Pub/Sub.
package redispubsub

import (
	"context"
	"errors"
	"runtime/debug"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/shandysiswandi/gostarter/pkg/messaging"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
)

// ErrInactiveClient indicates a closed or inactive Redis client.
var ErrInactiveClient = errors.New("inactive client")

// Client represents a Redis Pub/Sub client with support for synchronous and asynchronous publishing.
type Client struct {
	syncPublisher bool
	client        redis.UniversalClient
	log           logger.Logger
	cancels       []context.CancelFunc
	wg            sync.WaitGroup
}

// NewClient initializes a new Redis Pub/Sub client with the given URL and options.
func NewClient(url string, opts ...Option) (*Client, error) {
	client := &Client{log: logger.NewNoopLogger()}

	for _, opt := range opts {
		opt(client)
	}

	if client.client == nil {
		client.client = redis.NewClient(&redis.Options{Addr: url})
	}

	if err := client.client.Ping(context.Background()).Err(); err != nil {
		return nil, ErrInactiveClient
	}

	return client, nil
}

// Close terminates all active subscriptions and closes the Redis client connection.
func (rc *Client) Close() error {
	if rc.client == nil {
		return nil
	}

	for _, cancel := range rc.cancels {
		cancel()
	}

	rc.wg.Wait()

	return rc.client.Close()
}

// Publish sends a message to a specific topic. If syncPublisher is enabled, publishing is synchronous.
func (rc *Client) Publish(ctx context.Context, topic string, data *messaging.Data) error {
	if rc.client == nil {
		return ErrInactiveClient
	}

	if err := data.Validate(); err != nil {
		return err
	}

	if rc.syncPublisher {
		return rc.client.Publish(ctx, topic, data.Msg).Err()
	}

	go func() {
		if err := rc.client.Publish(ctx, topic, data.Msg).Err(); err != nil {
			rc.log.Error(ctx, "async publish failed", err)
		}
	}()

	return nil
}

// BulkPublish sends multiple messages to a topic using a Redis pipeline.
func (rc *Client) BulkPublish(ctx context.Context, topic string, datas []*messaging.Data) error {
	if rc.client == nil {
		return ErrInactiveClient
	}

	pipe := rc.client.Pipeline()
	for _, data := range datas {
		if err := data.Validate(); err != nil {
			return err
		}
		pipe.Publish(ctx, topic, data.Msg)
	}

	if rc.syncPublisher {
		_, err := pipe.Exec(ctx)

		return err
	}

	go func() {
		if _, err := pipe.Exec(ctx); err != nil {
			rc.log.Error(ctx, "async bulk publish failed", err)
		}
	}()

	return nil
}

// Subscribe subscribes to a topic and calls the handler when a message is received.
func (rc *Client) Subscribe(ctx context.Context, topic, subID string, h messaging.SubscriberFunc) error {
	if rc.client == nil {
		return ErrInactiveClient
	}

	subCtx, cancel := context.WithCancel(ctx)
	rc.cancels = append(rc.cancels, cancel)
	rc.wg.Add(1)

	go func() {
		defer rc.wg.Done()
		pubsub := rc.client.Subscribe(subCtx, topic)

		defer func() {
			if r := recover(); r != nil {
				rc.log.Error(subCtx, "recovered from subscriber panic", nil, logger.KeyVal("cause", r))
				debug.PrintStack()
			}
		}()

		for {
			select {
			case msg := <-pubsub.Channel():
				if err := h(subCtx, &messaging.Data{Msg: []byte(msg.Payload)}); err != nil {
					rc.log.Error(subCtx, "message handler failed", err,
						logger.KeyVal("topic", topic),
						logger.KeyVal("subscription", subID),
					)
				}
			case <-subCtx.Done():
				return
			}
		}
	}()

	return nil
}
