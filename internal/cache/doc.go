// Package cache implements a lightweight, thread-safe, TTL-based in-memory
// cache for Vault secret payloads.
//
// The cache is keyed by Vault path and stores the flattened key/value map
// returned by the secrets fetcher. A TTL of zero disables caching entirely,
// which is useful for security-sensitive environments where stale secrets
// must never be used.
//
// The cache is intentionally session-scoped: it is created once per
// vaultpipe invocation and discarded when the process exits, so no secret
// material is persisted between runs.
package cache
