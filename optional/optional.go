package optional

import (
	"encoding/json"
)

// Optional represents a value that may be explicitly set or left absent.
// When used with `omitzero`, absent values are omitted from JSON.
type Optional[T any] struct {
	Set   bool
	Value T
}

// Some creates an Optional with a value set.
func Some[T any](v T) Optional[T] {
	return Optional[T]{Set: true, Value: v}
}

// None creates an Optional that is explicitly absent.
func None[T any]() Optional[T] {
	return Optional[T]{Set: false}
}

// IsZero allows `omitzero` to omit absent values.
func (o Optional[T]) IsZero() bool {
	return !o.Set
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.Set {
		return []byte("null"), nil
	}
	return json.Marshal(o.Value)
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.Set = true
	return json.Unmarshal(data, &o.Value)
}
