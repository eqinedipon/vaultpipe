package env_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/env"
)

// TestResolve_WithSanitizedKeys verifies that Resolve works correctly when
// keys have been sanitized (uppercased, dashes replaced) before lookup.
func TestResolve_WithSanitizedKeys(t *testing.T) {
	raw := map[string]string{
		"db-password": "hunter2",
		"api_key":     "abc123",
	}
	sanitized := env.SanitizeMap(raw)

	out, err := env.Resolve(
		[]string{"DB_PASSWORD", "API_KEY"},
		sanitized,
		nil,
		env.Strict,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASSWORD"] != "hunter2" {
		t.Errorf("DB_PASSWORD: got %q", out["DB_PASSWORD"])
	}
	if out["API_KEY"] != "abc123" {
		t.Errorf("API_KEY: got %q", out["API_KEY"])
	}
}

// TestResolve_WithInjector verifies that a resolved map can be fed into an
// Injector and produces the expected environment slice.
func TestResolve_WithInjector(t *testing.T) {
	secrets := map[string]string{"SECRET_TOKEN": "tok-xyz"}
	base := map[string]string{"LOG_LEVEL": "info", "SECRET_TOKEN": "old"}

	resolved, err := env.Resolve(
		[]string{"SECRET_TOKEN", "LOG_LEVEL"},
		secrets,
		base,
	)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}

	inj := env.NewInjector(base, resolved)
	environ := inj.Environ()

	em := make(map[string]string, len(environ))
	for _, kv := range environ {
		for i := 0; i < len(kv); i++ {
			if kv[i] == '=' {
				em[kv[:i]] = kv[i+1:]
				break
			}
		}
	}

	if em["SECRET_TOKEN"] != "tok-xyz" {
		t.Errorf("SECRET_TOKEN: got %q, want tok-xyz", em["SECRET_TOKEN"])
	}
	if em["LOG_LEVEL"] != "info" {
		t.Errorf("LOG_LEVEL: got %q, want info", em["LOG_LEVEL"])
	}
}
