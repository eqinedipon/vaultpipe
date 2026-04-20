// Package env handles environment variable injection and snapshotting.
package env

import (
	"os"
	"strings"
)

// Snapshot captures the current process environment as a key-value map.
// It can be used to restore or diff the environment at a later point.
type Snapshot struct {
	vars map[string]string
}

// TakeSnapshot captures the current environment from os.Environ.
func TakeSnapshot() *Snapshot {
	raw := os.Environ()
	vars := make(map[string]string, len(raw))
	for _, entry := range raw {
		key, value, _ := strings.Cut(entry, "=")
		vars[key] = value
	}
	return &Snapshot{vars: vars}
}

// NewSnapshotFromMap creates a Snapshot from an explicit map, useful in tests.
func NewSnapshotFromMap(m map[string]string) *Snapshot {
	copy := make(map[string]string, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return &Snapshot{vars: copy}
}

// Get returns the value for a key and whether it was present.
func (s *Snapshot) Get(key string) (string, bool) {
	v, ok := s.vars[key]
	return v, ok
}

// Keys returns all variable names in the snapshot.
func (s *Snapshot) Keys() []string {
	keys := make([]string, 0, len(s.vars))
	for k := range s.vars {
		keys = append(keys, k)
	}
	return keys
}

// Diff returns keys that differ (added, removed, or changed) between s and other.
func (s *Snapshot) Diff(other *Snapshot) map[string]string {
	diff := make(map[string]string)
	for k, v := range other.vars {
		if sv, ok := s.vars[k]; !ok || sv != v {
			diff[k] = v
		}
	}
	return diff
}

// Environ returns the snapshot as a slice of KEY=VALUE strings.
func (s *Snapshot) Environ() []string {
	out := make([]string, 0, len(s.vars))
	for k, v := range s.vars {
		out = append(out, k+"="+v)
	}
	return out
}
