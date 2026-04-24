package env

import (
	"strings"
	"testing"
)

func TestTruncateMap_ShortValuesUnchanged(t *testing.T) {
	m := map[string]string{"KEY": "short"}
	out := TruncateMap(m, TruncateOptions{MaxLen: 20})
	if out["KEY"] != "short" {
		t.Fatalf("expected 'short', got %q", out["KEY"])
	}
}

func TestTruncateMap_LongValueTruncated(t *testing.T) {
	m := map[string]string{"TOKEN": "supersecretvalue"}
	out := TruncateMap(m, TruncateOptions{MaxLen: 5, Suffix: "..."})
	if out["TOKEN"] != "super..." {
		t.Fatalf("expected 'super...', got %q", out["TOKEN"])
	}
}

func TestTruncateMap_DefaultSuffix(t *testing.T) {
	m := map[string]string{"K": "abcdefgh"}
	out := TruncateMap(m, TruncateOptions{MaxLen: 4})
	if !strings.HasSuffix(out["K"], "...") {
		t.Fatalf("expected default suffix '...', got %q", out["K"])
	}
	if out["K"] != "abcd..." {
		t.Fatalf("expected 'abcd...', got %q", out["K"])
	}
}

func TestTruncateMap_ZeroMaxLen_ReturnsOriginal(t *testing.T) {
	m := map[string]string{"K": "value"}
	out := TruncateMap(m, TruncateOptions{MaxLen: 0})
	if &out == &m {
		// same pointer is fine — documented behaviour
	}
	if out["K"] != "value" {
		t.Fatalf("expected unchanged value, got %q", out["K"])
	}
}

func TestTruncateMap_DoesNotMutateOriginal(t *testing.T) {
	m := map[string]string{"A": "longvalue123"}
	_ = TruncateMap(m, TruncateOptions{MaxLen: 4})
	if m["A"] != "longvalue123" {
		t.Fatal("original map was mutated")
	}
}

func TestTruncateValue_WithinLimit(t *testing.T) {
	v, err := TruncateValue("hello", 10, "...")
	if err != nil {
		t.Fatal(err)
	}
	if v != "hello" {
		t.Fatalf("expected 'hello', got %q", v)
	}
}

func TestTruncateValue_ExceedsLimit(t *testing.T) {
	v, err := TruncateValue("hello world", 5, "~")
	if err != nil {
		t.Fatal(err)
	}
	if v != "hello~" {
		t.Fatalf("expected 'hello~', got %q", v)
	}
}

func TestTruncateValue_NegativeMaxLen_ReturnsError(t *testing.T) {
	_, err := TruncateValue("value", -1, "...")
	if err == nil {
		t.Fatal("expected error for negative maxLen")
	}
}

func TestTruncateValue_ZeroMaxLen_Unchanged(t *testing.T) {
	v, err := TruncateValue("data", 0, "...")
	if err != nil {
		t.Fatal(err)
	}
	if v != "data" {
		t.Fatalf("expected 'data', got %q", v)
	}
}
