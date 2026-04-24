// Package env provides utilities for constructing and manipulating
// process environments when injecting secrets.
package env

// MergeOption controls how two environment maps are combined.
type MergeOption func(dst, src map[string]string)

// WithOverwrite causes src values to overwrite dst values for the same key.
// This is the default behaviour when no option is supplied.
func WithOverwrite() MergeOption {
	return func(dst, src map[string]string) {
		for k, v := range src {
			dst[k] = v
		}
	}
}

// WithNoOverwrite causes src values to be ignored when the key already
// exists in dst, preserving the original value.
func WithNoOverwrite() MergeOption {
	return func(dst, src map[string]string) {
		for k, v := range src {
			if _, exists := dst[k]; !exists {
				dst[k] = v
			}
		}
	}
}

// Merge combines one or more source maps into a new map, applying each
// MergeOption in order. When no options are provided WithOverwrite is used.
// Neither dst nor any src map is mutated.
func Merge(base map[string]string, sources []map[string]string, opts ...MergeOption) map[string]string {
	if len(opts) == 0 {
		opts = []MergeOption{WithOverwrite()}
	}

	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}

	for _, src := range sources {
		for _, opt := range opts {
			opt(result, src)
		}
	}

	return result
}

// MergeTwo is a convenience wrapper around Merge for the common case of
// combining exactly two maps.
func MergeTwo(base, override map[string]string, opts ...MergeOption) map[string]string {
	return Merge(base, []map[string]string{override}, opts...)
}
