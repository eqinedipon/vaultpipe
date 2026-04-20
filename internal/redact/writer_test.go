package redact_test

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/your-org/vaultpipe/internal/mask"
	"github.com/your-org/vaultpipe/internal/redact"
)

func TestWrite_RedactsSecret(t *testing.T) {
	m := mask.New([]string{"s3cr3t"})
	var buf bytes.Buffer
	w := redact.New(&buf, m)

	_, err := fmt.Fprint(w, "the password is s3cr3t ok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "s3cr3t") {
		t.Errorf("secret leaked in output: %q", buf.String())
	}
	if !strings.Contains(buf.String(), "[REDACTED]") {
		t.Errorf("expected [REDACTED] in output, got: %q", buf.String())
	}
}

func TestWrite_NoSecret_PassThrough(t *testing.T) {
	m := mask.New([]string{"s3cr3t"})
	var buf bytes.Buffer
	w := redact.New(&buf, m)

	_, err := fmt.Fprint(w, "nothing sensitive here")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "nothing sensitive here" {
		t.Errorf("unexpected mutation: %q", buf.String())
	}
}

func TestWrite_ReportsOriginalLen(t *testing.T) {
	m := mask.New([]string{"s3cr3t"})
	var buf bytes.Buffer
	w := redact.New(&buf, m)

	input := "s3cr3t"
	n, err := w.Write([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(input) {
		t.Errorf("expected n=%d, got %d", len(input), n)
	}
}

func TestAddSecret_RedactsAfterAdd(t *testing.T) {
	m := mask.New(nil)
	var buf bytes.Buffer
	w := redact.New(&buf, m)

	w.AddSecret("lateSecret")
	fmt.Fprint(w, "value=lateSecret")

	if strings.Contains(buf.String(), "lateSecret") {
		t.Errorf("late-added secret leaked: %q", buf.String())
	}
}

func TestWrite_ConcurrentSafe(t *testing.T) {
	m := mask.New([]string{"tok3n"})
	var buf syncBuffer
	w := redact.New(&buf, m)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Fprint(w, "tok3n")
		}()
	}
	wg.Wait()

	if strings.Contains(buf.String(), "tok3n") {
		t.Error("secret leaked under concurrent writes")
	}
}

type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (sb *syncBuffer) Write(p []byte) (int, error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.Write(p)
}

func (sb *syncBuffer) String() string {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.String()
}
