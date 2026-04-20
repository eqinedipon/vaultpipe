package env

import "strings"

// FilterOption controls which environment variables are included.
type FilterOption func(key string) bool

// AllowPrefix returns a FilterOption that passes keys with the given prefix.
func AllowPrefix(prefix string) FilterOption {
	return func(key string) bool {
		return strings.HasPrefix(key, prefix)
	}
}

// DenyList returns a FilterOption that blocks the specified keys.
func DenyList(keys ...string) FilterOption {
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[k] = struct{}{}
	}
	return func(key string) bool {
		_, blocked := set[key]
		return !blocked
	}
}

// Filter returns a new Snapshot containing only variables for which ALL
// provided FilterOptions return true.
func (s *Snapshot) Filter(opts ...FilterOption) *Snapshot {
	out := make(map[string]string)
	for k, v := range s.vars {
		pass := true
		for _, fn := range opts {
			if !fn(k) {
				pass = false
				break
			}
		}
		if pass {
			out[k] = v
		}
	}
	return &Snapshot{vars: out}
}

// Merge returns a new Snapshot that combines s with overrides.
// Keys present in overrides take precedence.
func (s *Snapshot) Merge(overrides map[string]string) *Snapshot {
	out := make(map[string]string, len(s.vars)+len(overrides))
	for k, v := range s.vars {
		out[k] = v
	}
	for k, v := range overrides {
		out[k] = v
	}
	return &Snapshot{vars: out}
}
