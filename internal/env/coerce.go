// Package env provides utilities for environment variable management.
package env

import (
	"fmt"
	"strconv"
)

// CoerceOption configures coercion behaviour.
type CoerceOption func(*coerceConfig)

type coerceConfig struct {
	boolTruthy  []string
	boolFalsy   []string
	stringifyAll bool
}

func defaultCoerceConfig() *coerceConfig {
	return &coerceConfig{
		boolTruthy:  []string{"true", "1", "yes", "on"},
		boolFalsy:   []string{"false", "0", "no", "off"},
		stringifyAll: false,
	}
}

// StringifyAll forces all values to string without type-specific formatting.
func StringifyAll() CoerceOption {
	return func(c *coerceConfig) { c.stringifyAll = true }
}

// CoerceValue converts an arbitrary value to a string suitable for use as an
// environment variable value. Booleans, integers, floats and strings are
// handled explicitly; everything else falls back to fmt.Sprintf.
func CoerceValue(v any, opts ...CoerceOption) (string, error) {
	cfg := defaultCoerceConfig()
	for _, o := range opts {
		o(cfg)
	}
	if cfg.stringifyAll {
		return fmt.Sprintf("%v", v), nil
	}
	switch val := v.(type) {
	case string:
		return val, nil
	case bool:
		if val {
			return "true", nil
		}
		return "false", nil
	case int:
		return strconv.Itoa(val), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32), nil
	case nil:
		return "", nil
	default:
		return fmt.Sprintf("%v", val), nil
	}
}

// CoerceMap converts a map of arbitrary values to a map of strings, returning
// the first coercion error encountered.
func CoerceMap(m map[string]any, opts ...CoerceOption) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		s, err := CoerceValue(v, opts...)
		if err != nil {
			return nil, fmt.Errorf("coerce key %q: %w", k, err)
		}
		out[k] = s
	}
	return out, nil
}
