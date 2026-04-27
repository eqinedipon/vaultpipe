package env

import (
	"testing"
)

func TestRedactedEnv_SafeReplacesSecrets(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"DB_PASS":  "s3cr3t",
		"API_KEY":  "tok_abc123",
	}
	r := NewRedactedEnv(env, []string{"DB_PASS", "API_KEY"})
	safe := r.Safe()

	if safe["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", safe["APP_HOST"])
	}
	if safe["DB_PASS"] != "[REDACTED]" {
		t.Errorf("expected DB_PASS to be redacted, got %q", safe["DB_PASS"])
	}
	if safe["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected API_KEY to be redacted, got %q", safe["API_KEY"])
	}
}

func TestRedactedEnv_RawReturnsAllValues(t *testing.T) {
	env := map[string]string{
		"DB_PASS": "s3cr3t",
	}
	r := NewRedactedEnv(env, []string{"DB_PASS"})
	raw := r.Raw()

	if raw["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected raw DB_PASS=s3cr3t, got %q", raw["DB_PASS"])
	}
}

func TestRedactedEnv_DoesNotMutateInput(t *testing.T) {
	original := map[string]string{"TOKEN": "abc"}
	NewRedactedEnv(original, []string{"TOKEN"})
	if original["TOKEN"] != "abc" {
		t.Error("NewRedactedEnv mutated the input map")
	}
}

func TestRedactedEnv_CaseInsensitiveKeys(t *testing.T) {
	env := map[string]string{"db_pass": "secret"}
	r := NewRedactedEnv(env, []string{"DB_PASS"})
	safe := r.Safe()
	if safe["db_pass"] != "[REDACTED]" {
		t.Errorf("expected case-insensitive redaction, got %q", safe["db_pass"])
	}
}

func TestRedactedEnv_WithPlaceholder(t *testing.T) {
	env := map[string]string{"SECRET": "value"}
	r := NewRedactedEnv(env, []string{"SECRET"}).WithPlaceholder("***")
	safe := r.Safe()
	if safe["SECRET"] != "***" {
		t.Errorf("expected custom placeholder ***, got %q", safe["SECRET"])
	}
}

func TestRedactedEnv_IsSecret(t *testing.T) {
	r := NewRedactedEnv(map[string]string{}, []string{"MY_TOKEN"})
	if !r.IsSecret("MY_TOKEN") {
		t.Error("expected MY_TOKEN to be a secret")
	}
	if !r.IsSecret("my_token") {
		t.Error("expected case-insensitive IsSecret to return true")
	}
	if r.IsSecret("APP_HOST") {
		t.Error("expected APP_HOST to not be a secret")
	}
}

func TestRedactedEnv_EmptySecretKeys(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	r := NewRedactedEnv(env, nil)
	safe := r.Safe()
	if safe["FOO"] != "bar" {
		t.Errorf("expected FOO=bar with no secret keys, got %q", safe["FOO"])
	}
}
