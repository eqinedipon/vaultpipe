// Package env provides utilities for managing process environments.
package env

import (
	"strings"
	"unicode"
)

// SanitizeKey returns a valid POSIX environment variable name by replacing
// illegal characters with underscores and uppercasing the result.
// Empty strings are returned unchanged.
func SanitizeKey(key string) string {
	if key == "" {
		return ""
	}
	var b strings.Builder
	for i, r := range key {
		switch {
		case r == '_':
			b.WriteRune('_')
		case unicode.IsLetter(r):
			b.WriteRune(unicode.ToUpper(r))
		case unicode.IsDigit(r) && i > 0:
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}

// SanitizeMap returns a new map with all keys sanitized via SanitizeKey.
// If two keys collide after sanitization the last value (in iteration order)
// wins; callers that need determinism should sort keys beforehand.
func SanitizeMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[SanitizeKey(k)] = v
	}
	return out
}
