// Package redact provides an io.Writer decorator that transparently
// redacts Vault secret values from any output stream.
//
// Usage:
//
//	m := mask.New(secretValues)
//	safeStdout := redact.New(os.Stdout, m)
//
// Any bytes written to safeStdout that contain a registered secret
// are replaced with "[REDACTED]" before reaching the underlying
// writer. This prevents accidental leakage of secrets through
// process stdout or stderr.
package redact
