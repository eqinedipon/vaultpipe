// Package env provides utilities for injecting, transforming, and validating
// environment variables sourced from Vault secrets.
package env

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError is returned when one or more required keys fail validation.
type ValidationError struct {
	Missing  []string
	Invalid  []string
}

func (e *ValidationError) Error() string {
	var parts []string
	if len(e.Missing) > 0 {
		parts = append(parts, fmt.Sprintf("missing required keys: %s", strings.Join(e.Missing, ", ")))
	}
	if len(e.Invalid) > 0 {
		parts = append(parts, fmt.Sprintf("invalid values for keys: %s", strings.Join(e.Invalid, ", ")))
	}
	return strings.Join(parts, "; ")
}

// ValidateOption configures the behaviour of Validate.
type ValidateOption func(*validateConfig)

type validateConfig struct {
	required []string
	noEmpty  bool
}

// RequireKeys returns a ValidateOption that asserts each named key is present
// in the environment map.
func RequireKeys(keys ...string) ValidateOption {
	return func(c *validateConfig) {
		c.required = append(c.required, keys...)
	}
}

// NoEmptyValues returns a ValidateOption that treats keys with empty string
// values as invalid.
func NoEmptyValues() ValidateOption {
	return func(c *validateConfig) {
		c.noEmpty = true
	}
}

// Validate checks env against the supplied options and returns a
// *ValidationError if any constraint is violated, or nil on success.
func Validate(env map[string]string, opts ...ValidateOption) error {
	cfg := &validateConfig{}
	for _, o := range opts {
		o(cfg)
	}

	var missing, invalid []string

	for _, key := range cfg.required {
		val, ok := env[key]
		if !ok {
			missing = append(missing, key)
			continue
		}
		if cfg.noEmpty && val == "" {
			invalid = append(invalid, key)
		}
	}

	if cfg.noEmpty {
		for k, v := range env {
			if v == "" {
				// only flag keys not already captured via required check
				alreadyFlagged := false
				for _, inv := range invalid {
					if inv == k {
						alreadyFlagged = true
						break
					}
				}
				if !alreadyFlagged {
					invalid = append(invalid, k)
				}
			}
		}
	}

	if len(missing) == 0 && len(invalid) == 0 {
		return nil
	}
	return &ValidationError{Missing: missing, Invalid: invalid}
}

// IsValidationError reports whether err is a *ValidationError.
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}
