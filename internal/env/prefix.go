// Package env provides utilities for environment variable management.
package env

import "strings"

// PrefixMapper strips or adds a prefix to environment variable keys.
type PrefixMapper struct {
	prefix string
}

// NewPrefixMapper creates a PrefixMapper that operates on the given prefix.
func NewPrefixMapper(prefix string) *PrefixMapper {
	return &PrefixMapper{prefix: prefix}
}

// Strip returns a new map with the prefix removed from all matching keys.
// Keys that do not carry the prefix are dropped.
func (p *PrefixMapper) Strip(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if strings.HasPrefix(k, p.prefix) {
			newKey := strings.TrimPrefix(k, p.prefix)
			if newKey != "" {
				out[newKey] = v
			}
		}
	}
	return out
}

// Add returns a new map with the prefix prepended to every key.
func (p *PrefixMapper) Add(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[p.prefix+k] = v
	}
	return out
}

// FilterByPrefix returns only entries whose keys start with the given prefix.
// Unlike Strip, keys are kept unchanged.
func FilterByPrefix(m map[string]string, prefix string) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		if strings.HasPrefix(k, prefix) {
			out[k] = v
		}
	}
	return out
}
