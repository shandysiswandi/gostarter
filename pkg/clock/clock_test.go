package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *TimeClocker
	}{
		{
			name: "Success",
			want: &TimeClocker{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := New()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeClocker_Now(t *testing.T) {
	tests := []struct {
		name string
		tr   *TimeClocker
		want time.Time
	}{
		{
			name: "Success",
			tr:   &TimeClocker{},
			want: time.Now(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.tr.Now()
			assert.Equal(t, tt.want.Truncate(time.Second), got.Truncate(time.Second))
		})
	}
}
