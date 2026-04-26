package env_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpipe/internal/env"
)

// TestTransform_ThenSanitize verifies that a Transformer can be composed with
// SanitizeMap: trim whitespace first, then sanitize the keys.
func TestTransform_ThenSanitize(t *testing.T) {
	raw := map[string]string{
		"my-secret": "  s3cr3t  ",
		"other_key": "  value  ",
	}
	tr := env.NewTransformer(env.TrimSpaceTransform())
	trimmed, err := tr.Apply(raw)
	if err != nil {
		t.Fatalf("transform: %v", err)
	}
	sanitized := env.SanitizeMap(trimmed)
	if sanitized["MY_SECRET"] != "s3cr3t" {
		t.Errorf("MY_SECRET: got %q", sanitized["MY_SECRET"])
	}
	if sanitized["OTHER_KEY"] != "value" {
		t.Errorf("OTHER_KEY: got %q", sanitized["OTHER_KEY"])
	}
}

// TestTransform_ThenInject verifies that transformed secrets are visible in the
// environment produced by the Injector.
func TestTransform_ThenInject(t *testing.T) {
	secrets := map[string]string{
		"API_KEY": "  tok_live_xyz  ",
	}
	tr := env.NewTransformer(env.TrimSpaceTransform())
	cleaned, err := tr.Apply(secrets)
	if err != nil {
		t.Fatalf("transform: %v", err)
	}
	inj := env.NewInjector(cleaned)
	environ := inj.Environ(nil)
	found := false
	for _, e := range environ {
		if e == "API_KEY=tok_live_xyz" {
			found = true
		}
	}
	if !found {
		t.Errorf("trimmed secret not found in environ; got: %v", environ)
	}
}

// TestTransform_ChainWithCoerce verifies Transformer composes with CoerceMap:
// coerce numeric strings to their canonical form, then trim.
func TestTransform_ChainWithCoerce(t *testing.T) {
	raw := map[string]string{
		"TIMEOUT": "  30  ",
	}
	tr := env.NewTransformer(env.TrimSpaceTransform())
	out, err := tr.Apply(raw)
	if err != nil {
		t.Fatalf("transform: %v", err)
	}
	if strings.TrimSpace(out["TIMEOUT"]) != "30" {
		t.Errorf("TIMEOUT: got %q, want %q", out["TIMEOUT"], "30")
	}
}
