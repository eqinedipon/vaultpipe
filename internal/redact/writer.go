// Package redact provides an io.Writer wrapper that masks secret values
// before they are written to an underlying writer (e.g. stdout/stderr).
package redact

import (
	"io"
	"sync"

	"github.com/your-org/vaultpipe/internal/mask"
)

// Writer wraps an io.Writer and redacts any registered secret values
// from bytes before forwarding them to the underlying writer.
type Writer struct {
	mu     sync.RWMutex
	w      io.Writer
	masker *mask.Masker
}

// New returns a Writer that redacts secrets known to m before writing
// to w. m must not be nil.
func New(w io.Writer, m *mask.Masker) *Writer {
	if m == nil {
		panic("redact: masker must not be nil")
	}
	return &Writer{w: w, masker: m}
}

// Write redacts secret values from p and writes the result to the
// underlying writer. The number of bytes reported as written always
// matches len(p) so callers behave correctly even though the
// underlying byte count may differ after redaction.
func (rw *Writer) Write(p []byte) (int, error) {
	rw.mu.RLock()
	redacted := rw.masker.Redact(string(p))
	rw.mu.RUnlock()

	if _, err := io.WriteString(rw.w, redacted); err != nil {
		return 0, err
	}
	return len(p), nil
}

// AddSecret registers an additional secret value to be redacted.
// Safe for concurrent use.
func (rw *Writer) AddSecret(secret string) {
	rw.mu.Lock()
	rw.masker.Add(secret)
	rw.mu.Unlock()
}
