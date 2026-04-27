package env_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/your-org/vaultpipe/internal/env"
)

// TestDefaults_ThenSanitize verifies that defaults are applied and then keys
// are sanitized in a typical pipeline order.
func TestDefaults_ThenSanitize(t *testing.T) {
	base := map[string]string{"existing-key": "hello"}
	specs := []env.DefaultSpec{
		{Key: "existing-key", Value: "ignored"},
		{Key: "new-key", Value: "world"},
	}

	withDefaults, err := env.ApplyDefaults(base, specs)
	if err != nil {
		t.Fatalf("ApplyDefaults: %v", err)
	}

	sanitized := env.SanitizeMap(withDefaults)

	want := map[string]string{
		"EXISTING_KEY": "hello",
		"NEW_KEY":      "world",
	}
	if diff := cmp.Diff(want, sanitized); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

// TestDefaults_ThenValidate ensures that after applying defaults all required
// keys pass validation.
func TestDefaults_ThenValidate(t *testing.T) {
	base := map[string]string{"PORT": "8080"}
	specs := []env.DefaultSpec{
		{Key: "HOST", Value: "localhost"},
		{Key: "LOG_LEVEL", Value: "info"},
	}

	withDefaults, err := env.ApplyDefaults(base, specs)
	if err != nil {
		t.Fatalf("ApplyDefaults: %v", err)
	}

	err = env.Validate(withDefaults, env.RequireKeys("PORT", "HOST", "LOG_LEVEL"))
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

// TestDefaults_ChainIntegration exercises ApplyDefaults inside a Chain step.
func TestDefaults_ChainIntegration(t *testing.T) {
	defaultSpecs := []env.DefaultSpec{
		{Key: "TIMEOUT", Value: "30s"},
		{Key: "RETRIES", Value: "3"},
	}

	chain := env.NewChain(
		env.WrapTransformer(env.NewTransformer(env.TrimSpaceTransform)),
		env.Step(func(m map[string]string) (map[string]string, error) {
			return env.ApplyDefaults(m, defaultSpecs)
		}),
	)

	input := map[string]string{"APP": "  vaultpipe  ", "RETRIES": "5"}
	got, err := chain.Run(input)
	if err != nil {
		t.Fatalf("chain.Run: %v", err)
	}

	if got["APP"] != "vaultpipe" {
		t.Errorf("expected trimmed APP, got %q", got["APP"])
	}
	if got["TIMEOUT"] != "30s" {
		t.Errorf("expected TIMEOUT=30s, got %q", got["TIMEOUT"])
	}
	// RETRIES already set — must not be overwritten
	if got["RETRIES"] != "5" {
		t.Errorf("expected RETRIES=5 (preserved), got %q", got["RETRIES"])
	}
}
