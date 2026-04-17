// Package audit provides structured audit logging for vaultpipe operations.
//
// Events are written as newline-delimited JSON (NDJSON) to a configurable
// io.Writer (defaulting to stderr). Each event captures the UTC timestamp,
// the action performed (fetch, exec, error), the relevant path, and — for
// fetch events — the list of secret keys that were retrieved (values are
// never logged).
//
// Usage:
//
//	l := audit.NewLogger(os.Stderr)
//	l.LogFetch("secret/data/myapp", []string{"DB_PASSWORD"})
//	l.LogExec("/usr/local/bin/myapp")
package audit
