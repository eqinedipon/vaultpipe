// Package env provides utilities for constructing, manipulating, and injecting
// process environment variables when streaming secrets from Vault.
//
// # Sub-features
//
//   - Injector   – merges a base environment with secret overrides and exposes
//     the result as a []string suitable for exec.Cmd.Env.
//   - Snapshot   – captures an environment map at a point in time and can diff
//     it against a later state to detect leakage or drift.
//   - Filter     – allow/deny rules for trimming the environment before
//     passing it to child processes.
//   - Expand     – shell-style variable expansion that resolves ${VAR}
//     references against secrets, base env, and os.Environ.
//   - Sanitize   – normalises arbitrary key strings to valid POSIX env-var
//     names (uppercase, replace illegal chars).
//   - DotEnv     – parses .env files and merges them into the environment.
//   - Prefix     – strips or adds a prefix to all keys in a map, and filters
//     by prefix.
//   - Merge      – combines multiple environment maps with configurable
//     overwrite semantics.
//   - Resolve    – resolves a required key list from secrets, base env, and
//     the OS environment, with optional strict mode.
package env
