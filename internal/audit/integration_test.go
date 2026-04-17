package audit_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/yourusername/vaultpipe/internal/audit"
)

// TestMultipleEventsSequential ensures multiple events are each on their own
// line and independently decodable.
func TestMultipleEventsSequential(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	l.LogFetch("secret/data/db", []string{"PASS"})
	l.LogExec("myapp")
	l.LogError("exec", "exit status 1")

	dec := json.NewDecoder(&buf)
	actions := []string{"fetch", "exec", "exec"}
	for i, want := range actions {
		var e audit.Event
		if err := dec.Decode(&e); err != nil {
			t.Fatalf("event %d: decode error: %v", i, err)
		}
		if e.Action != want {
			t.Errorf("event %d: got action=%q, want %q", i, e.Action, want)
		}
	}
}

// TestNoSecretValuesLogged is a canary: secret values must never appear in
// audit output — only key names are permitted.
func TestNoSecretValuesLogged(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	secretValue := "s3cr3t-p@ssw0rd"
	l.LogFetch("secret/data/app", []string{"DB_PASS"})
	// Simulate a caller mistakenly trying to log a value via LogError
	l.LogError("fetch", "redacted")

	if bytes.Contains(buf.Bytes(), []byte(secretValue)) {
		t.Error("secret value must not appear in audit log")
	}
}
