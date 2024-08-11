package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClock(t *testing.T) {
	tests := []struct {
		name string
		want *Clock
	}{
		{
			name: "BumpSuccess",
			want: &Clock{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewClock()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClock_Now(t *testing.T) {
	tests := []struct {
		name string
		c    *Clock
		want time.Time
	}{
		{
			name: "Success",
			c:    &Clock{},
			want: time.Now(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.c.Now()
			assert.Equal(t, tt.want.Format(time.RFC3339), got.Format(time.RFC3339))
		})
	}
}
