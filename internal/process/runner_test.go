package process

import (
	"os"
	"testing"
)

func TestRun_Success(t *testing.T) {
	env := os.Environ()
	r := NewRunner(env)

	code, err := r.Run("echo", []string{"hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestRun_NonZeroExit(t *testing.T) {
	r := NewRunner(os.Environ())

	code, err := r.Run("sh", []string{"-c", "exit 42"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 42 {
		t.Fatalf("expected exit code 42, got %d", code)
	}
}

func TestRun_CommandNotFound(t *testing.T) {
	r := NewRunner(os.Environ())

	_, err := r.Run("__no_such_command_vaultpipe__", nil)
	if err == nil {
		t.Fatal("expected error for missing command, got nil")
	}
}

func TestRun_EnvInjected(t *testing.T) {
	env := append(os.Environ(), "VAULTPIPE_TEST_VAR=injected")
	r := NewRunner(env)

	// sh -c 'test "$VAULTPIPE_TEST_VAR" = "injected"' exits 0 on match.
	code, err := r.Run("sh", []string{"-c", `test "$VAULTPIPE_TEST_VAR" = "injected"`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("env var not injected; exit code %d", code)
	}
}
