package googlepubsub

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/shandysiswandi/gostarter/pkg/messaging"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func mockServerAndClient() (*pstest.Server, *pubsub.Client, *grpc.ClientConn) {
	srv := pstest.NewServer()
	conn, _ := grpc.NewClient(srv.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client, _ := pubsub.NewClient(context.Background(), "test", option.WithGRPCConn(conn))
	return srv, client, conn
}

func TestNewClient(t *testing.T) {
	os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")
	type args struct {
		ctx       context.Context
		projectID string
		opts      []Option
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr error
	}{
		{
			name: "Error",
			args: args{
				ctx:       context.Background(),
				projectID: "",
				opts: []Option{
					WithAutoAck(),
					WithSyncPublisher(),
					WithLogger(logger.NewNoopLogger()),
				},
			},
			want:    nil,
			wantErr: pubsub.ErrEmptyProjectID,
		},
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				projectID: "-",
				opts: []Option{
					WithAutoAck(),
					WithSyncPublisher(),
					WithLogger(logger.NewNoopLogger()),
				},
			},
			want:    &Client{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewClient(tt.args.ctx, tt.args.projectID, tt.args.opts...)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
		mockFn  func() (*Client, func())
	}{
		{
			name:    "InactiveClient",
			wantErr: nil,
			mockFn: func() (*Client, func()) {
				return &Client{}, nil
			},
		},
		{
			name:    "Success",
			wantErr: nil,
			mockFn: func() (*Client, func()) {
				srv, client, conn := mockServerAndClient()

				return &Client{client: client, cancels: []context.CancelFunc{func() {}}}, func() {
					client.Close()
					conn.Close()
					srv.Close()
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, closer := tt.mockFn()
			if closer != nil {
				defer closer()
			}

			err := c.Close()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestClient_Publish(t *testing.T) {
	type args struct {
		ctx   func() context.Context
		topic string
		data  *messaging.Data
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
		mockFn  func(a args) (*Client, func())
	}{
		{
			name: "InactiveClient",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				topic: "topic",
				data:  &messaging.Data{},
			},
			wantErr: ErrInactiveClient.Error(),
			mockFn: func(a args) (*Client, func()) {
				return &Client{}, nil
			},
		},
		{
			name: "InvalidateData",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				topic: "topic",
				data:  &messaging.Data{},
			},
			wantErr: messaging.ErrMessageNil.Error(),
			mockFn: func(a args) (*Client, func()) {
				srv, client, conn := mockServerAndClient()

				return &Client{client: client, cancels: []context.CancelFunc{func() {}}}, func() {
					client.Close()
					conn.Close()
					srv.Close()
				}
			},
		},
		{
			name: "ErrorTopicNotFound",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				topic: "topic",
				data:  &messaging.Data{Msg: []byte{}},
			},
			wantErr: ErrTopicNotFound.Error(),
			mockFn: func(a args) (*Client, func()) {
				srv, client, conn := mockServerAndClient()

				return &Client{client: client}, func() {
					client.Close()
					conn.Close()
					srv.Close()
				}
			},
		},
		{
			name: "ErrorCheckTopicExists",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()

					return ctx
				},
				topic: "topic",
				data:  &messaging.Data{Msg: []byte{}},
			},
			wantErr: "rpc error: code = Canceled desc = context canceled",
			mockFn: func(a args) (*Client, func()) {
				srv, client, conn := mockServerAndClient()

				return &Client{client: client}, func() {
					client.Close()
					conn.Close()
					srv.Close()
				}
			},
		},
		{
			name: "SuccessSyncPublish",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				topic: "topic",
				data:  &messaging.Data{Msg: []byte{}},
			},
			wantErr: "",
			mockFn: func(a args) (*Client, func()) {
				srv, client, conn := mockServerAndClient()

				client.CreateTopic(a.ctx(), a.topic)

				return &Client{client: client, syncPublisher: true}, func() {
					client.Close()
					conn.Close()
					srv.Close()
				}
			},
		},
		{
			name: "SuccessAsyncPublish",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				topic: "topic",
				data:  &messaging.Data{Msg: []byte{}},
			},
			wantErr: "",
			mockFn: func(a args) (*Client, func()) {
				srv, client, conn := mockServerAndClient()

				client.CreateTopic(a.ctx(), a.topic)

				return &Client{client: client}, func() {
					client.Close()
					conn.Close()
					srv.Close()
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, closer := tt.mockFn(tt.args)
			if closer != nil {
				defer closer()
			}

			err := c.Publish(tt.args.ctx(), tt.args.topic, tt.args.data)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
			}
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
		wantErr string
		mockFn  func(a args) (*Client, func())
	}{
		{
			name: "InactiveClient",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				datas: []*messaging.Data{},
			},
			wantErr: ErrInactiveClient.Error(),
			mockFn: func(a args) (*Client, func()) {
				return &Client{}, nil
			},
		},
		{
			name: "ErrorTopicNotFound",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				datas: []*messaging.Data{},
			},
			wantErr: ErrTopicNotFound.Error(),
			mockFn: func(a args) (*Client, func()) {
				srv, client, conn := mockServerAndClient()

				return &Client{client: client}, func() {
					client.Close()
					conn.Close()
					srv.Close()
				}
			},
		},
		{
			name: "ErrorValidateData",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				datas: []*messaging.Data{{
					Msg:        nil,
					Attributes: map[string]string{},
				}},
			},
			wantErr: messaging.ErrMessageNil.Error(),
			mockFn: func(a args) (*Client, func()) {
				srv, client, conn := mockServerAndClient()

				client.CreateTopic(a.ctx, a.topic)

				return &Client{client: client, syncPublisher: true}, func() {
					client.Close()
					conn.Close()
					srv.Close()
				}
			},
		},
		{
			name: "SuccessSyncBulkPublish",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				datas: []*messaging.Data{{
					Msg:        []byte{},
					Attributes: map[string]string{},
				}},
			},
			wantErr: "",
			mockFn: func(a args) (*Client, func()) {
				srv, client, conn := mockServerAndClient()

				client.CreateTopic(a.ctx, a.topic)

				return &Client{client: client, syncPublisher: true}, func() {
					client.Close()
					conn.Close()
					srv.Close()
				}
			},
		},
		{
			name: "SuccessAsyncBulkPublish",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				datas: []*messaging.Data{{
					Msg:        []byte{},
					Attributes: map[string]string{},
				}},
			},
			wantErr: "",
			mockFn: func(a args) (*Client, func()) {
				srv, client, conn := mockServerAndClient()

				client.CreateTopic(a.ctx, a.topic)

				return &Client{client: client}, func() {
					client.Close()
					conn.Close()
					srv.Close()
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, closer := tt.mockFn(tt.args)
			if closer != nil {
				defer closer()
			}

			err := c.BulkPublish(tt.args.ctx, tt.args.topic, tt.args.datas)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}

func TestClient_Subscribe(t *testing.T) {
	type args struct {
		ctx   context.Context
		topic string
		subID string
		h     messaging.SubscriberFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
		mockFn  func(a args) (*Client, func())
	}{
		{
			name: "InactiveClient",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				subID: "sub",
				h: func(ctx context.Context, data *messaging.Data) error {
					return nil
				},
			},
			wantErr: ErrInactiveClient.Error(),
			mockFn: func(a args) (*Client, func()) {
				return &Client{}, nil
			},
		},
		{
			name: "ErrorTopicNotFound",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				subID: "sub",
				h: func(ctx context.Context, data *messaging.Data) error {
					return nil
				},
			},
			wantErr: ErrTopicNotFound.Error(),
			mockFn: func(a args) (*Client, func()) {
				srv, client, conn := mockServerAndClient()

				return &Client{client: client}, func() {
					client.Close()
					conn.Close()
					srv.Close()
				}
			},
		},
		{
			name: "ErrorSubscriptionNotFound",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				subID: "sub",
				h: func(ctx context.Context, data *messaging.Data) error {
					return nil
				},
			},
			wantErr: ErrSubscriptionNotFound.Error(),
			mockFn: func(a args) (*Client, func()) {
				srv, client, conn := mockServerAndClient()

				client.CreateTopic(a.ctx, a.topic)

				return &Client{client: client}, func() {
					client.Close()
					conn.Close()
					srv.Close()
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx:   context.Background(),
				topic: "topic",
				subID: "sub",
				h: func(ctx context.Context, data *messaging.Data) error {
					return nil
				},
			},
			wantErr: "",
			mockFn: func(a args) (*Client, func()) {
				srv, client, conn := mockServerAndClient()

				topic, _ := client.CreateTopic(a.ctx, a.topic)
				client.CreateSubscription(a.ctx, a.subID, pubsub.SubscriptionConfig{
					Topic: topic,
				})

				return &Client{client: client, log: logger.NewNoopLogger()}, func() {
					client.Close()
					conn.Close()
					srv.Close()
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, closer := tt.mockFn(tt.args)
			if closer != nil {
				defer closer()
			}

			err := c.Subscribe(tt.args.ctx, tt.args.topic, tt.args.subID, tt.args.h)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
			}

			time.Sleep(time.Millisecond)
		})
	}
}
