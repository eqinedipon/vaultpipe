package env_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/env"
)

// TestTypeCast_ThenSanitize verifies that TypeCaster output can be fed
// directly into SanitizeMap without loss of normalized values.
func TestTypeCast_ThenSanitize(t *testing.T) {
	tc := env.NewTypeCaster([]env.TypeCastRule{
		{Key: "feature-enabled", TypeName: "bool"},
		{Key: "max-retries", TypeName: "int"},
	})

	src := map[string]string{
		"feature-enabled": "yes",
		"max-retries":     "3.0",
		"app-name":        "vaultpipe",
	}

	casted, err := tc.Cast(src)
	if err != nil {
		t.Fatalf("Cast: %v", err)
	}

	sanitized := env.SanitizeMap(casted)

	if sanitized["FEATURE_ENABLED"] != "true" {
		t.Errorf("FEATURE_ENABLED: got %q, want \"true\"", sanitized["FEATURE_ENABLED"])
	}
	if sanitized["MAX_RETRIES"] != "3" {
		t.Errorf("MAX_RETRIES: got %q, want \"3\"", sanitized["MAX_RETRIES"])
	}
	if sanitized["APP_NAME"] != "vaultpipe" {
		t.Errorf("APP_NAME: got %q, want \"vaultpipe\"", sanitized["APP_NAME"])
	}
}

// TestTypeCast_ThenInject ensures that type-cast values are correctly
// injected into a process environment via the Injector.
func TestTypeCast_ThenInject(t *testing.T) {
	tc := env.NewTypeCaster([]env.TypeCastRule{
		{Key: "VERBOSE", TypeName: "bool"},
		{Key: "WORKERS", TypeName: "int"},
	})

	secrets := map[string]string{
		"VERBOSE": "on",
		"WORKERS": "4.0",
	}

	normalized, err := tc.Cast(secrets)
	if err != nil {
		t.Fatalf("Cast: %v", err)
	}

	inj := env.NewInjector(normalized)
	result := inj.Environ(nil)

	m := make(map[string]string)
	for _, kv := range result {
		for i := 0; i < len(kv); i++ {
			if kv[i] == '=' {
				m[kv[:i]] = kv[i+1:]
				break
			}
		}
	}

	if m["VERBOSE"] != "true" {
		t.Errorf("VERBOSE: got %q, want \"true\"", m["VERBOSE"])
	}
	if m["WORKERS"] != "4" {
		t.Errorf("WORKERS: got %q, want \"4\"", m["WORKERS"])
	}
}
