package mask_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/mask"
)

func TestRedact_ReplacesSecret(t *testing.T) {
	m := mask.New([]string{"s3cr3t"})
	out := m.Redact("the password is s3cr3t, keep it safe")
	if want := "the password is ***REDACTED***, keep it safe"; out != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}

func TestRedact_MultipleSecrets(t *testing.T) {
	m := mask.New([]string{"alpha", "beta"})
	out := m.Redact("alpha and beta should both be hidden")
	if want := "***REDACTED*** and ***REDACTED*** should both be hidden"; out != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}

func TestRedact_NoMatch_Unchanged(t *testing.T) {
	m := mask.New([]string{"hidden"})
	input := "nothing sensitive here"
	if out := m.Redact(input); out != input {
		t.Fatalf("expected unchanged, got %q", out)
	}
}

func TestNew_IgnoresEmptyStrings(t *testing.T) {
	m := mask.New([]string{"", "real", ""})
	if m.Len() != 1 {
		t.Fatalf("expected 1 secret, got %d", m.Len())
	}
}

func TestAdd_AppendsSecrets(t *testing.T) {
	m := mask.New([]string{"first"})
	m.Add("second", "")
	if m.Len() != 2 {
		t.Fatalf("expected 2 secrets, got %d", m.Len())
	}
	out := m.Redact("first and second")
	if want := "***REDACTED*** and ***REDACTED***"; out != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}

func TestRedact_RepeatedOccurrences(t *testing.T) {
	m := mask.New([]string{"tok"})
	out := m.Redact("tok tok tok")
	if want := "***REDACTED*** ***REDACTED*** ***REDACTED***"; out != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}
