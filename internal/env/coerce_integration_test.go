package env_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/env"
)

// TestCoerceMap_ThenSanitize verifies that coerced values survive the
// SanitizeMap pipeline without data loss.
func TestCoerceMap_ThenSanitize(t *testing.T) {
	raw := map[string]any{
		"my-key":   "secret",
		"count":    int(3),
		"enabled":  true,
		"ratio":    float64(0.75),
	}

	coerced, err := env.CoerceMap(raw)
	if err != nil {
		t.Fatalf("CoerceMap: %v", err)
	}

	sanitized := env.SanitizeMap(coerced)

	cases := map[string]string{
		"MY_KEY":  "secret",
		"COUNT":   "3",
		"ENABLED": "true",
		"RATIO":   "0.75",
	}
	for k, want := range cases {
		if got, ok := sanitized[k]; !ok || got != want {
			t.Errorf("sanitized[%q] = %q (ok=%v), want %q", k, got, ok, want)
		}
	}
}

// TestCoerceMap_ThenInject verifies coerced values can be injected into a
// process environment via the Injector.
func TestCoerceMap_ThenInject(t *testing.T) {
	raw := map[string]any{
		"DB_PORT": int(5432),
		"VERBOSE": false,
	}

	coerced, err := env.CoerceMap(raw)
	if err != nil {
		t.Fatalf("CoerceMap: %v", err)
	}

	inj := env.NewInjector(coerced)
	result := inj.Environ(nil)

	rm := toEnvMap(result)
	if rm["DB_PORT"] != "5432" {
		t.Errorf("DB_PORT: got %q, want \"5432\"", rm["DB_PORT"])
	}
	if rm["VERBOSE"] != "false" {
		t.Errorf("VERBOSE: got %q, want \"false\"", rm["VERBOSE"])
	}
}
