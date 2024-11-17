package googlepubsub

import (
	"context"
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

func TestWithAutoAck(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
		mockFn  func() (*Client, error)
	}{
		{
			name:    "Success",
			wantErr: nil,
			mockFn: func() (*Client, error) {
				return NewClient(context.Background(), "pubsub", WithAutoAck())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn()
			assert.Equal(t, tt.wantErr, err)
			assert.NotNil(t, got)
		})
	}
}

func TestWithSyncPublisher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
		mockFn  func() (*Client, error)
	}{
		{
			name:    "Success",
			wantErr: nil,
			mockFn: func() (*Client, error) {
				return NewClient(context.Background(), "pubsub", WithSyncPublisher())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn()
			assert.Equal(t, tt.wantErr, err)
			assert.NotNil(t, got)
		})
	}
}

func TestWithLogger(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
		mockFn  func() (*Client, error)
	}{
		{
			name:    "Success",
			wantErr: nil,
			mockFn: func() (*Client, error) {
				return NewClient(context.Background(), "pubsub", WithLogger(logger.NewNoopLogger()))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn()
			assert.Equal(t, tt.wantErr, err)
			assert.NotNil(t, got)
		})
	}
}

func TestWithClientOption(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
		mockFn  func() (*Client, error)
	}{
		{
			name:    "Success",
			wantErr: nil,
			mockFn: func() (*Client, error) {
				return NewClient(context.Background(), "pubsub", WithClientOption(option.WithAudiences("ok")))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn()
			assert.Equal(t, tt.wantErr, err)
			assert.NotNil(t, got)
		})
	}
}
