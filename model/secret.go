package model

import (
	"encoding/json"
	xrf "xrf197ilz35aq0"
)

type Secret[T xrf.Serializable] struct {
	data T
}

func (s Secret[T]) String() string {
	return "[REDACTED]"
}

func (s Secret[T]) Data() T {
	return s.data
}

func (s Secret[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Data)
}

func (s Secret[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.data)
}

func NewSecret[T xrf.Serializable](data T) *Secret[T] {
	return &Secret[T]{data: data}
}
