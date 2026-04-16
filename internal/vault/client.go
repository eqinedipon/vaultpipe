package vault

import (
	"context"
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods for secret retrieval.
type Client struct {
	api *vaultapi.Client
}

// Config holds configuration for connecting to Vault.
type Config struct {
	Address string
	Token   string
	RoleID  string
	SecretID string
}

// NewClient creates a new Vault client from the provided config.
// Falls back to environment variables (VAULT_ADDR, VAULT_TOKEN) if fields are empty.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()

	addr := cfg.Address
	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
	}
	if addr != "" {
		vcfg.Address = addr
	}

	api, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("vault: failed to create api client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token != "" {
		api.SetToken(token)
	}

	return &Client{api: api}, nil
}

// GetSecrets reads a KV v2 secret at the given mount and path,
// returning a map of key→value strings.
func (c *Client) GetSecrets(ctx context.Context, mount, path string) (map[string]string, error) {
	secret, err := c.api.KVv2(mount).Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("vault: failed to read secret %s/%s: %w", mount, path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("vault: no data found at %s/%s", mount, path)
	}

	result := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		str, ok := v.(string)
		if !ok {
			str = fmt.Sprintf("%v", v)
		}
		result[k] = str
	}
	return result, nil
}
