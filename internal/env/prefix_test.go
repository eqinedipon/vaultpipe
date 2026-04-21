package env

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStrip_RemovesPrefixFromMatchingKeys(t *testing.T) {
	m := map[string]string{
		"VAULT_DB_HOST": "localhost",
		"VAULT_DB_PORT": "5432",
		"OTHER_KEY":     "ignored",
	}
	pm := NewPrefixMapper("VAULT_")
	got := pm.Strip(m)
	want := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Strip() mismatch (-want +got):\n%s", diff)
	}
}

func TestStrip_DropsPrefixOnlyKey(t *testing.T) {
	// A key equal to the prefix itself should be dropped (empty remainder).
	m := map[string]string{
		"PREFIX_": "value",
		"PREFIX_A": "keep",
	}
	pm := NewPrefixMapper("PREFIX_")
	got := pm.Strip(m)
	want := map[string]string{"A": "keep"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Strip() mismatch (-want +got):\n%s", diff)
	}
}

func TestAdd_PrependsPrefixToAllKeys(t *testing.T) {
	m := map[string]string{
		"HOST": "db",
		"PORT": "5432",
	}
	pm := NewPrefixMapper("APP_")
	got := pm.Add(m)
	want := map[string]string{
		"APP_HOST": "db",
		"APP_PORT": "5432",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Add() mismatch (-want +got):\n%s", diff)
	}
}

func TestAdd_EmptyMap(t *testing.T) {
	pm := NewPrefixMapper("X_")
	got := pm.Add(map[string]string{})
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestFilterByPrefix_ReturnsOnlyMatching(t *testing.T) {
	m := map[string]string{
		"SECRET_A": "1",
		"SECRET_B": "2",
		"OTHER":    "3",
	}
	got := FilterByPrefix(m, "SECRET_")
	want := map[string]string{
		"SECRET_A": "1",
		"SECRET_B": "2",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("FilterByPrefix() mismatch (-want +got):\n%s", diff)
	}
}

func TestFilterByPrefix_NoMatch(t *testing.T) {
	m := map[string]string{"FOO": "bar"}
	got := FilterByPrefix(m, "NOPE_")
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}
