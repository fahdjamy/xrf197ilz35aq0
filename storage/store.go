package storage

import (
	"context"
	"xrf197ilz35aq0"
)

// Store is a port
type Store interface {
	SetContext(ctx context.Context)
	Save(collection string, obj xrf197ilz35aq0.Serializable) (any, error)
	FindById(collection string, id int64) (*xrf197ilz35aq0.Serializable, error)
}
