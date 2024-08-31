package emitter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEvent_GetTopic(t *testing.T) {
	tests := []struct {
		name   string
		event  Event
		expect string
	}{
		{
			name:   "Valid topic",
			event:  Event{topic: "UserSignup"},
			expect: "UserSignup",
		},
		{
			name:   "Empty topic",
			event:  Event{topic: ""},
			expect: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.event.GetTopic()
			assert.Equal(t, tt.expect, got, "GetTopic() should return the correct topic")
		})
	}
}

func TestEvent_GetArgs(t *testing.T) {
	tests := []struct {
		name   string
		event  Event
		expect []any
	}{
		{
			name:   "Non-empty args",
			event:  Event{args: []any{1, "foo", 3.14}},
			expect: []any{1, "foo", 3.14},
		},
		{
			name:   "Empty args",
			event:  Event{args: []any{}},
			expect: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.event.GetArgs()
			assert.Equal(t, tt.expect, got, "GetArgs() should return the correct arguments")
		})
	}
}

func TestEvent_GetTimestamp(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		event    Event
		expected time.Time
	}{
		{
			name:     "Valid timestamp",
			event:    Event{timestamp: now},
			expected: now,
		},
		{
			name:     "Zero timestamp",
			event:    Event{timestamp: time.Time{}},
			expected: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.event.GetTimestamp()
			assert.Equal(t, tt.expected, got, "GetTimestamp() should return the correct timestamp")
		})
	}
}
