package storage

import "xrf197ilz35aq0"

type Store interface {
	Save(key string, obj xrf197ilz35aq0.Serializable) error
}
