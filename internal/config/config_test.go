package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpipe/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultpipe-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, `
vault_addr: http://127.0.0.1:8200
vault_token: root
mount: secret
secrets:
  - path: myapp/db
    keys:
      password: DB_PASSWORD
      user: DB_USER
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("vault_addr mismatch: %q", cfg.VaultAddr)
	}
	if len(cfg.Secrets) != 1 {
		t.Fatalf("expected 1 secret mapping, got %d", len(cfg.Secrets))
	}
	if cfg.Secrets[0].Keys["password"] != "DB_PASSWORD" {
		t.Errorf("key mapping mismatch")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_NoSecrets(t *testing.T) {
	path := writeTemp(t, `vault_addr: http://127.0.0.1:8200\nsecrets: []\n`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty secrets")
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://env-host:8200")
	t.Setenv("VAULT_TOKEN", "env-token")
	path := writeTemp(t, `
secrets:
  - path: app/cfg
    keys:
      key: MY_KEY
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "http://env-host:8200" {
		t.Errorf("expected env override for vault_addr, got %q", cfg.VaultAddr)
	}
	if cfg.VaultToken != "env-token" {
		t.Errorf("expected env override for vault_token, got %q", cfg.VaultToken)
	}
}
