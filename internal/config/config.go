package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// SecretMapping defines a single Vault secret path and which keys to pull.
type SecretMapping struct {
	Path string            `yaml:"path"`
	Keys map[string]string `yaml:"keys"` // vault key -> env var name
}

// Config holds the top-level vaultpipe configuration.
type Config struct {
	VaultAddr  string          `yaml:"vault_addr"`
	VaultToken string          `yaml:"vault_token"`
	Mount      string          `yaml:"mount"`
	Secrets    []SecretMapping `yaml:"secrets"`
}

// Load reads a YAML config file from the given path and returns a Config.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	// Allow env vars to override token and addr.
	if v := os.Getenv("VAULT_ADDR"); v != "" && cfg.VaultAddr == "" {
		cfg.VaultAddr = v
	}
	if v := os.Getenv("VAULT_TOKEN"); v != "" && cfg.VaultToken == "" {
		cfg.VaultToken = v
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Secrets) == 0 {
		return fmt.Errorf("at least one secret mapping is required")
	}
	for i, s := range c.Secrets {
		if s.Path == "" {
			return fmt.Errorf("secret[%d]: path is required", i)
		}
		if len(s.Keys) == 0 {
			return fmt.Errorf("secret[%d]: at least one key mapping is required", i)
		}
	}
	return nil
}
