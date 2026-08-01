package types

import (
	"encoding/json"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Optional wraps a value that may be absent from a JSON payload, present with
// a value, or present and explicitly null. Plain pointers cannot tell the
// first case apart from the third, which is what PATCH semantics hinge on.
type Optional[T any] struct {
	Value T
	// Set reports whether the key was present in the payload at all.
	Set bool
	// Null reports whether the key was present and explicitly null.
	Null bool
}

// Some returns an Optional holding a value.
func Some[T any](value T) Optional[T] {
	return Optional[T]{Value: value, Set: true}
}

// Null returns an Optional representing an explicit JSON null.
func Null[T any]() Optional[T] {
	return Optional[T]{Set: true, Null: true}
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.Set = true

	if string(data) == "null" {
		o.Null = true
		return nil
	}

	return json.Unmarshal(data, &o.Value)
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if o.Null {
		return []byte("null"), nil
	}

	return json.Marshal(o.Value)
}

// Ptr returns the value as a pointer, or nil for an explicit null. Only
// meaningful when Set is true.
func (o Optional[T]) Ptr() *T {
	if o.Null {
		return nil
	}

	return &o.Value
}

// WhenSet applies ordinary validation rules to the wrapped value, but only
// when the key was present in the payload. An explicit null presents the zero
// value to the rules, so pair it with NotBlankWhenSet on NOT NULL columns.
func WhenSet[T any](rules ...validation.Rule) validation.Rule {
	return validation.By(func(value any) error {
		opt, ok := value.(Optional[T])
		if !ok || !opt.Set {
			return nil
		}

		return validation.Validate(opt.Value, rules...)
	})
}

// NotBlankWhenSet validates an Optional[string] field backed by a NOT NULL
// column: an absent key is fine, but a key that is present must carry a
// non-blank value. Absent keys are never rejected, so the same rule works for
// both PATCH and PUT payloads.
var NotBlankWhenSet = validation.By(func(value any) error {
	opt, ok := value.(Optional[string])
	if !ok || !opt.Set {
		return nil
	}

	if opt.Null || strings.TrimSpace(opt.Value) == "" {
		return validation.ErrRequired
	}

	return nil
})