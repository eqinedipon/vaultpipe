package env

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestApplyDefaults_FillsMissingKeys(t *testing.T) {
	dst := map[string]string{"A": "1"}
	specs := []DefaultSpec{{Key: "A", Value: "99"}, {Key: "B", Value: "2"}}

	got, err := ApplyDefaults(dst, specs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]string{"A": "1", "B": "2"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestApplyDefaults_OverrideExisting(t *testing.T) {
	dst := map[string]string{"A": "1"}
	specs := []DefaultSpec{{Key: "A", Value: "99"}}

	got, err := ApplyDefaults(dst, specs, OverrideExisting())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["A"] != "99" {
		t.Errorf("expected A=99, got %q", got["A"])
	}
}

func TestApplyDefaults_EmptyKeyReturnsError(t *testing.T) {
	_, err := ApplyDefaults(map[string]string{}, []DefaultSpec{{Key: "", Value: "x"}})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestApplyDefaults_DoesNotMutateDst(t *testing.T) {
	dst := map[string]string{"X": "original"}
	specs := []DefaultSpec{{Key: "Y", Value: "new"}}

	_, err := ApplyDefaults(dst, specs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := dst["Y"]; ok {
		t.Error("ApplyDefaults mutated the dst map")
	}
}

func TestApplyDefaults_EmptySpecs_ReturnsCopy(t *testing.T) {
	dst := map[string]string{"K": "v"}
	got, err := ApplyDefaults(dst, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(dst, got); diff != "" {
		t.Errorf("expected identical copy, diff:\n%s", diff)
	}
}

func TestDefaultsFromMap_ConvertsMap(t *testing.T) {
	m := map[string]string{"FOO": "bar", "BAZ": "qux"}
	specs := DefaultsFromMap(m)
	if len(specs) != 2 {
		t.Fatalf("expected 2 specs, got %d", len(specs))
	}
	reconstructed := make(map[string]string, len(specs))
	for _, s := range specs {
		reconstructed[s.Key] = s.Value
	}
	if diff := cmp.Diff(m, reconstructed); diff != "" {
		t.Errorf("round-trip mismatch:\n%s", diff)
	}
}
