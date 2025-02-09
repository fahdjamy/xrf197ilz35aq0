package storage

import (
	"context"
	"encoding"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisStorage struct {
	client *redis.Client
}

func (r RedisStorage) Get(key string, ctx context.Context) (any, error) {
	return r.client.Get(ctx, key).Result()
}

func (r RedisStorage) Set(key string, value any, expiration time.Duration, ctx context.Context) error {
	if value == nil {
		return fmt.Errorf("value is nil")
	}
	_, ok := value.(encoding.BinaryMarshaler)
	if !ok {
		return fmt.Errorf("value is not a BinaryMarshaler")
	}
	return r.client.Set(ctx, key, value, expiration).Err()
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}
