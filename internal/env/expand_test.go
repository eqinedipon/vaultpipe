package env

import (
	"os"
	"testing"
)

func TestExpand_ResolvesSecretReference(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "db.internal",
		"DB_PORT": "5432",
	}
	exp := NewExpander(secrets, nil)

	got := exp.Expand("postgres://${DB_HOST}:${DB_PORT}/app")
	want := "postgres://db.internal:5432/app"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestExpand_FallsBackToBase(t *testing.T) {
	base := []string{"REGION=us-east-1"}
	exp := NewExpander(nil, base)

	got := exp.Expand("region-${REGION}")
	want := "region-us-east-1"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestExpand_FallsBackToOSEnv(t *testing.T) {
	t.Setenv("_VP_TEST_VAR", "hello")
	exp := NewExpander(nil, nil)

	got := exp.Expand("say ${_VP_TEST_VAR}")
	want := "say hello"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestExpand_UnknownVarPreserved(t *testing.T) {
	exp := NewExpander(nil, nil)
	// Ensure the var is not set in the OS environment.
	os.Unsetenv("_VP_UNKNOWN_XYZ")

	got := exp.Expand("value=${_VP_UNKNOWN_XYZ}")
	want := "value=${_VP_UNKNOWN_XYZ}"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestExpand_SecretTakesPrecedenceOverBase(t *testing.T) {
	secrets := map[string]string{"HOST": "secret-host"}
	base := []string{"HOST=base-host"}
	exp := NewExpander(secrets, base)

	got := exp.Expand("${HOST}")
	if got != "secret-host" {
		t.Errorf("expected secret value, got %q", got)
	}
}

func TestExpandAll_ExpandsAllValues(t *testing.T) {
	secrets := map[string]string{
		"BASE_URL": "https://api.example.com",
		"FULL_URL": "${BASE_URL}/v1/resource",
	}
	exp := NewExpander(secrets, nil)
	out := exp.ExpandAll(secrets)

	if out["FULL_URL"] != "https://api.example.com/v1/resource" {
		t.Errorf("unexpected FULL_URL: %q", out["FULL_URL"])
	}
	// Original map must not be mutated.
	if secrets["FULL_URL"] != "${BASE_URL}/v1/resource" {
		t.Error("original secrets map was mutated")
	}
}

func TestExpandAll_EmptySecrets(t *testing.T) {
	exp := NewExpander(nil, nil)
	out := exp.ExpandAll(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
