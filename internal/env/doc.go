// Package env provides utilities for constructing, transforming, validating,
// and injecting environment variable maps into child processes.
//
// # Core types
//
//   - Injector: merges a base environment with secret overrides and exposes
//     the result as a []string slice compatible with os/exec.
//   - Snapshot: an immutable point-in-time copy of an environment map that
//     supports diffing to detect added or changed keys.
//   - Expander: resolves ${VAR} references inside values against a combined
//     secrets + base environment lookup.
//   - Watcher: polls a secret source on a configurable interval and calls an
//     onChange callback when the snapshot differs from the previous fetch.
//   - Transformer: applies a chain of per-entry functions (e.g. TrimSpace,
//     UpperKey) to produce a transformed copy of a map.
//   - Chain: composes multiple ChainStep functions — including Transformer,
//     Validate, Sanitize, Truncate, and Coerce passes — into a single
//     ordered pipeline.
//
// # Helpers
//
// SanitizeKey / SanitizeMap, CoerceMap, TruncateMap, Merge / MergeTwo,
// ParseDotEnv / LoadDotEnvFile, Resolve, NewPrefixMapper, and Validate
// are standalone functions that can be used independently or composed via
// Chain.
package env
