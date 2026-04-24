// Package env provides utilities for managing process environment variables.
package env

import "fmt"

// TruncateOptions controls how values are truncated.
type TruncateOptions struct {
	// MaxLen is the maximum byte length of a value before truncation.
	// Values at or below this length are returned unchanged.
	MaxLen int
	// Suffix is appended to truncated values to indicate truncation.
	// Defaults to "..." if empty.
	Suffix string
}

// TruncateMap returns a copy of m where any value exceeding opts.MaxLen bytes
// is truncated and the configured suffix is appended. Keys are never modified.
// If opts.MaxLen is zero or negative, m is returned unchanged (no copy).
func TruncateMap(m map[string]string, opts TruncateOptions) map[string]string {
	if opts.MaxLen <= 0 {
		return m
	}
	suffix := opts.Suffix
	if suffix == "" {
		suffix = "..."
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		if len(v) > opts.MaxLen {
			out[k] = v[:opts.MaxLen] + suffix
		} else {
			out[k] = v
		}
	}
	return out
}

// TruncateValue truncates a single string value to maxLen bytes, appending
// suffix. If the value is within maxLen it is returned unchanged.
// An error is returned if maxLen is negative.
func TruncateValue(v string, maxLen int, suffix string) (string, error) {
	if maxLen < 0 {
		return "", fmt.Errorf("env: TruncateValue: maxLen must be >= 0, got %d", maxLen)
	}
	if maxLen == 0 || len(v) <= maxLen {
		return v, nil
	}
	if suffix == "" {
		suffix = "..."
	}
	return v[:maxLen] + suffix, nil
}
