// Package audit provides lightweight audit logging for secret access events.
package audit

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Path      string    `json:"path,omitempty"`
	Keys      []string  `json:"keys,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// Logger writes audit events as newline-delimited JSON.
type Logger struct {
	enc *json.Encoder
}

// NewLogger creates a Logger writing to w. Pass nil to use stderr.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{enc: json.NewEncoder(w)}
}

// LogFetch records a secret fetch event.
func (l *Logger) LogFetch(path string, keys []string) {
	l.write(Event{
		Timestamp: time.Now().UTC(),
		Action:    "fetch",
		Path:      path,
		Keys:      keys,
	})
}

// LogExec records a process execution event.
func (l *Logger) LogExec(cmd string) {
	l.write(Event{
		Timestamp: time.Now().UTC(),
		Action:    "exec",
		Path:      cmd,
	})
}

// LogError records an error event.
func (l *Logger) LogError(action, errMsg string) {
	l.write(Event{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Error:     errMsg,
	})
}

func (l *Logger) write(e Event) {
	// best-effort; ignore encode errors
	_ = l.enc.Encode(e)
}
