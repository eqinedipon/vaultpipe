package env

import "testing"

func TestFilter_AllowPrefix(t *testing.T) {
	s := NewSnapshotFromMap(map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"HOME":     "/home/user",
	})

	filtered := s.Filter(AllowPrefix("APP_"))
	keys := filtered.Keys()

	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(keys), keys)
	}
	if _, ok := filtered.Get("HOME"); ok {
		t.Error("HOME should have been filtered out")
	}
	if v, ok := filtered.Get("APP_HOST"); !ok || v != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", v)
	}
}

func TestFilter_DenyList(t *testing.T) {
	s := NewSnapshotFromMap(map[string]string{
		"SECRET_TOKEN": "s3cr3t",
		"DB_URL":       "postgres://",
		"LOG_LEVEL":    "debug",
	})

	filtered := s.Filter(DenyList("SECRET_TOKEN"))

	if _, ok := filtered.Get("SECRET_TOKEN"); ok {
		t.Error("SECRET_TOKEN should be denied")
	}
	if v, _ := filtered.Get("DB_URL"); v != "postgres://" {
		t.Errorf("DB_URL should pass through, got %q", v)
	}
}

func TestFilter_MultipleOptions_AllMustPass(t *testing.T) {
	s := NewSnapshotFromMap(map[string]string{
		"APP_SECRET": "hidden",
		"APP_HOST":   "localhost",
		"OTHER":      "value",
	})

	filtered := s.Filter(AllowPrefix("APP_"), DenyList("APP_SECRET"))
	keys := filtered.Keys()

	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d: %v", len(keys), keys)
	}
	if v, _ := filtered.Get("APP_HOST"); v != "localhost" {
		t.Errorf("expected APP_HOST=localhost")
	}
}

func TestMerge_OverridesTakePrecedence(t *testing.T) {
	base := NewSnapshotFromMap(map[string]string{
		"A": "original",
		"B": "keep",
	})

	merged := base.Merge(map[string]string{
		"A": "overridden",
		"C": "new",
	})

	if v, _ := merged.Get("A"); v != "overridden" {
		t.Errorf("expected A=overridden, got %q", v)
	}
	if v, _ := merged.Get("B"); v != "keep" {
		t.Errorf("expected B=keep, got %q", v)
	}
	if v, _ := merged.Get("C"); v != "new" {
		t.Errorf("expected C=new, got %q", v)
	}
}

func TestMerge_DoesNotMutateOriginal(t *testing.T) {
	base := NewSnapshotFromMap(map[string]string{"A": "1"})
	_ = base.Merge(map[string]string{"A": "mutated"})

	if v, _ := base.Get("A"); v != "1" {
		t.Error("original snapshot was mutated")
	}
}
