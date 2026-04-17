package vault

import (
	"fmt"

	"github.com/yourusername/vaultpipe/internal/cache"
)

// CachedClient wraps a Client with an in-memory secret cache.
type CachedClient struct {
	client *Client
	cache  *cache.Cache
}

// NewCachedClient creates a CachedClient using the provided Client and Cache.
func NewCachedClient(c *Client, ca *cache.Cache) *CachedClient {
	return &CachedClient{client: c, cache: ca}
}

// GetSecrets returns secrets for path, serving from cache when available.
func (cc *CachedClient) GetSecrets(mount, path string) (map[string]string, error) {
	key := cacheKey(mount, path)
	if vals, ok := cc.cache.Get(key); ok {
		return vals, nil
	}
	vals, err := cc.client.GetSecrets(mount, path)
	if err != nil {
		return nil, err
	}
	cc.cache.Set(key, vals)
	return vals, nil
}

// InvalidatePath removes a single path from the cache.
func (cc *CachedClient) InvalidatePath(mount, path string) {
	cc.cache.Invalidate(cacheKey(mount, path))
}

func cacheKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}
