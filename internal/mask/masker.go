// Package mask provides utilities for redacting secret values from strings
// and log output to prevent accidental disclosure.
package mask

import "strings"

// Masker holds a set of secret values and can redact them from strings.
type Masker struct {
	secrets []string
}

// New creates a Masker preloaded with the provided secret values.
// Empty strings are ignored.
func New(secrets []string) *Masker {
	filtered := make([]string, 0, len(secrets))
	for _, s := range secrets {
		if s != "" {
			filtered = append(filtered, s)
		}
	}
	return &Masker{secrets: filtered}
}

// Redact replaces every occurrence of a known secret value within input
// with the literal string "***REDACTED***".
func (m *Masker) Redact(input string) string {
	for _, secret := range m.secrets {
		input = strings.ReplaceAll(input, secret, "***REDACTED***")
	}
	return input
}

// Add appends additional secret values to the masker at runtime.
func (m *Masker) Add(values ...string) {
	for _, v := range values {
		if v != "" {
			m.secrets = append(m.secrets, v)
		}
	}
}

// Len returns the number of tracked secret values.
func (m *Masker) Len() int {
	return len(m.secrets)
}
