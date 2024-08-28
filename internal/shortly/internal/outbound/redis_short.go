package outbound

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/shandysiswandi/gostarter/internal/shortly/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/codec"
)

type RedisShort struct {
	client *redis.Client
	json   codec.Codec
	keyMap string
}

func NewRedisShort(client *redis.Client, json codec.Codec) *RedisShort {
	return &RedisShort{client: client, json: json, keyMap: "shortly"}
}

func (rs *RedisShort) Get(ctx context.Context, key string) (*domain.Short, error) {
	str, err := rs.client.HGet(ctx, rs.keyMap, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil //nolint:nilnil // no rows is not an error, just a nil result
	}

	if err != nil {
		return nil, err
	}

	var short domain.Short
	err = rs.json.Decode([]byte(str), &short)
	if err != nil {
		return nil, err
	}

	return &short, nil
}

func (rs *RedisShort) Set(ctx context.Context, value domain.Short) error {
	bt, err := rs.json.Encode(value)
	if err != nil {
		return err
	}

	return rs.client.HSet(ctx, rs.keyMap, value.Key, bt).Err()
}
