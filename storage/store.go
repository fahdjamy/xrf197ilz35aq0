package storage

import (
	"context"
	"xrf197ilz35aq0/internal"
)

// Store is a port
type Store interface {
	Save(collection string, obj internal.Serializable, ctx context.Context) (any, error)
	FindById(collection string, id int64, ctx context.Context) (*internal.Serializable, error)
}
