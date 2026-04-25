package env

import (
	"testing"
)

func TestCoerceValue_String(t *testing.T) {
	got, err := CoerceValue("hello")
	if err != nil || got != "hello" {
		t.Fatalf("expected \"hello\", got %q err %v", got, err)
	}
}

func TestCoerceValue_Bool(t *testing.T) {
	for _, tc := range []struct{ in bool; want string }{
		{true, "true"}, {false, "false"},
	} {
		got, err := CoerceValue(tc.in)
		if err != nil || got != tc.want {
			t.Errorf("CoerceValue(%v) = %q, %v; want %q", tc.in, got, err, tc.want)
		}
	}
}

func TestCoerceValue_Int(t *testing.T) {
	got, err := CoerceValue(42)
	if err != nil || got != "42" {
		t.Fatalf("expected \"42\", got %q err %v", got, err)
	}
}

func TestCoerceValue_Int64(t *testing.T) {
	got, err := CoerceValue(int64(9876543210))
	if err != nil || got != "9876543210" {
		t.Fatalf("expected \"9876543210\", got %q err %v", got, err)
	}
}

func TestCoerceValue_Float64(t *testing.T) {
	got, err := CoerceValue(float64(3.14))
	if err != nil || got != "3.14" {
		t.Fatalf("expected \"3.14\", got %q err %v", got, err)
	}
}

func TestCoerceValue_Nil(t *testing.T) {
	got, err := CoerceValue(nil)
	if err != nil || got != "" {
		t.Fatalf("expected \"\", got %q err %v", got, err)
	}
}

func TestCoerceValue_FallbackStringer(t *testing.T) {
	type custom struct{ V int }
	got, err := CoerceValue(custom{V: 7})
	if err != nil || got != "{7}" {
		t.Fatalf("expected \"{7}\", got %q err %v", got, err)
	}
}

func TestCoerceValue_StringifyAll(t *testing.T) {
	got, err := CoerceValue(true, StringifyAll())
	if err != nil || got != "true" {
		t.Fatalf("StringifyAll bool: got %q err %v", got, err)
	}
}

func TestCoerceMap_AllTypes(t *testing.T) {
	input := map[string]any{
		"STR":   "value",
		"NUM":   int(10),
		"FLOAT": float64(1.5),
		"BOOL":  true,
		"NIL":   nil,
	}
	out, err := CoerceMap(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := map[string]string{
		"STR":   "value",
		"NUM":   "10",
		"FLOAT": "1.5",
		"BOOL":  "true",
		"NIL":   "",
	}
	for k, want := range expected {
		if got := out[k]; got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestCoerceMap_Empty(t *testing.T) {
	out, err := CoerceMap(map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}
