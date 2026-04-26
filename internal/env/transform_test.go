package env

import (
	"errors"
	"strings"
	"testing"
)

func TestTransformer_TrimSpace(t *testing.T) {
	tr := NewTransformer(TrimSpaceTransform())
	src := map[string]string{
		"FOO": "  hello  ",
		"BAR": "\tworld\n",
	}
	out, err := tr.Apply(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "hello" {
		t.Errorf("FOO: got %q, want %q", out["FOO"], "hello")
	}
	if out["BAR"] != "world" {
		t.Errorf("BAR: got %q, want %q", out["BAR"], "world")
	}
}

func TestTransformer_DoesNotMutateSrc(t *testing.T) {
	tr := NewTransformer(TrimSpaceTransform())
	src := map[string]string{"K": "  v  "}
	_, _ = tr.Apply(src)
	if src["K"] != "  v  " {
		t.Error("Apply mutated the source map")
	}
}

func TestTransformer_ChainedFunctions(t *testing.T) {
	upper := func(_, v string) (string, error) { return strings.ToUpper(v), nil }
	append_ := func(_, v string) (string, error) { return v + "!", nil }
	tr := NewTransformer(upper, append_)
	out, err := tr.Apply(map[string]string{"X": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != "HELLO!" {
		t.Errorf("got %q, want %q", out["X"], "HELLO!")
	}
}

func TestTransformer_ErrorAbortsEarly(t *testing.T) {
	boom := func(k, _ string) (string, error) {
		return "", errors.New("boom")
	}
	tr := NewTransformer(boom)
	_, err := tr.Apply(map[string]string{"A": "val"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "A") {
		t.Errorf("expected key name in error, got: %v", err)
	}
}

func TestUpperKeyTransform_MatchingSuffix(t *testing.T) {
	tr := NewTransformer(UpperKeyTransform("_FLAG"))
	out, err := tr.Apply(map[string]string{
		"FEATURE_FLAG": "enabled",
		"SECRET_KEY":   "abc123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FEATURE_FLAG"] != "ENABLED" {
		t.Errorf("FEATURE_FLAG: got %q, want %q", out["FEATURE_FLAG"], "ENABLED")
	}
	if out["SECRET_KEY"] != "abc123" {
		t.Errorf("SECRET_KEY should be unchanged, got %q", out["SECRET_KEY"])
	}
}

func TestTransformer_EmptyMap(t *testing.T) {
	tr := NewTransformer(TrimSpaceTransform())
	out, err := tr.Apply(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}
