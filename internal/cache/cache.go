// Package cache provides an in-memory, TTL-based cache for Vault secrets
// to reduce redundant API calls within a single vaultpipe session.
package cache

import (
	"sync"
	"time"
)

// entry holds a cached secret map and its expiry time.
type entry struct {
	values    map[string]string
	expiresAt time.Time
}

// Cache is a thread-safe, TTL-based in-memory secret cache.
type Cache struct {
	mu    sync.RWMutex
	items map[string]entry
	ttl   time.Duration
}

// New creates a Cache with the given TTL. A zero TTL disables caching.
func New(ttl time.Duration) *Cache {
	return &Cache{
		items: make(map[string]entry),
		ttl:   ttl,
	}
}

// Get returns the cached values for key, and whether the entry was found and valid.
func (c *Cache) Get(key string) (map[string]string, bool) {
	if c.ttl == 0 {
		return nil, false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.items[key]
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	// Return a shallow copy to prevent mutation of cached data.
	copy := make(map[string]string, len(e.values))
	for k, v := range e.values {
		copy[k] = v
	}
	return copy, true
}

// Set stores values under key for the cache TTL.
func (c *Cache) Set(key string, values map[string]string) {
	if c.ttl == 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	copy := make(map[string]string, len(values))
	for k, v := range values {
		copy[k] = v
	}
	c.items[key] = entry{values: copy, expiresAt: time.Now().Add(c.ttl)}
}

// Invalidate removes a single key from the cache.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Flush clears all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]entry)
}
