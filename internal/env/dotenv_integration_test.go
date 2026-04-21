package env_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpipe/internal/env"
)

// TestDotEnvMergedWithInjector verifies that values loaded from a .env file
// are available through the Injector and that Vault secrets override them.
func TestDotEnvMergedWithInjector(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	content := []byte("BASE_VAR=from_dotenv\nSHARED=dotenv_value\n")
	if err := os.WriteFile(envFile, content, 0o600); err != nil {
		t.Fatal(err)
	}

	dotEnvVars, err := env.LoadDotEnvFile(envFile)
	if err != nil {
		t.Fatalf("LoadDotEnvFile: %v", err)
	}

	// Vault secrets override dotenv values for the same key.
	vaultSecrets := map[string]string{
		"SHARED": "vault_value",
		"VAULT_ONLY": "secret123",
	}

	base := os.Environ()
	for k, v := range dotEnvVars {
		base = append(base, k+"="+v)
	}

	inj := env.NewInjector(base, vaultSecrets)
	environ := inj.Environ()

	m := toEnvMap(environ)

	if m["BASE_VAR"] != "from_dotenv" {
		t.Errorf("BASE_VAR: got %q, want %q", m["BASE_VAR"], "from_dotenv")
	}
	if m["SHARED"] != "vault_value" {
		t.Errorf("SHARED: got %q, want %q", m["SHARED"], "vault_value")
	}
	if m["VAULT_ONLY"] != "secret123" {
		t.Errorf("VAULT_ONLY: got %q, want %q", m["VAULT_ONLY"], "secret123")
	}
}

// TestDotEnvSanitizeKeys verifies that keys loaded from a .env file are
// sanitized before being merged, matching the behaviour of SanitizeMap.
func TestDotEnvSanitizeKeys(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("my-key=hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	raw, err := env.LoadDotEnvFile(envFile)
	if err != nil {
		t.Fatalf("LoadDotEnvFile: %v", err)
	}

	sanitized := env.SanitizeMap(raw)
	if sanitized["MY_KEY"] != "hello" {
		t.Errorf("expected MY_KEY=hello, got %v", sanitized)
	}
}
