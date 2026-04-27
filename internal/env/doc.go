// Package env provides utilities for constructing, transforming, and injecting
// environment variable maps into child processes.
//
// # Sub-features
//
//   - Injector   – merges secret maps with a base environment and applies them
//     to exec.Cmd instances.
//   - Snapshot   – captures an immutable point-in-time view of an env map and
//     computes diffs between snapshots.
//   - Filter     – allow-list / deny-list filtering by key prefix or name.
//   - Expand     – shell-style variable expansion within values.
//   - Sanitize   – normalises keys to UPPER_SNAKE_CASE.
//   - DotEnv     – parses .env files into maps.
//   - Prefix     – strips or prepends key prefixes.
//   - Merge      – merges multiple maps with configurable overwrite semantics.
//   - Resolve    – resolves keys against secrets, base env, and OS env.
//   - Truncate   – caps values at a configurable maximum length.
//   - Validate   – asserts required keys are present and non-empty.
//   - Coerce     – converts non-string values (bool, int, float) to strings.
//   - Watch      – polls for env changes and invokes a callback on diff.
//   - Transform  – applies user-defined transformation functions to maps.
//   - Chain      – composes multiple env pipeline steps in order.
//   - Defaults   – fills missing keys from a set of DefaultSpec entries
//     without mutating the source map.
package env
