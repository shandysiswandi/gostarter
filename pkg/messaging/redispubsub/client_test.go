package redispubsub

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/shandysiswandi/gostarter/pkg/messaging"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	type args struct {
		url  string
		opts func() []Option
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr error
	}{
		{
			name: "ErrorInactiveClient",
			args: args{
				url: "-",
				opts: func() []Option {
					opt := []Option{}

					opt = append(opt, WithSyncPublisher())

					return opt
				},
			},
			want:    nil,
			wantErr: ErrInactiveClient,
		},
		{
			name: "Success",
			args: args{
				url: "-",
				opts: func() []Option {
					opt := []Option{}
					db, mock := redismock.NewClientMock()
					mock.ExpectPing().SetVal("ok")
					opt = append(opt, WithExistingClient(db))

					return opt
				},
			},
			want: &Client{
				syncPublisher: false,
				client:        &redis.Client{},
				log:           logger.NewNoopLogger(),
				cancels:       nil,
				wg:            sync.WaitGroup{},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewClient(tt.args.url, tt.args.opts()...)
			assert.Equal(t, tt.wantErr, err)
			if tt.want != nil {
				assert.Equal(t, tt.want.cancels, got.cancels)
				assert.Equal(t, tt.want.syncPublisher, got.syncPublisher)
				// assert.Equal(t, tt.want.wg, got.wg)
				assert.Equal(t, tt.want.log, got.log)
				assert.NotNil(t, tt.want.client)
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
		mockFn  func() *Client
	}{
		{
			name:    "NoActiveClient",
			wantErr: nil,
			mockFn: func() *Client {
				return &Client{}
			},
		},
		{
			name:    "Success",
			wantErr: nil,
			mockFn: func() *Client {
				db, _ := redismock.NewClientMock()

				return &Client{client: db, cancels: []context.CancelFunc{func() {}}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.mockFn().Close()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestClient_Publish(t *testing.T) {
	type args struct {
		ctx   context.Context
		topic string
		data  *messaging.Data
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) *Client
	}{
		{
			name: "ErrorInactiveClient",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				data:  nil,
			},
			wantErr: ErrInactiveClient,
			mockFn: func(a args) *Client {
				return &Client{}
			},
		},
		{
			name: "InvalidDataNil",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				data:  nil,
			},
			wantErr: messaging.ErrDataNil,
			mockFn: func(a args) *Client {
				db, _ := redismock.NewClientMock()

				return &Client{
					client: db,
				}
			},
		},
		{
			name: "InvalidDataMessageNil",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				data:  &messaging.Data{},
			},
			wantErr: messaging.ErrMessageNil,
			mockFn: func(a args) *Client {
				db, _ := redismock.NewClientMock()

				return &Client{
					client: db,
				}
			},
		},
		{
			name: "ErrorSyncPublisher",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				data:  &messaging.Data{Msg: make([]byte, 0)},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) *Client {
				db, mock := redismock.NewClientMock()

				mock.ExpectPublish(a.topic, a.data.Msg).SetErr(assert.AnError)

				return &Client{
					client:        db,
					syncPublisher: true,
				}
			},
		},
		{
			name: "ErrorAsyncPublisher",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				data:  &messaging.Data{Msg: make([]byte, 0)},
			},
			wantErr: nil,
			mockFn: func(a args) *Client {
				db, mock := redismock.NewClientMock()

				mock.ExpectPublish(a.topic, a.data.Msg).SetErr(assert.AnError)

				return &Client{
					client:        db,
					syncPublisher: false,
					log:           logger.NewNoopLogger(),
				}
			},
		},
		{
			name: "SuccessSyncPublisher",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				data:  &messaging.Data{Msg: make([]byte, 0)},
			},
			wantErr: nil,
			mockFn: func(a args) *Client {
				db, mock := redismock.NewClientMock()

				mock.ExpectPublish(a.topic, a.data.Msg).SetVal(11)

				return &Client{
					client:        db,
					syncPublisher: true,
					log:           logger.NewNoopLogger(),
				}
			},
		},
		{
			name: "SuccessSyncPublisher",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				data:  &messaging.Data{Msg: make([]byte, 0)},
			},
			wantErr: nil,
			mockFn: func(a args) *Client {
				db, mock := redismock.NewClientMock()

				mock.ExpectPublish(a.topic, a.data.Msg).SetVal(12)

				return &Client{
					client:        db,
					syncPublisher: false,
					log:           logger.NewNoopLogger(),
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.mockFn(tt.args).Publish(tt.args.ctx, tt.args.topic, tt.args.data)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestClient_BulkPublish(t *testing.T) {
	type args struct {
		ctx   context.Context
		topic string
		datas []*messaging.Data
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) *Client
	}{
		{
			name: "ErrorInactiveClient",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				datas: nil,
			},
			wantErr: ErrInactiveClient,
			mockFn: func(a args) *Client {
				return &Client{}
			},
		},
		{
			name: "ErrorInvalidData",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				datas: []*messaging.Data{{Msg: nil}},
			},
			wantErr: messaging.ErrMessageNil,
			mockFn: func(a args) *Client {
				db, _ := redismock.NewClientMock()

				return &Client{client: db}
			},
		},
		{
			name: "ErrorSyncPublisher",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				datas: []*messaging.Data{{Msg: make([]byte, 0)}},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) *Client {
				db, mock := redismock.NewClientMock()

				mock.ExpectPublish(a.topic, a.datas[0].Msg).SetErr(assert.AnError)

				return &Client{client: db, syncPublisher: true}
			},
		},
		{
			name: "SuccessSyncPublisher",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				datas: []*messaging.Data{{Msg: make([]byte, 0)}},
			},
			wantErr: nil,
			mockFn: func(a args) *Client {
				db, mock := redismock.NewClientMock()

				mock.ExpectPublish(a.topic, a.datas[0].Msg).SetVal(1)

				return &Client{client: db, syncPublisher: true}
			},
		},
		{
			name: "ErrorAsyncPublisher",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				datas: []*messaging.Data{{Msg: make([]byte, 0)}},
			},
			wantErr: nil,
			mockFn: func(a args) *Client {
				db, mock := redismock.NewClientMock()

				mock.ExpectPublish(a.topic, a.datas[0].Msg).SetErr(assert.AnError)

				return &Client{client: db, syncPublisher: false, log: logger.NewNoopLogger()}
			},
		},
		{
			name: "SuccessAsyncPublisher",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				datas: []*messaging.Data{{Msg: make([]byte, 0)}},
			},
			wantErr: nil,
			mockFn: func(a args) *Client {
				db, mock := redismock.NewClientMock()

				mock.ExpectPublish(a.topic, a.datas[0].Msg).SetVal(10)

				return &Client{client: db, syncPublisher: false}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.mockFn(tt.args).BulkPublish(tt.args.ctx, tt.args.topic, tt.args.datas)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestClient_Subscribe(t *testing.T) {
	type args struct {
		ctx     context.Context
		topic   string
		subID   string
		handler messaging.SubscriberFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) *Client
	}{
		{
			name: "ErrorInactiveClient",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				subID: "subID",
				handler: func(ctx context.Context, data *messaging.Data) error {
					return nil
				},
			},
			wantErr: ErrInactiveClient,
			mockFn: func(a args) *Client {
				return &Client{}
			},
		},
		{
			// since redismock not support command  Subscribe we cannot mock this
			name: "Success",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				subID: "subID",
				handler: func(ctx context.Context, data *messaging.Data) error {
					return nil
				},
			},
			wantErr: nil,
			mockFn: func(a args) *Client {
				db, _ := redismock.NewClientMock()

				return &Client{client: db}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.mockFn(tt.args)

			err := c.Subscribe(tt.args.ctx, tt.args.topic, tt.args.subID, tt.args.handler)
			assert.Equal(t, tt.wantErr, err)

			time.Sleep(time.Millisecond) // for goroutine
		})
	}
}
