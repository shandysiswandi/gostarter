// Package googlepubsub provides a client for interacting with Google Cloud Pub/Sub.
// It includes functionality for publishing and subscribing to messages with various configuration options.
package googlepubsub

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/shandysiswandi/gostarter/pkg/messaging"
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
	clientOptions        []option.ClientOption // Options for the Pub/Sub client.
	autoAck              bool                  // Whether to automatically acknowledge messages.
	autoCreateTopic      bool                  // Automatically create topic if not existing.
	autoCreateSubscriber bool                  // Automatically create subscription if not existing.
	syncPublisher        bool                  // Whether to use synchronous publishing.
	client               *pubsub.Client
	subscriptions        map[string]*SubscriberHandler
	mu                   sync.RWMutex
}

// NewClient creates a new Google Cloud Pub/Sub client with the provided configuration options.
//
// ctx: The context to use for client initialization.
// projectID: The Google Cloud project ID.
// opts: Configuration options for the client.
//
// Returns a pointer to the created Client or an error if the client could not be created.
func NewClient(ctx context.Context, projectID string, opts ...Option) (*Client, error) {
	client := &Client{
		subscriptions:        make(map[string]*SubscriberHandler),
		clientOptions:        []option.ClientOption{},
		autoAck:              false,
		autoCreateTopic:      false,
		autoCreateSubscriber: false,
		syncPublisher:        false,
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
//
// Returns an error if the client could not be closed.
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
//
// ctx: The context to use for the publish operation.
// topic: The name of the Pub/Sub topic.
// message: The message to be published.
//
// Returns an error if the message could not be published.
func (c *Client) Publish(ctx context.Context, topic string, message []byte) error {
	if c.client == nil {
		return ErrInactiveClient
	}

	t, err := c.getTopic(ctx, topic)
	if err != nil {
		return err
	}

	msg := &pubsub.Message{Data: message, Attributes: make(map[string]string)}
	if c.syncPublisher {
		_, err := t.Publish(ctx, msg).Get(ctx)

		return err
	}

	_ = t.Publish(ctx, msg)

	return nil
}

// BulkPublish sends multiple messages to the specified topic.
//
// ctx: The context to use for the publish operation.
// topic: The name of the Pub/Sub topic.
// messages: A slice of messages to be published.
//
// Returns an error if any of the messages could not be published.
func (c *Client) BulkPublish(ctx context.Context, topic string, messages [][]byte) error {
	if c.client == nil {
		return ErrInactiveClient
	}

	t, err := c.getTopic(ctx, topic)
	if err != nil {
		return err
	}

	if c.syncPublisher {
		return c.doSyncBulkPublish(ctx, t, messages)
	}

	return c.doAsyncBulkPublish(ctx, t, messages)
}

// Subscribe subscribes to a specified topic with a given subscription ID and handler function.
//
// ctx: The context to use for the subscription operation.
// topic: The name of the Pub/Sub topic.
// subscriptionID: The ID for the subscription.
// handler: The function to handle incoming messages.
//
// Returns a SubscriptionHandler to manage the subscription and an error if the subscription could not be created.
//
//nolint:ireturn // ignore this linter in this file
func (c *Client) Subscribe(ctx context.Context, topic, subscriptionID string, handler messaging.SubscriberHandlerFunc) (
	messaging.SubscriptionHandler, error,
) {
	if c.client == nil {
		return nil, ErrInactiveClient
	}

	_, err := c.getTopic(ctx, topic)
	if err != nil {
		return nil, err
	}

	subscription := c.client.Subscription(subscriptionID)
	exists, err := subscription.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrSubscriptionNotFound
	}

	ctx, cancel := context.WithCancel(ctx)
	sh := &SubscriberHandler{cancelFunc: cancel, client: c, name: subscriptionID}
	c.addHandler(subscriptionID, sh)

	go c.subscribingMessage(ctx, subscription, topic, subscriptionID, handler)

	return sh, nil
}

// getTopic retrieves a reference to a Pub/Sub topic and checks its existence.
//
// ctx: The context to use for the topic operation.
// topic: The name of the Pub/Sub topic.
//
// Returns a reference to the topic and an error if the topic does not exist or another error occurred.
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

// doSyncBulkPublish sends multiple messages to the specified topic synchronously.
//
// ctx: The context to use for the publish operation.
// topic: The Pub/Sub topic to publish messages to.
// messages: A slice of messages to be published.
//
// Returns an error if any of the messages could not be published.
func (c *Client) doSyncBulkPublish(ctx context.Context, topic *pubsub.Topic, messages [][]byte) error {
	for i, msg := range messages {
		pubMsg := &pubsub.Message{Data: msg, Attributes: make(map[string]string)}

		result := topic.Publish(ctx, pubMsg)
		if _, err := result.Get(ctx); err != nil {
			return fmt.Errorf("failed to publish message at index %d: %w", i, err)
		}
	}

	return nil
}

// doAsyncBulkPublish sends multiple messages to the specified topic asynchronously.
//
// ctx: The context to use for the publish operation.
// topic: The Pub/Sub topic to publish messages to.
// messages: A slice of messages to be published.
//
// Returns an error if any of the messages could not be published.
func (c *Client) doAsyncBulkPublish(ctx context.Context, topic *pubsub.Topic, messages [][]byte) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var publishErrors []error
	errChan := make(chan error, len(messages)) // Buffer size to hold all possible errors

	// Start a goroutine to handle logging errors
	go func() {
		for err := range errChan {
			if err != nil {
				mu.Lock()
				publishErrors = append(publishErrors, err) // Collect errors
				mu.Unlock()
				log.Printf("Publishing error: %v", err) // Log the error
			}
		}
	}()

	for _, msg := range messages {
		wg.Add(1)
		go func(m []byte) {
			defer wg.Done()
			res := topic.Publish(ctx, &pubsub.Message{Data: m})
			if _, err := res.Get(ctx); err != nil {
				errChan <- err // Send error to channel
			}
		}(msg)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan) // Close error channel after all goroutines are done

	// Check if any errors were collected
	//nolint:godox,err113 // will add later
	if len(publishErrors) > 0 {
		// TODO: here to handle retry with configuration, for now only log all errors
		// Return an aggregate error if there were any issues
		return fmt.Errorf("errors occurred during bulk publish: %v", publishErrors)
	}

	return nil
}

// addHandler adds a new SubscriberHandler to the client's subscription management.
//
// key: The unique identifier for the subscription handler.
// handler: The SubscriberHandler instance to be added.
//
// This method is thread-safe, as it uses a mutex to synchronize access to the subscriptions map.
func (c *Client) addHandler(key string, handler *SubscriberHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.subscriptions[key] = handler
}

// removeHandler removes a SubscriberHandler from the client's subscription management.
//
// key: The unique identifier for the subscription handler to be removed.
//
// This method is thread-safe, as it uses a mutex to synchronize access to the subscriptions map.
func (c *Client) removeHandler(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.subscriptions, key)
}

// subscribingMessage handles the receiving and processing of messages from a given Pub/Sub subscription.
//
// ctx: The context used to control cancellation and timeouts for receiving messages.
// subs: A pointer to the pubsub.Subscription object representing the Pub/Sub subscription.
// topic: The name of the topic associated with the subscription.
// subscription: The ID of the subscription within the Pub/Sub system.
// handler: A messaging.Handler function that processes the message data.
func (c *Client) subscribingMessage(ctx context.Context, subs *pubsub.Subscription, topic, subscription string,
	handler messaging.SubscriberHandlerFunc,
) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("subscribingMessage recover: %v", r)
			debug.PrintStack()
		}
	}()

	f := func(ctx context.Context, m *pubsub.Message) {
		err := handler(ctx, topic, subscription, m.Data)
		if err != nil {
			log.Printf("Failed to handle message topic(%s) subscription(%s): %v", topic, subscription, err)
		}
	}

	err := subs.Receive(ctx, f)
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("Subscription receive error: %v", err)
	}
}
