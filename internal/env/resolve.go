// Package env provides utilities for constructing, manipulating, and
// injecting environment variable sets into child processes.
package env

import (
	"fmt"
	"strings"
)

// ResolveOption configures the behaviour of Resolve.
type ResolveOption func(*resolveConfig)

type resolveConfig struct {
	allowEmpty bool
	strict     bool
}

// AllowEmpty permits variables that resolve to an empty string. By default
// Resolve drops keys whose final value is the empty string.
func AllowEmpty() ResolveOption {
	return func(c *resolveConfig) { c.allowEmpty = true }
}

// Strict causes Resolve to return an error if any required key (one whose
// value begins with "required:") cannot be satisfied.
func Strict() ResolveOption {
	return func(c *resolveConfig) { c.strict = true }
}

// Resolve builds a final environment map by layering sources in priority order
// (lowest → highest): base OS environment, dotenv overrides, Vault secrets.
//
// Each value is expanded via the provided Expander so that cross-references
// such as DB_URL=postgres://${DB_USER}:${DB_PASS}@host/db are resolved.
//
// When Strict is set, any value that still starts with the sentinel prefix
// "required:" after expansion is treated as an unresolved required variable
// and causes an error.
func Resolve(
	base map[string]string,
	dotenv map[string]string,
	secrets map[string]string,
	expander *Expander,
	opts ...ResolveOption,
) (map[string]string, error) {
	cfg := &resolveConfig{}
	for _, o := range opts {
		o(cfg)
	}

	// Layer the sources: base < dotenv < secrets.
	merged := Merge(
		[]map[string]string{base, dotenv, secrets},
		WithOverwrite(),
	)

	resolved := make(map[string]string, len(merged))
	var errs []string

	for k, v := range merged {
		expanded := expander.Expand(v)

		// Strict mode: detect unresolved required placeholders.
		if cfg.strict && strings.HasPrefix(expanded, "required:") {
			errs = append(errs, fmt.Sprintf(
				"required variable %q is not set (value: %q)", k, expanded,
			))
			continue
		}

		// Drop empty values unless AllowEmpty is set.
		if expanded == "" && !cfg.allowEmpty {
			continue
		}

		resolved[k] = expanded
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("resolve: unresolved required variables:\n  %s",
			strings.Join(errs, "\n  "))
	}

	return resolved, nil
}
