// Package env provides utilities for managing environment variables.
package env

import (
	"strings"
	"sync"
)

// RedactedEnv wraps a map of environment variables and provides a redacted
// view suitable for logging — secret values are replaced with a placeholder.
type RedactedEnv struct {
	mu       sync.RWMutex
	env      map[string]string
	secrets  map[string]struct{}
	placeholder string
}

const defaultPlaceholder = "[REDACTED]"

// NewRedactedEnv creates a RedactedEnv from the given env map. Keys listed in
// secretKeys will have their values replaced with the placeholder in any
// exported view.
func NewRedactedEnv(env map[string]string, secretKeys []string) *RedactedEnv {
	secrets := make(map[string]struct{}, len(secretKeys))
	for _, k := range secretKeys {
		secrets[strings.ToUpper(k)] = struct{}{}
	}
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return &RedactedEnv{
		env:         copy,
		secrets:     secrets,
		placeholder: defaultPlaceholder,
	}
}

// WithPlaceholder returns a new RedactedEnv with a custom placeholder string.
func (r *RedactedEnv) WithPlaceholder(p string) *RedactedEnv {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.placeholder = p
	return r
}

// Safe returns a copy of the environment map with secret values redacted.
func (r *RedactedEnv) Safe() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]string, len(r.env))
	for k, v := range r.env {
		if _, secret := r.secrets[strings.ToUpper(k)]; secret {
			out[k] = r.placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// Raw returns a copy of the environment map with all values intact.
// Callers must ensure this is not passed to any logging sink.
func (r *RedactedEnv) Raw() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]string, len(r.env))
	for k, v := range r.env {
		out[k] = v
	}
	return out
}

// IsSecret reports whether the given key is treated as a secret.
func (r *RedactedEnv) IsSecret(key string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.secrets[strings.ToUpper(key)]
	return ok
}
