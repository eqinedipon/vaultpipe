package env_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/env"
)

// TestFlattenMap_ThenSanitize verifies that nested secret maps from Vault can be
// flattened into a dot-separated key space and then sanitized into valid
// environment variable names in a single pipeline.
func TestFlattenMap_ThenSanitize(t *testing.T) {
	nested := map[string]any{
		"database": map[string]any{
			"host": "db.internal",
			"port": "5432",
			"credentials": map[string]any{
				"username": "admin",
				"password": "s3cr3t",
			},
		},
		"api-key": "tok_abc123",
	}

	flat := env.FlattenMap(nested, env.FlattenKeys("_"))
	sanitized := env.SanitizeMap(flat)

	cases := []struct {
		key   string
		want  string
	}{
		{"DATABASE_HOST", "db.internal"},
		{"DATABASE_PORT", "5432"},
		{"DATABASE_CREDENTIALS_USERNAME", "admin"},
		{"DATABASE_CREDENTIALS_PASSWORD", "s3cr3t"},
		{"API_KEY", "tok_abc123"},
	}

	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			got, ok := sanitized[tc.key]
			if !ok {
				t.Fatalf("key %q not found in sanitized map; keys: %v", tc.key, keys(sanitized))
			}
			if got != tc.want {
				t.Errorf("sanitized[%q] = %q; want %q", tc.key, got, tc.want)
			}
		})
	}
}

// TestFlattenMap_ThenInject ensures that a flattened and sanitized secret map
// can be injected into a process environment via the Injector without collision.
func TestFlattenMap_ThenInject(t *testing.T) {
	nested := map[string]any{
		"service": map[string]any{
			"token": "bearer_xyz",
			"timeout": "30s",
		},
	}

	flat := env.FlattenMap(nested, env.FlattenKeys("_"))
	sanitized := env.SanitizeMap(flat)

	base := []string{"PATH=/usr/bin", "HOME=/root", "SERVICE_TOKEN=old_value"}
	injector := env.NewInjector(sanitized)
	result := injector.Environ(base)

	resultMap := toEnvMap(result)

	if got := resultMap["SERVICE_TOKEN"]; got != "bearer_xyz" {
		t.Errorf("SERVICE_TOKEN = %q; want %q", got, "bearer_xyz")
	}
	if got := resultMap["SERVICE_TIMEOUT"]; got != "30s" {
		t.Errorf("SERVICE_TIMEOUT = %q; want %q", got, "30s")
	}
	// Base-only keys must be preserved.
	if got := resultMap["PATH"]; got != "/usr/bin" {
		t.Errorf("PATH = %q; want %q", got, "/usr/bin")
	}
}

// TestFlattenMap_ThenValidate confirms that required keys are present after
// flattening a deeply nested Vault response.
func TestFlattenMap_ThenValidate(t *testing.T) {
	nested := map[string]any{
		"app": map[string]any{
			"secret": "abc",
		},
	}

	flat := env.FlattenMap(nested, env.FlattenKeys("_"))
	sanitized := env.SanitizeMap(flat)

	err := env.Validate(sanitized, env.RequireKeys("APP_SECRET"))
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}

	// Missing key should fail validation.
	err = env.Validate(sanitized, env.RequireKeys("APP_SECRET", "APP_MISSING"))
	if err == nil {
		t.Fatal("expected validation error for missing key APP_MISSING")
	}
	if !env.IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T: %v", err, err)
	}
}

// keys returns the sorted key slice of a string map for diagnostic output.
func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
