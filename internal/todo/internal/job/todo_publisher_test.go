package job

import (
	"context"
	"testing"
	"time"

	"github.com/shandysiswandi/goreng/mocker"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_todoPublisher_Start(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
		mockFn  func() *todoPublisher
	}{
		{
			name:    "ErrorEncode",
			wantErr: nil,
			mockFn: func() *todoPublisher {
				tel := telemetry.NewTelemetry()
				jsonMock := mocker.NewMockCodec(t)
				msgMock := mocker.NewMockMessagingClient(t)

				jsonMock.EXPECT().
					Encode(mock.Anything).
					Return(nil, assert.AnError).
					Once()

				return &todoPublisher{
					cjson: jsonMock,
					mc:    msgMock,
					tel:   tel,
					topic: "topic",
				}
			},
		},
		{
			name:    "ErrorPublish",
			wantErr: nil,
			mockFn: func() *todoPublisher {
				tel := telemetry.NewTelemetry()
				jsonMock := mocker.NewMockCodec(t)
				msgMock := mocker.NewMockMessagingClient(t)

				jsonMock.EXPECT().
					Encode(mock.Anything).
					Return([]byte{}, nil).
					Times(10)

				msgMock.EXPECT().
					Publish(mock.Anything, mock.Anything, mock.Anything).
					Return(assert.AnError)

				return &todoPublisher{
					cjson: jsonMock,
					mc:    msgMock,
					tel:   tel,
					topic: "topic",
				}
			},
		},
		{
			name:    "ErrorBulkPublish",
			wantErr: nil,
			mockFn: func() *todoPublisher {
				tel := telemetry.NewTelemetry()
				jsonMock := mocker.NewMockCodec(t)
				msgMock := mocker.NewMockMessagingClient(t)

				jsonMock.EXPECT().
					Encode(mock.Anything).
					Return([]byte{}, nil).
					Times(10)

				msgMock.EXPECT().
					Publish(mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				msgMock.EXPECT().
					BulkPublish(mock.Anything, mock.Anything, mock.Anything).
					Return(assert.AnError)

				return &todoPublisher{
					cjson: jsonMock,
					mc:    msgMock,
					tel:   tel,
					topic: "topic",
				}
			},
		},
		{
			name:    "Success",
			wantErr: nil,
			mockFn: func() *todoPublisher {
				tel := telemetry.NewTelemetry()
				jsonMock := mocker.NewMockCodec(t)
				msgMock := mocker.NewMockMessagingClient(t)

				jsonMock.EXPECT().
					Encode(mock.Anything).
					Return([]byte{}, nil).
					Times(10)

				msgMock.EXPECT().
					Publish(mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				msgMock.EXPECT().
					BulkPublish(mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				return &todoPublisher{
					cjson: jsonMock,
					mc:    msgMock,
					tel:   tel,
					topic: "topic",
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.mockFn().Start()
			assert.Equal(t, tt.wantErr, err)
			time.Sleep(time.Second) //  wait actual Goroutine
		})
	}
}

func Test_todoPublisher_Stop(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr error
		mockFn  func() *todoPublisher
	}{
		{
			name:    "Success",
			ctx:     context.Background(),
			wantErr: nil,
			mockFn: func() *todoPublisher {
				return &todoPublisher{
					tel: telemetry.NewTelemetry(),
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.mockFn().Stop(tt.ctx)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
