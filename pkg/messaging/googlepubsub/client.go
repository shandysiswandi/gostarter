// Package googlepubsub provides a client for interacting with Google Cloud Pub/Sub.
// It includes functionality for publishing and subscribing to messages with various configuration options.
package googlepubsub

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/shandysiswandi/gostarter/pkg/messaging"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"google.golang.org/api/option"
)

var (
	// ErrInactiveClient is returned when an operation is attempted with a closed or inactive client.
	ErrInactiveClient = errors.New("inactive client")

	// ErrTopicNotFound is returned when attempting to access a non-existent topic.
	ErrTopicNotFound = errors.New("topic does not exist")

	// ErrSubscriptionNotFound is returned when attempting to access a non-existent subscription.
	ErrSubscriptionNotFound = errors.New("subscription does not exist")
)

// Client represents a Google Cloud Pub/Sub client with additional configuration options.
type Client struct {
	clientOptions []option.ClientOption // Options for the Pub/Sub client.
	autoAck       bool                  // Whether to automatically acknowledge messages.
	syncPublisher bool                  // Whether to use synchronous publishing.
	client        *pubsub.Client
	log           logger.Logger
	cancels       []context.CancelFunc
	wg            sync.WaitGroup
}

// NewClient creates a new Google Cloud Pub/Sub client with the provided configuration options.
func NewClient(ctx context.Context, projectID string, opts ...Option) (*Client, error) {
	client := &Client{
		clientOptions: []option.ClientOption{},
		autoAck:       false,
		syncPublisher: false,
		log:           logger.NewNoopLogger(),
	}

	for _, opt := range opts {
		opt(client)
	}

	cli, err := pubsub.NewClient(ctx, projectID, client.clientOptions...)
	if err != nil {
		return nil, err
	}
	client.client = cli

	return client, nil
}

// Close closes the Pub/Sub client and releases any resources associated with it.
func (c *Client) Close() error {
	if c.client == nil {
		return nil
	}

	for _, cancel := range c.cancels {
		cancel()
	}

	c.wg.Wait()

	return c.client.Close()
}

// Publish sends a single message to the specified topic.
func (c *Client) Publish(ctx context.Context, topic string, data *messaging.Data) error {
	if c.client == nil {
		return ErrInactiveClient
	}

	if err := data.Validate(); err != nil {
		return err
	}

	t, err := c.getTopic(ctx, topic)
	if err != nil {
		return err
	}

	msg := &pubsub.Message{Data: data.Msg, Attributes: data.Attributes}
	if c.syncPublisher {
		_, err := t.Publish(ctx, msg).Get(ctx)

		return err
	}

	go func() {
		_, err = t.Publish(ctx, msg).Get(ctx)
		if err != nil {
			c.log.Error(ctx, "async publish failed", err)
		}
	}()

	return nil
}

// BulkPublish sends multiple messages to the specified topic.
func (c *Client) BulkPublish(ctx context.Context, topic string, datas []*messaging.Data) error {
	if c.client == nil {
		return ErrInactiveClient
	}

	t, err := c.getTopic(ctx, topic)
	if err != nil {
		return err
	}

	if c.syncPublisher {
		return c.doBulkPublish(ctx, t, datas)
	}

	go func() {
		err := c.doBulkPublish(ctx, t, datas)
		if err != nil {
			c.log.Error(ctx, "async bulk publish failed", err)
		}
	}()

	return nil
}

// Subscribe subscribes to a specified topic with a given subscription ID and handler function.
func (c *Client) Subscribe(ctx context.Context, topic, subID string, h messaging.SubscriberFunc) error {
	if c.client == nil {
		return ErrInactiveClient
	}

	_, err := c.getTopic(ctx, topic)
	if err != nil {
		return err
	}

	subscription := c.client.Subscription(subID)
	exists, err := subscription.Exists(ctx)
	if err != nil {
		return err
	}

	if !exists {
		return ErrSubscriptionNotFound
	}

	subCtx, cancel := context.WithCancel(ctx)
	c.cancels = append(c.cancels, cancel)
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()

		defer func() {
			if r := recover(); r != nil {
				c.log.Error(subCtx, "recovered from subscriber panic", nil, logger.KeyVal("cause", r))
				debug.PrintStack()
			}
		}()

		f := func(subCtx context.Context, m *pubsub.Message) {
			err := h(subCtx, &messaging.Data{Msg: m.Data, Attributes: m.Attributes})
			c.doMessageResult(subCtx, topic, subID, err, m)
		}

		err := subscription.Receive(subCtx, f)
		if err != nil {
			c.log.Error(subCtx, "failed when receive message", err)
		}
	}()

	return nil
}

// getTopic retrieves a reference to a Pub/Sub topic and checks its existence.
func (c *Client) getTopic(ctx context.Context, topic string) (*pubsub.Topic, error) {
	t := c.client.Topic(topic)

	exists, err := t.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrTopicNotFound
	}

	return t, nil
}

// doBulkPublish sends multiple messages to the specified topic.
func (c *Client) doBulkPublish(ctx context.Context, topic *pubsub.Topic, datas []*messaging.Data) error {
	for i, data := range datas {
		if err := data.Validate(); err != nil {
			return err
		}

		result := topic.Publish(ctx, &pubsub.Message{Data: data.Msg, Attributes: data.Attributes})
		if _, err := result.Get(ctx); err != nil {
			return fmt.Errorf("failed to publish message at index %d: %w", i, err)
		}
	}

	return nil
}

// processMessageResult handles message acknowledgment based on processing result.
func (c *Client) doMessageResult(ctx context.Context, topic, subID string, err error, m *pubsub.Message) {
	if err != nil {
		c.log.Error(ctx, "message handler failed", err,
			logger.KeyVal("topic", topic),
			logger.KeyVal("subscription", subID),
		)
		if !c.autoAck {
			m.Nack()

			return
		}
	}

	m.Ack()
}
