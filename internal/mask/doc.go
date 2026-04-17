// Package mask implements secret-value redaction for vaultpipe.
//
// When secrets fetched from Vault are injected into a child process, the same
// values may inadvertently appear in log lines, error messages, or audit
// events. The Masker type tracks every known secret value and replaces any
// occurrence with the placeholder string "***REDACTED***" before the text is
// written to any output stream.
//
// Usage:
//
//	m := mask.New(secretValues)
//	safeMsg := m.Redact(potentiallyLeakyString)
package mask
