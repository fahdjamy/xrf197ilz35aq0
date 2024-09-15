package custom

import (
	"encoding/json"
)

type Field interface {
	string | float64 | complex64 | int | int8 | int16 | int32 | int64
}

type Secret[T Field] struct {
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

func NewSecret[T Field](data T) *Secret[T] {
	return &Secret[T]{data: data}
}
