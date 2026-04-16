package vault

import (
	"context"
	"fmt"
	"strings"
)

// SecretMap is a flat map of key -> string value resolved from Vault.
type SecretMap map[string]string

// GetSecrets fetches secrets at the given path and returns a flat SecretMap.
// KV v2 data envelope (data.data) is unwrapped automatically.
func (c *Client) GetSecrets(ctx context.Context, path string) (SecretMap, error) {
	secret, err := c.logical.ReadWithContext(ctx, kvPath(path))
	if err != nil {
		return nil, fmt.Errorf("vault read %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("vault read %q: no data returned", path)
	}

	data := secret.Data
	// Unwrap KV v2 envelope
	if inner, ok := data["data"]; ok {
		if m, ok := inner.(map[string]interface{}); ok {
			data = m
		}
	}

	result := make(SecretMap, len(data))
	for k, v := range data {
		result[k] = coerceString(v)
	}
	return result, nil
}

// GetMultiple fetches secrets from multiple paths and merges them.
// Later paths take precedence over earlier ones on key collision.
func (c *Client) GetMultiple(ctx context.Context, paths []string) (SecretMap, error) {
	merged := make(SecretMap)
	for _, p := range paths {
		sm, err := c.GetSecrets(ctx, p)
		if err != nil {
			return nil, err
		}
		for k, v := range sm {
			merged[k] = v
		}
	}
	return merged, nil
}

// kvPath normalises a secret path for KV v2 mounts by inserting /data/.
func kvPath(path string) string {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 2 {
		return parts[0] + "/data/" + parts[1]
	}
	return path
}

func coerceString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}
