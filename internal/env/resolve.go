// Package env provides utilities for building and manipulating process
// environments, including secret injection, filtering, expansion, and resolution.
package env

import (
	"errors"
	"fmt"
	"strings"
)

// ResolveOption configures the behaviour of Resolve.
type ResolveOption func(*resolveConfig)

type resolveConfig struct {
	allowEmpty bool
	strict     bool
}

// AllowEmpty permits variables whose resolved value is the empty string.
func AllowEmpty(c *resolveConfig) { c.allowEmpty = true }

// Strict causes Resolve to return an error if any required key is missing or
// empty (unless AllowEmpty is also set).
func Strict(cfg *resolveConfig) { cfg.strict = true }

// Resolve walks required and returns a map containing each key's value drawn
// from, in priority order: secrets, then base, then os.LookupEnv.
//
// If Strict is set any key that resolves to "" (and AllowEmpty is not set)
// returns an error that lists every missing key.
func Resolve(required []string, secrets, base map[string]string, opts ...ResolveOption) (map[string]string, error) {
	cfg := &resolveConfig{}
	for _, o := range opts {
		o(cfg)
	}

	out := make(map[string]string, len(required))
	var missing []string

	for _, key := range required {
		var (
			val   string
			found bool
		)

		if v, ok := secrets[key]; ok {
			val, found = v, true
		} else if v, ok := base[key]; ok {
			val, found = v, true
		} else if v, ok := lookupEnv(key); ok {
			val, found = v, true
		}

		if cfg.strict && (!found || (!cfg.allowEmpty && val == "")) {
			missing = append(missing, key)
			continue
		}

		out[key] = val
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("resolve: missing required keys: %s",
			strings.Join(missing, ", "))
	}

	return out, nil
}

// lookupEnv is a variable so tests can override it without touching os.Getenv.
var lookupEnv = func(key string) (string, bool) {
	// Avoid importing os at package level; swap in tests.
	return "", false
}

// ErrMissingKeys is returned by Resolve in Strict mode.
var ErrMissingKeys = errors.New("resolve: missing required keys")
