// Package env provides a suite of utilities for constructing, transforming,
// and injecting environment variables into child processes.
//
// # Core capabilities
//
//   - Injector: merges secret maps over a base environment slice.
//   - Snapshot / Diff: capture and compare environment state.
//   - Filter: allow/deny lists for environment key selection.
//   - Expand: variable interpolation using secret and base values.
//   - Sanitize: normalise arbitrary keys to valid POSIX variable names.
//   - DotEnv: parse .env files into string maps.
//   - Prefix: add or strip key prefixes.
//   - Merge: combine multiple string maps with configurable overwrite rules.
//   - Resolve: ordered lookup across secrets, base env, and OS environment.
//   - Truncate: cap long values to a configurable maximum length.
//   - Validate: enforce required keys and non-empty value constraints.
//   - Coerce: convert arbitrary typed values (bool, int, float, …) to strings
//     suitable for use as environment variable values.
//
// All functions are safe for concurrent use unless noted otherwise.
package env
