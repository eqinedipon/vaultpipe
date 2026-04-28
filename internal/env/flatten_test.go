package env

import (
	"testing"
)

func TestFlattenMap_SimpleKeys(t *testing.T) {
	src := map[string]any{
		"host": "localhost",
		"port": 5432,
	}
	got, err := FlattenMap(src, FlattenOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["HOST"] != "localhost" {
		t.Errorf("HOST: got %q, want %q", got["HOST"], "localhost")
	}
	if got["PORT"] != "5432" {
		t.Errorf("PORT: got %q, want %q", got["PORT"], "5432")
	}
}

func TestFlattenMap_NestedKeys(t *testing.T) {
	src := map[string]any{
		"db": map[string]any{
			"host": "db.internal",
			"port": 5432,
		},
	}
	got, err := FlattenMap(src, FlattenOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_HOST"] != "db.internal" {
		t.Errorf("DB_HOST: got %q, want %q", got["DB_HOST"], "db.internal")
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("DB_PORT: got %q, want %q", got["DB_PORT"], "5432")
	}
}

func TestFlattenMap_WithPrefix(t *testing.T) {
	src := map[string]any{"token": "abc123"}
	got, err := FlattenMap(src, FlattenOptions{Prefix: "vault", UpperCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["VAULT_TOKEN"] != "abc123" {
		t.Errorf("VAULT_TOKEN: got %q, want %q", got["VAULT_TOKEN"], "abc123")
	}
}

func TestFlattenMap_NilValue(t *testing.T) {
	src := map[string]any{"empty": nil}
	got, err := FlattenMap(src, FlattenOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := got["EMPTY"]; !ok || v != "" {
		t.Errorf("EMPTY: got %q (ok=%v), want empty string", v, ok)
	}
}

func TestFlattenMap_CustomSeparator(t *testing.T) {
	src := map[string]any{
		"aws": map[string]any{"region": "us-east-1"},
	}
	got, err := FlattenMap(src, FlattenOptions{Separator: ".", UpperCase: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["aws.region"] != "us-east-1" {
		t.Errorf("aws.region: got %q, want %q", got["aws.region"], "us-east-1")
	}
}

func TestFlattenKeys_Sorted(t *testing.T) {
	src := map[string]any{"z": 1, "a": 2, "m": 3}
	keys, err := FlattenKeys(src, FlattenOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}

func TestFlattenMap_DeepNesting(t *testing.T) {
	src := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
		},
	}
	got, err := FlattenMap(src, FlattenOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["A_B_C"] != "deep" {
		t.Errorf("A_B_C: got %q, want %q", got["A_B_C"], "deep")
	}
}
