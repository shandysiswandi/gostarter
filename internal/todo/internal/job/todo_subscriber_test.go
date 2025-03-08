package job

import (
	"context"
	"testing"

	"github.com/shandysiswandi/goreng/messaging"
	"github.com/shandysiswandi/goreng/mocker"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_todoSubscriber_Start(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
		mockFn  func() *todoSubscriber
	}{
		{
			name:    "ErrorSubscribe",
			wantErr: assert.AnError,
			mockFn: func() *todoSubscriber {
				tel := telemetry.NewTelemetry()
				msgMock := mocker.NewMockMessagingClient(t)

				msgMock.EXPECT().
					Subscribe(mock.Anything, "topic", "subscription", mock.Anything).
					Return(assert.AnError)

				return &todoSubscriber{
					tel:          tel,
					mc:           msgMock,
					topic:        "topic",
					subscription: "subscription",
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.mockFn().Start()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_todoSubscriber_do(t *testing.T) {
	type args struct {
		ctx  context.Context
		data *messaging.Data
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) *todoSubscriber
	}{
		{
			name: "ErrorDecode",
			args: args{
				ctx:  context.Background(),
				data: &messaging.Data{},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) *todoSubscriber {
				tel := telemetry.NewTelemetry()
				jsonMock := mocker.NewMockCodec(t)
				msgMock := mocker.NewMockMessagingClient(t)

				jsonMock.EXPECT().
					Decode(a.data.Msg, mock.Anything).
					Return(assert.AnError)

				return &todoSubscriber{
					cjson:        jsonMock,
					mc:           msgMock,
					tel:          tel,
					createUC:     nil,
					topic:        "",
					subscription: "",
				}
			},
		},
		{
			name: "ErrorDomainCreate",
			args: args{
				ctx:  context.Background(),
				data: &messaging.Data{},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) *todoSubscriber {
				tel := telemetry.NewTelemetry()
				jsonMock := mocker.NewMockCodec(t)
				createUC := mockz.NewMockCreate(t)

				jsonMock.EXPECT().
					Decode(a.data.Msg, mock.Anything).
					Return(nil)

				createUC.EXPECT().
					Call(a.ctx, mock.Anything).
					Return(nil, assert.AnError)

				return &todoSubscriber{
					cjson:    jsonMock,
					tel:      tel,
					createUC: createUC,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx:  context.Background(),
				data: &messaging.Data{},
			},
			wantErr: nil,
			mockFn: func(a args) *todoSubscriber {
				tel := telemetry.NewTelemetry()
				jsonMock := mocker.NewMockCodec(t)
				createUC := mockz.NewMockCreate(t)

				jsonMock.EXPECT().
					Decode(a.data.Msg, mock.Anything).
					Return(nil)

				createUC.EXPECT().
					Call(a.ctx, mock.Anything).
					Return(&domain.CreateOutput{}, nil)

				return &todoSubscriber{
					cjson:    jsonMock,
					tel:      tel,
					createUC: createUC,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.mockFn(tt.args).do(tt.args.ctx, tt.args.data)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_todoSubscriber_Stop(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr error
		mockFn  func() *todoSubscriber
	}{
		{
			name:    "Success",
			ctx:     context.Background(),
			wantErr: nil,
			mockFn: func() *todoSubscriber {
				return &todoSubscriber{
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
