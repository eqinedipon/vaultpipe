package env

import (
	"testing"
)

func TestMerge_OverwriteByDefault(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	src := map[string]string{"B": "99", "C": "3"}

	got := MergeTwo(base, src)

	if got["A"] != "1" {
		t.Errorf("A: want 1, got %s", got["A"])
	}
	if got["B"] != "99" {
		t.Errorf("B: want 99, got %s", got["B"])
	}
	if got["C"] != "3" {
		t.Errorf("C: want 3, got %s", got["C"])
	}
}

func TestMerge_NoOverwrite_PreservesBase(t *testing.T) {
	base := map[string]string{"A": "original", "B": "2"}
	src := map[string]string{"A": "overridden", "C": "new"}

	got := MergeTwo(base, src, WithNoOverwrite())

	if got["A"] != "original" {
		t.Errorf("A: want original, got %s", got["A"])
	}
	if got["C"] != "new" {
		t.Errorf("C: want new, got %s", got["C"])
	}
}

func TestMerge_DoesNotMutateBase(t *testing.T) {
	base := map[string]string{"A": "1"}
	src := map[string]string{"A": "2", "B": "3"}

	_ = MergeTwo(base, src)

	if base["A"] != "1" {
		t.Errorf("base mutated: A = %s", base["A"])
	}
	if _, ok := base["B"]; ok {
		t.Error("base mutated: unexpected key B")
	}
}

func TestMerge_MultipleSources(t *testing.T) {
	base := map[string]string{"A": "base"}
	s1 := map[string]string{"A": "s1", "B": "s1"}
	s2 := map[string]string{"B": "s2", "C": "s2"}

	got := Merge(base, []map[string]string{s1, s2})

	if got["A"] != "s1" {
		t.Errorf("A: want s1, got %s", got["A"])
	}
	if got["B"] != "s2" {
		t.Errorf("B: want s2, got %s", got["B"])
	}
	if got["C"] != "s2" {
		t.Errorf("C: want s2, got %s", got["C"])
	}
}

func TestMerge_EmptyBase(t *testing.T) {
	got := MergeTwo(map[string]string{}, map[string]string{"X": "1"})
	if got["X"] != "1" {
		t.Errorf("X: want 1, got %s", got["X"])
	}
}

func TestMerge_EmptySource(t *testing.T) {
	base := map[string]string{"A": "1"}
	got := MergeTwo(base, map[string]string{})
	if got["A"] != "1" {
		t.Errorf("A: want 1, got %s", got["A"])
	}
	if len(got) != 1 {
		t.Errorf("unexpected extra keys in result: %v", got)
	}
}
