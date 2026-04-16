// Package vault provides a thin client around the HashiCorp Vault API
// tailored for vaultpipe's use-case: reading KV secrets and returning
// them as flat string maps suitable for process environment injection.
//
// Supported features:
//   - KV v2 path normalisation and envelope unwrapping
//   - Merging secrets from multiple paths with last-write-wins semantics
//   - Token and environment-variable based authentication
//
// Usage:
//
//	client, err := vault.NewClient(vault.Config{
//		Address: "https://vault.example.com",
//		Token:   os.Getenv("VAULT_TOKEN"),
//	})
//	secrets, err := client.GetMultiple(ctx, []string{"secret/base", "secret/myapp"})
package vault
