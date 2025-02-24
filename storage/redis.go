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

func (r RedisStorage) Delete(key string, ctx context.Context) (int64, error) {
	// UNLINK key [key ...]: Similar to DEL, but performs the actual memory de-allocation in a separate thread.
	//This is crucial for large keys to avoid blocking the Redis server. Best practice for deleting large keys.
	// H-DEL key field [field ...]: Specifically for Hashes. Deletes one or more fields within a hash
	// L-REM key count value: Specifically for Lists. Removes elements from a list.
	// SREM key member [member ...]: Specifically for Sets. Removes one or more members from a set
	// Z-REM key member [member ...]: Specifically for Sorted Sets. Removes one or more members from a sorted set
	//delCount, err := r.client.Unlink(ctx, key).Result() // UNLINK example (better for *very* large values)
	delCount, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return delCount, nil
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}
