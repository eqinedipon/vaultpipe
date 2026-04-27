package env_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/env"
)

// TestRedactedEnv_WithSanitizedKeys verifies that keys sanitized via
// SanitizeMap are correctly redacted when referenced by their sanitized names.
func TestRedactedEnv_WithSanitizedKeys(t *testing.T) {
	raw := map[string]string{
		"db-password": "hunter2",
		"app-host":    "localhost",
	}
	sanitized := env.SanitizeMap(raw)

	// Secret keys are referenced by their sanitized form.
	r := env.NewRedactedEnv(sanitized, []string{"DB_PASSWORD"})
	safe := r.Safe()

	if safe["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected DB_PASSWORD redacted, got %q", safe["DB_PASSWORD"])
	}
	if safe["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", safe["APP_HOST"])
	}
}

// TestRedactedEnv_SafeCanBeInjected verifies that the Safe() output can be
// fed directly into an Injector without leaking secret values.
func TestRedactedEnv_SafeCanBeInjected(t *testing.T) {
	base := []string{"PATH=/usr/bin", "HOME=/root"}
	secrets := map[string]string{
		"DB_PASS": "real-secret",
		"APP_ENV": "production",
	}

	inj := env.NewInjector(base, secrets)
	full := make(map[string]string)
	for _, kv := range inj.Environ() {
		for i := 0; i < len(kv); i++ {
			if kv[i] == '=' {
				full[kv[:i]] = kv[i+1:]
				break
			}
		}
	}

	r := env.NewRedactedEnv(full, []string{"DB_PASS"})
	safe := r.Safe()

	if safe["DB_PASS"] != "[REDACTED]" {
		t.Errorf("expected DB_PASS redacted in injected env, got %q", safe["DB_PASS"])
	}
	if safe["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", safe["APP_ENV"])
	}
}
