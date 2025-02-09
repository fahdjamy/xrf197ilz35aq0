package storage

import (
	"context"
	"time"
)

type Cache interface {
	Get(key string, ctx context.Context) (any, error)
	Set(key string, value any, expiration time.Duration, ctx context.Context) error
}
