// Package env provides utilities for constructing, transforming, and injecting
// environment variables into child processes.
//
// # Core concepts
//
// Injector merges secret maps with a base environment and exposes the result as
// a []string suitable for exec.Cmd.Env.
//
// Snapshot captures the current environment so that callers can detect drift
// after secrets are refreshed.
//
// Filter, Merge, Resolve, Expand, Sanitize, Coerce, Truncate, Validate,
// Prefix, DotEnv, and Watch are composable building blocks that each address a
// single concern.
//
// Transformer applies an ordered chain of TransformFuncs to a secret map,
// enabling lightweight value transformations (trim, case conversion, etc.)
// before secrets are injected into the process environment.
package env
