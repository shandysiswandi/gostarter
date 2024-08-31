package outbound

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/shandysiswandi/gostarter/internal/shortly/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	codecMock "github.com/shandysiswandi/gostarter/pkg/codec/mocker"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisShort(t *testing.T) {
	type args struct {
		client *redis.Client
		json   codec.Codec
	}
	tests := []struct {
		name string
		args args
		want *RedisShort
	}{
		{
			name: "Success",
			args: args{},
			want: &RedisShort{keyMap: "shortly"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewRedisShort(tt.args.client, tt.args.json)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRedisShort_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Short
		wantErr error
		mockFn  func(a args) *RedisShort
	}{
		{
			name:    "NotFound",
			args:    args{ctx: context.TODO(), key: "key"},
			want:    nil,
			wantErr: nil,
			mockFn: func(a args) *RedisShort {
				cli, rmock := redismock.NewClientMock()

				rmock.ExpectHGet("shortly", a.key).RedisNil()

				return &RedisShort{
					client: cli,
					keyMap: "shortly",
				}
			},
		},
		{
			name:    "ErrorRedisHGet",
			args:    args{ctx: context.TODO(), key: "key"},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *RedisShort {
				cli, rmock := redismock.NewClientMock()

				rmock.ExpectHGet("shortly", a.key).SetErr(assert.AnError)

				return &RedisShort{
					client: cli,
					keyMap: "shortly",
				}
			},
		},
		{
			name:    "ErrorJSONDecode",
			args:    args{ctx: context.TODO(), key: "key"},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *RedisShort {
				cli, rmock := redismock.NewClientMock()
				mjson := new(codecMock.MockCodec)

				str := `{}`
				v := domain.Short{}

				rmock.ExpectHGet("shortly", a.key).SetVal(str)
				mjson.EXPECT().Decode([]byte(str), &v).Return(assert.AnError)

				return &RedisShort{
					client: cli,
					keyMap: "shortly",
					json:   mjson,
				}
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.TODO(), key: "key"},
			want:    &domain.Short{Key: "key"},
			wantErr: nil,
			mockFn: func(a args) *RedisShort {
				cli, rmock := redismock.NewClientMock()
				mjson := new(codecMock.MockCodec)

				str := `{"key":"key"}`
				ds := domain.Short{}

				rmock.ExpectHGet("shortly", a.key).SetVal(str)
				mjson.EXPECT().Decode([]byte(str), &ds).
					Run(func(data []byte, v interface{}) {
						if val, ok := v.(*domain.Short); ok {
							*val = domain.Short{Key: "key"}
						}
					}).
					Return(nil)

				return &RedisShort{
					client: cli,
					keyMap: "shortly",
					json:   mjson,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn(tt.args).Get(tt.args.ctx, tt.args.key)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRedisShort_Set(t *testing.T) {
	type args struct {
		ctx   context.Context
		value domain.Short
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) *RedisShort
	}{
		{
			name:    "ErrorJSONEncode",
			args:    args{ctx: context.TODO(), value: domain.Short{}},
			wantErr: assert.AnError,
			mockFn: func(a args) *RedisShort {
				cli, _ := redismock.NewClientMock()
				mjson := new(codecMock.MockCodec)

				mjson.EXPECT().Encode(a.value).Return(nil, assert.AnError)

				return &RedisShort{
					client: cli,
					keyMap: "shortly",
					json:   mjson,
				}
			},
		},
		{
			name:    "ErrorRedisHSet",
			args:    args{ctx: context.TODO(), value: domain.Short{}},
			wantErr: assert.AnError,
			mockFn: func(a args) *RedisShort {
				cli, rmock := redismock.NewClientMock()
				mjson := new(codecMock.MockCodec)

				mjson.EXPECT().Encode(a.value).Return([]byte(`{}`), nil)

				rmock.ExpectHSet("shortly", a.value.Key, []byte(`{}`)).SetErr(assert.AnError)

				return &RedisShort{
					client: cli,
					keyMap: "shortly",
					json:   mjson,
				}
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.TODO(), value: domain.Short{Key: "keyval"}},
			wantErr: nil,
			mockFn: func(a args) *RedisShort {
				cli, rmock := redismock.NewClientMock()
				mjson := new(codecMock.MockCodec)

				mjson.EXPECT().Encode(a.value).Return([]byte(`{"key":"keyval"}`), nil)

				rmock.ExpectHSet("shortly", a.value.Key, []byte(`{"key":"keyval"}`)).SetVal(1)

				return &RedisShort{
					client: cli,
					keyMap: "shortly",
					json:   mjson,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.mockFn(tt.args).Set(tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
