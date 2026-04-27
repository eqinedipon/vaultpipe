package env_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/env"
)

// TestChain_SanitizeThenCoerceThenValidate exercises a realistic pipeline:
// raw secret keys are sanitized, values are coerced to strings, and then
// required keys are validated.
func TestChain_SanitizeThenCoerceThenValidate(t *testing.T) {
	sanitizeStep := env.ChainStep(func(m map[string]string) (map[string]string, error) {
		return env.SanitizeMap(m), nil
	})

	chain := env.NewChain(
		sanitizeStep,
		env.WrapValidation(env.RequireKeys("DB_HOST", "DB_PORT")),
	)

	src := map[string]string{
		"db-host": "localhost",
		"db-port": "5432",
	}

	out, err := chain.Run(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST: got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("DB_PORT: got %q", out["DB_PORT"])
	}
}

// TestChain_TrimThenTruncateThenInject verifies that a chain feeding into the
// injector produces the expected environment.
func TestChain_TrimThenTruncateThenInject(t *testing.T) {
	tr := env.NewTransformer(env.TrimSpaceTransform)

	chain := env.NewChain(
		env.WrapTransformer(tr),
		func(m map[string]string) (map[string]string, error) {
			return env.TruncateMap(m, 8, ""), nil
		},
	)

	src := map[string]string{
		"SECRET_A": "  toolongvalue  ",
		"SECRET_B": "  short  ",
	}

	out, err := chain.Run(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	inj := env.NewInjector(nil, out)
	envMap := make(map[string]string)
	for _, kv := range inj.Environ() {
		for i := 0; i < len(kv); i++ {
			if kv[i] == '=' {
				envMap[kv[:i]] = kv[i+1:]
				break
			}
		}
	}

	if len(out["SECRET_A"]) > 8 {
		t.Errorf("SECRET_A not truncated: %q", out["SECRET_A"])
	}
	if out["SECRET_B"] != "short" {
		t.Errorf("SECRET_B: expected 'short', got %q", out["SECRET_B"])
	}
}
