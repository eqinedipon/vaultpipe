package env

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a single environment value.
type TransformFunc func(key, value string) (string, error)

// Transformer applies a chain of TransformFuncs to every key/value pair in a map.
type Transformer struct {
	fns []TransformFunc
}

// NewTransformer creates a Transformer with the provided transform functions
// applied in order.
func NewTransformer(fns ...TransformFunc) *Transformer {
	return &Transformer{fns: fns}
}

// Apply runs every registered TransformFunc over each entry in src and returns
// a new map with the transformed values. The original map is never mutated.
// If any transform returns an error the call is aborted and the error is
// returned with the offending key included in the message.
func (t *Transformer) Apply(src map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		current := v
		for _, fn := range t.fns {
			result, err := fn(k, current)
			if err != nil {
				return nil, fmt.Errorf("transform error on key %q: %w", k, err)
			}
			current = result
		}
		out[k] = current
	}
	return out, nil
}

// TrimSpaceTransform returns a TransformFunc that trims leading and trailing
// whitespace from every value.
func TrimSpaceTransform() TransformFunc {
	return func(_, v string) (string, error) {
		return strings.TrimSpace(v), nil
	}
}

// UpperKeyTransform returns a TransformFunc that upper-cases the value when the
// key has the given suffix (case-insensitive). Useful for flag-style secrets.
func UpperKeyTransform(suffix string) TransformFunc {
	suffix = strings.ToUpper(suffix)
	return func(k, v string) (string, error) {
		if strings.HasSuffix(strings.ToUpper(k), suffix) {
			return strings.ToUpper(v), nil
		}
		return v, nil
	}
}
