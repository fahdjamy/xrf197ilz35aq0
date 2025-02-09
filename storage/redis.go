package storage

import (
	"context"
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
	return r.client.Set(ctx, key, value, expiration).Err()
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}
