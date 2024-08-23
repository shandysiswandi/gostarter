package emitter

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEventEmitter(t *testing.T) {
	tests := []struct {
		name                    string
		options                 []Option
		inputArgsValidation     []any
		inputTopicValidation    string
		expectedArgsValidation  error
		expectedTopicValidation error
	}{
		{
			name:                    "Default options",
			options:                 nil,
			inputArgsValidation:     []any{},
			inputTopicValidation:    "",
			expectedArgsValidation:  nil,
			expectedTopicValidation: nil,
		},
		{
			name: "Custom all options",
			options: []Option{
				WithArgumentsValidation(func(args []any) error {
					if len(args) < 3 {
						return assert.AnError
					}
					return nil
				}),
				WithTopicValidation(func(topic string) error {
					if topic == "fail" {
						return assert.AnError
					}
					return nil
				}),
				WithBufferSize(3),
				WithTimeProvider(nil),
			},
			inputArgsValidation:     []any{"1", "2"},
			inputTopicValidation:    "fail",
			expectedArgsValidation:  assert.AnError,
			expectedTopicValidation: assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewEventEmitter(tt.options...)
			assert.NotNil(t, got)
			assert.Equal(t, tt.expectedArgsValidation, got.options.ArgumentsValidation(tt.inputArgsValidation))
			assert.Equal(t, tt.expectedTopicValidation, got.options.TopicValidation(tt.inputTopicValidation))
		})
	}
}

func TestEventEmitter_Emit(t *testing.T) {
	tests := []struct {
		name          string
		topic         string
		args          []any
		expectedError error
		mockFn        func(topic string, args ...any) *EventEmitter
	}{
		{
			name:          "No validation",
			topic:         "testTopic",
			args:          []any{"arg1", 2},
			expectedError: nil,
			mockFn: func(topic string, args ...any) *EventEmitter {
				return &EventEmitter{}
			},
		},
		{
			name:          "Invalid topic validation",
			topic:         "testTopic",
			args:          []any{"arg1", 2},
			expectedError: assert.AnError,
			mockFn: func(topic string, args ...any) *EventEmitter {
				return &EventEmitter{
					options: Options{
						TopicValidation:     func(topic string) error { return assert.AnError },
						ArgumentsValidation: func(a []any) error { return nil },
					},
				}
			},
		},
		{
			name:          "Invalid args validation",
			topic:         "testTopic",
			args:          []any{"arg1", 2},
			expectedError: assert.AnError,
			mockFn: func(topic string, args ...any) *EventEmitter {
				return &EventEmitter{
					options: Options{
						TopicValidation:     func(topic string) error { return nil },
						ArgumentsValidation: func(args []any) error { return assert.AnError },
					},
				}
			},
		},
		{
			name:          "No listeners",
			topic:         "testTopic",
			args:          []any{"arg1", 2},
			expectedError: nil,
			mockFn: func(topic string, args ...any) *EventEmitter {
				mTime := &MockTimeProvider{}
				mTime.EXPECT().Now().Return(time.Time{})

				return &EventEmitter{
					options: Options{
						TopicValidation:     func(topic string) error { return nil },
						ArgumentsValidation: func(args []any) error { return nil },
						TimeProvider:        mTime,
					},
				}
			},
		},
		{
			name:          "Channel full",
			topic:         "testTopic",
			args:          []any{"arg1", 2},
			expectedError: nil,
			mockFn: func(topic string, args ...any) *EventEmitter {
				mTime := &MockTimeProvider{}
				mTime.EXPECT().Now().Return(time.Time{})

				listeners := make(map[string][]chan Event)
				listenerCh := make(chan Event, 1)
				listenerCh <- Event{} // Fill the channel
				listeners["testTopic"] = []chan Event{listenerCh}

				return &EventEmitter{
					listeners: listeners,
					options: Options{
						TopicValidation:     func(topic string) error { return nil },
						ArgumentsValidation: func(args []any) error { return nil },
						TimeProvider:        mTime,
					},
				}
			},
		},
		{
			name:          "Success",
			topic:         "testTopic",
			args:          []any{"arg1", 2},
			expectedError: nil,
			mockFn: func(topic string, args ...any) *EventEmitter {
				mTime := &MockTimeProvider{}
				mTime.EXPECT().Now().Return(time.Time{})

				listeners := make(map[string][]chan Event)
				listenerCh := make(chan Event, 2)
				listenerCh <- Event{}
				listeners["testTopic"] = []chan Event{listenerCh}

				go func() {
					event := <-listenerCh
					log.Println("success event", event)
				}()

				return &EventEmitter{
					listeners: listeners,
					options: Options{
						TopicValidation:     func(topic string) error { return nil },
						ArgumentsValidation: func(args []any) error { return nil },
						TimeProvider:        mTime,
					},
				}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ee := tt.mockFn(tt.topic, tt.args...)
			err := ee.Emit(tt.topic, tt.args...)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEventEmitter_AddListener(t *testing.T) {
	tests := []struct {
		name         string
		bufferSize   int
		expectedSize int
	}{
		{
			name:         "Default buffer size",
			bufferSize:   0,
			expectedSize: 3, // Default buffer size
		},
		{
			name:         "Custom buffer size",
			bufferSize:   5,
			expectedSize: 5,
		},
		{
			name:         "Add multiple listeners",
			bufferSize:   3,
			expectedSize: 3,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ee := &EventEmitter{
				options: Options{BufferSize: tt.bufferSize},
			}

			ch1 := ee.AddListener("topic1")
			ch2 := ee.AddListener("topic1")
			ch3 := ee.AddListener("topic1")

			assert.Equal(t, tt.expectedSize, cap(ch1))
			assert.Equal(t, tt.expectedSize, cap(ch2))
			assert.Equal(t, tt.expectedSize, cap(ch3))

			ee.mu.RLock()
			listeners, ok := ee.listeners["topic1"]
			ee.mu.RUnlock()

			assert.True(t, ok)
			assert.Len(t, listeners, 3)
			assert.Equal(t, ch1, (<-chan Event)(listeners[0]))
			assert.Equal(t, ch2, (<-chan Event)(listeners[1]))
			assert.Equal(t, ch3, (<-chan Event)(listeners[2]))
		})
	}
}

func TestEventEmitter_RemoveListener(t *testing.T) {
	tests := []struct {
		name              string
		topic             string
		setupListeners    func() (*EventEmitter, []<-chan Event)
		expectedRemaining int
	}{
		{
			name:  "Remove all listeners for a topic",
			topic: "topic1",
			setupListeners: func() (*EventEmitter, []<-chan Event) {
				ee := &EventEmitter{
					listeners: make(map[string][]chan Event),
				}
				ee.AddListener("topic1")
				ee.AddListener("topic1")
				return ee, nil
			},
			expectedRemaining: 0,
		},
		{
			name:  "Remove specific listeners for a topic",
			topic: "topic1",
			setupListeners: func() (*EventEmitter, []<-chan Event) {
				ee := &EventEmitter{
					listeners: make(map[string][]chan Event),
				}
				ch1 := ee.AddListener("topic1")
				ee.AddListener("topic1")
				return ee, []<-chan Event{ch1}
			},
			expectedRemaining: 1,
		},
		{
			name:  "Remove listener from non-existent topic",
			topic: "nonexistent",
			setupListeners: func() (*EventEmitter, []<-chan Event) {
				ee := &EventEmitter{
					listeners: make(map[string][]chan Event),
				}
				return ee, nil
			},
			expectedRemaining: 0,
		},
		{
			name:  "Remove specific listeners and leave others",
			topic: "topic1",
			setupListeners: func() (*EventEmitter, []<-chan Event) {
				ee := &EventEmitter{
					listeners: make(map[string][]chan Event),
				}
				ch1 := ee.AddListener("topic1")
				ee.AddListener("topic1")
				ch3 := ee.AddListener("topic1")
				return ee, []<-chan Event{ch1, ch3}
			},
			expectedRemaining: 1,
		},
		{
			name:  "Remove specific listeners and with no others",
			topic: "topic1",
			setupListeners: func() (*EventEmitter, []<-chan Event) {
				ee := &EventEmitter{
					listeners: make(map[string][]chan Event),
				}
				ch1 := ee.AddListener("topic1")
				ch2 := ee.AddListener("topic1")
				ch3 := ee.AddListener("topic1")
				return ee, []<-chan Event{ch1, ch2, ch3}
			},
			expectedRemaining: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ee, listenersToRemove := tt.setupListeners()

			// Remove listeners based on the test case
			ee.RemoveListener(tt.topic, listenersToRemove...)

			// Check the remaining listeners for the topic
			remainingListeners := ee.listeners[tt.topic]

			// Ensure the expected number of listeners remains
			assert.Equal(t, tt.expectedRemaining, len(remainingListeners))

			// If we're removing specific listeners, check that they're no longer in the list
			for _, removed := range listenersToRemove {
				for _, remaining := range remainingListeners {
					assert.NotEqual(t, remaining, removed, "Removed listener should not be present")
				}
			}

			// Check that the channels are properly closed
			if len(remainingListeners) == 0 {
				_, found := ee.listeners[tt.topic]
				assert.False(t, found, "Topic should be removed if no listeners remain")
			}
		})
	}
}

func TestEventEmitter_Close(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *EventEmitter
		expectError bool
	}{
		{
			name: "Close with listeners",
			setup: func() *EventEmitter {
				ee := &EventEmitter{
					listeners: make(map[string][]chan Event),
				}
				ee.AddListener("topic1")
				ee.AddListener("topic2")
				return ee
			},
			expectError: false,
		},
		{
			name: "Close with no listeners",
			setup: func() *EventEmitter {
				ee := &EventEmitter{
					listeners: make(map[string][]chan Event),
				}
				return ee
			},
			expectError: false,
		},
		{
			name: "Close an EventEmitter twice",
			setup: func() *EventEmitter {
				ee := &EventEmitter{
					listeners: make(map[string][]chan Event),
				}
				ee.AddListener("topic1")
				ee.Close()
				return ee
			},
			expectError: false,
		},
		{
			name: "Close an EventEmitter and check channel closure",
			setup: func() *EventEmitter {
				ee := &EventEmitter{
					listeners: make(map[string][]chan Event),
				}
				ch := ee.AddListener("topic1")
				ee.Close()

				// Attempt to receive from the closed channel to confirm it's closed.
				_, ok := <-ch
				assert.False(t, ok, "channel should be closed")

				return ee
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ee := tt.setup()

			err := ee.Close()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Nil(t, ee.listeners)
		})
	}
}
