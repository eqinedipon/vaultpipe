// Package env provides utilities for environment variable management,
// injection, filtering, and expansion.
package env

import (
	"os"
	"strings"
)

// Expander resolves ${VAR} and $VAR references within secret values
// using a combined lookup of provided secrets and the base environment.
type Expander struct {
	lookup func(string) (string, bool)
}

// NewExpander an Expander that resolves variable references by
// first checking secrets, then falling back to the base environment slice
in KEY=VALUE form).
func NewExpander(secrets map[string]string, base []string) *Expander {
	baseMap := make(map[string]string, len(base))
	for _, kv := range base {
		if idx := strings.IndexByte(kv, '='); idx >= 0 {
			baseMap[kv[:idx]] = kv[idx+1:]
		}
	}

	return &Expander{
		lookup: func(key string) (string, bool) {
			if v, ok := secrets[key]; ok {
				return v, true
			}
			if v, ok := baseMap[key]; ok {
				return v, true
			}
			return os.LookupEnv(key)
		},
	}
}

// Expand resolves variable references in s using the configured lookup.
// Supports both $VAR and ${VAR} syntax. Unknown variables are left as-is.
func (e *Expander) Expand(s string) string {
	return os.Expand(s, func(key string) string {
		if v, ok := e.lookup(key); ok {
			return v
		}
		// Preserve unresolved references rather than silently dropping them.
		if strings.Contains(key, ":") || key == "" {
			return ""
		}
		return "${" + key + "}"
	})
}

// ExpandAll applies Expand to every value in the provided map, returning
// a new map with resolved values. The original map is not mutated.
func (e *Expander) ExpandAll(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = e.Expand(v)
	}
	return out
}
