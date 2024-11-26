package storage

import (
	"context"
	"xrf197ilz35aq0/internal"
)

// Store is a port
type Store interface {
	SetContext(ctx context.Context)
	Save(collection string, obj internal.Serializable) (any, error)
	FindById(collection string, id int64) (*internal.Serializable, error)
}
