package env

import (
	"strings"
	"testing"
)

func TestResolve_PrefersSecrets(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "s3cr3t"}
	base := map[string]string{"DB_PASS": "plain"}

	out, err := Resolve([]string{"DB_PASS"}, secrets, base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected secret value, got %q", out["DB_PASS"])
	}
}

func TestResolve_FallsBackToBase(t *testing.T) {
	out, err := Resolve([]string{"APP_ENV"}, nil, map[string]string{"APP_ENV": "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected base value, got %q", out["APP_ENV"])
	}
}

func TestResolve_FallsBackToOSLookup(t *testing.T) {
	original := lookupEnv
	t.Cleanup(func() { lookupEnv = original })
	lookupEnv = func(key string) (string, bool) {
		if key == "FROM_OS" {
			return "os-value", true
		}
		return "", false
	}

	out, err := Resolve([]string{"FROM_OS"}, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FROM_OS"] != "os-value" {
		t.Errorf("expected os value, got %q", out["FROM_OS"])
	}
}

func TestResolve_Strict_MissingKeyReturnsError(t *testing.T) {
	_, err := Resolve([]string{"MISSING"}, nil, nil, Strict)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "MISSING") {
		t.Errorf("error should mention key name, got: %v", err)
	}
}

func TestResolve_Strict_AllowEmpty_DoesNotError(t *testing.T) {
	_, err := Resolve([]string{"EMPTY"}, map[string]string{"EMPTY": ""}, nil, Strict, AllowEmpty)
	if err != nil {
		t.Fatalf("expected no error with AllowEmpty, got: %v", err)
	}
}

func TestResolve_NonStrict_MissingKeyIsEmptyString(t *testing.T) {
	out, err := Resolve([]string{"GHOST"}, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := out["GHOST"]; !ok || v != "" {
		t.Errorf("expected empty string entry, got ok=%v val=%q", ok, v)
	}
}

func TestResolve_MultipleMissingListedTogether(t *testing.T) {
	_, err := Resolve([]string{"A", "B", "C"}, nil, nil, Strict)
	if err == nil {
		t.Fatal("expected error")
	}
	for _, key := range []string{"A", "B", "C"} {
		if !strings.Contains(err.Error(), key) {
			t.Errorf("error should mention %q, got: %v", key, err)
		}
	}
}
