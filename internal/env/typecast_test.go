package env

import (
	"testing"
)

func TestCast_BoolNormalization(t *testing.T) {
	tc := NewTypeCaster([]TypeCastRule{
		{Key: "FEATURE_FLAG", TypeName: "bool"},
	})

	cases := []struct {
		input string
		want  string
	}{
		{"true", "true"}, {"1", "true"}, {"yes", "true"}, {"on", "true"}, {"enabled", "true"},
		{"false", "false"}, {"0", "false"}, {"no", "false"}, {"off", "false"}, {"disabled", "false"},
	}

	for _, c := range cases {
		src := map[string]string{"FEATURE_FLAG": c.input}
		out, err := tc.Cast(src)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", c.input, err)
		}
		if got := out["FEATURE_FLAG"]; got != c.want {
			t.Errorf("input %q: got %q, want %q", c.input, got, c.want)
		}
	}
}

func TestCast_IntNormalization(t *testing.T) {
	tc := NewTypeCaster([]TypeCastRule{
		{Key: "TIMEOUT", TypeName: "int"},
	})

	src := map[string]string{"TIMEOUT": "30.0"}
	out, err := tc.Cast(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TIMEOUT"] != "30" {
		t.Errorf("got %q, want \"30\"", out["TIMEOUT"])
	}
}

func TestCast_FloatNormalization(t *testing.T) {
	tc := NewTypeCaster([]TypeCastRule{
		{Key: "RATIO", TypeName: "float"},
	})

	src := map[string]string{"RATIO": "0.500"}
	out, err := tc.Cast(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["RATIO"] != "0.5" {
		t.Errorf("got %q, want \"0.5\"", out["RATIO"])
	}
}

func TestCast_UnknownKeyPassThrough(t *testing.T) {
	tc := NewTypeCaster(nil)
	src := map[string]string{"MY_VAR": "hello"}
	out, err := tc.Cast(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["MY_VAR"] != "hello" {
		t.Errorf("expected pass-through, got %q", out["MY_VAR"])
	}
}

func TestCast_InvalidBool_ReturnsError(t *testing.T) {
	tc := NewTypeCaster([]TypeCastRule{
		{Key: "FLAG", TypeName: "bool"},
	})
	_, err := tc.Cast(map[string]string{"FLAG": "maybe"})
	if err == nil {
		t.Fatal("expected error for invalid bool, got nil")
	}
}

func TestCast_InvalidInt_ReturnsError(t *testing.T) {
	tc := NewTypeCaster([]TypeCastRule{
		{Key: "COUNT", TypeName: "int"},
	})
	_, err := tc.Cast(map[string]string{"COUNT": "abc"})
	if err == nil {
		t.Fatal("expected error for invalid int, got nil")
	}
}

func TestCast_KeyMatchIsCaseInsensitive(t *testing.T) {
	tc := NewTypeCaster([]TypeCastRule{
		{Key: "debug", TypeName: "bool"},
	})
	out, err := tc.Cast(map[string]string{"DEBUG": "yes"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DEBUG"] != "true" {
		t.Errorf("got %q, want \"true\"", out["DEBUG"])
	}
}

func TestCast_DoesNotMutateInput(t *testing.T) {
	tc := NewTypeCaster([]TypeCastRule{
		{Key: "ENABLED", TypeName: "bool"},
	})
	src := map[string]string{"ENABLED": "1", "OTHER": "val"}
	_, err := tc.Cast(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src["ENABLED"] != "1" {
		t.Error("Cast mutated the input map")
	}
}
