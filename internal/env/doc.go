// Package env provides utilities for constructing and manipulating process
// environments when streaming secrets from Vault.
//
// Key components:
//
//   - Injector   – merges base environment variables with secret key/value pairs.
//   - Snapshot   – captures an immutable view of an environment for diffing.
//   - Filter     – allow/deny rules applied to environment variable sets.
//   - Expander   – resolves ${VAR} references against secrets and the base env.
//   - SanitizeKey / SanitizeMap – convert arbitrary Vault secret keys into valid
//     POSIX environment variable names (uppercase, no illegal characters).
package env
