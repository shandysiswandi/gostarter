package redispubsub

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/stretchr/testify/assert"
)

func TestWithLogger(t *testing.T) {
	type args struct {
		log logger.Logger
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(a args) (*Client, error)
		want    *Client
		wantErr error
	}{
		{
			name: "ErrorInactiveClient",
			args: args{
				log: logger.NewNoopLogger(),
			},
			mockFn: func(a args) (*Client, error) {
				return NewClient("-", WithLogger(a.log))
			},
			wantErr: ErrInactiveClient,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn(tt.args)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWithSyncPublisher(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() (*Client, error)
		want    *Client
		wantErr error
	}{
		{
			name: "ErrorInactiveClient",
			mockFn: func() (*Client, error) {
				return NewClient("-", WithSyncPublisher())
			},
			wantErr: ErrInactiveClient,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn()
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWithExistingClient(t *testing.T) {
	type args struct {
		redisClient redis.UniversalClient
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(a args) (*Client, error)
		want    *Client
		wantErr error
	}{
		{
			name: "ErrorInactiveClient",
			args: args{
				redisClient: nil,
			},
			mockFn: func(a args) (*Client, error) {
				return NewClient("-", WithExistingClient(a.redisClient))
			},
			wantErr: ErrInactiveClient,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn(tt.args)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
