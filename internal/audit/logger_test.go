package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/audit"
)

func decode(t *testing.T, buf *bytes.Buffer) audit.Event {
	t.Helper()
	var e audit.Event
	if err := json.NewDecoder(buf).Decode(&e); err != nil {
		t.Fatalf("decode event: %v", err)
	}
	return e
}

func TestLogFetch(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)
	l.LogFetch("secret/data/app", []string{"DB_PASS", "API_KEY"})

	e := decode(t, &buf)
	if e.Action != "fetch" {
		t.Errorf("expected action=fetch, got %q", e.Action)
	}
	if e.Path != "secret/data/app" {
		t.Errorf("unexpected path: %q", e.Path)
	}
	if len(e.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(e.Keys))
	}
	if e.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestLogExec(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)
	l.LogExec("/usr/bin/env")

	e := decode(t, &buf)
	if e.Action != "exec" {
		t.Errorf("expected action=exec, got %q", e.Action)
	}
	if e.Path != "/usr/bin/env" {
		t.Errorf("unexpected path: %q", e.Path)
	}
}

func TestLogError(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)
	l.LogError("fetch", "permission denied")

	e := decode(t, &buf)
	if e.Action != "fetch" {
		t.Errorf("expected action=fetch, got %q", e.Action)
	}
	if !strings.Contains(e.Error, "permission denied") {
		t.Errorf("unexpected error field: %q", e.Error)
	}
}

func TestNewLogger_NilUsesStderr(t *testing.T) {
	// Just ensure it doesn't panic
	l := audit.NewLogger(nil)
	l.LogExec("test")
}
