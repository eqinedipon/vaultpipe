package signal_test

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/signal"
)

// fakeProcess wraps a real *os.Process obtained from a no-op command so we
// can introspect signals without spawning a real subprocess in every test.
func startSleepProc(t *testing.T) *os.Process {
	t.Helper()
	cmd := testCmd()
	if err := cmd.Start(); err != nil {
		t.Fatalf("start helper process: %v", err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })
	return cmd.Process
}

func TestForwarder_StartStop(t *testing.T) {
	proc := startSleepProc(t)
	f := signal.New(proc)
	f.Start()
	// Stop should not block indefinitely.
	done := make(chan struct{})
	go func() { f.Stop(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Stop() timed out")
	}
}

func TestForwarder_ForwardsSignal(t *testing.T) {
	proc := startSleepProc(t)

	// Use SIGUSR1 so we don't kill the process.
	f := signal.New(proc, syscall.SIGUSR1)
	f.Start()
	defer f.Stop()

	// Send SIGUSR1 to ourselves; the forwarder should relay it to proc.
	// We just verify no panic / deadlock occurs within the timeout.
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(50 * time.Millisecond)
}

func TestForwarder_DefaultSignals(t *testing.T) {
	proc := startSleepProc(t)
	// Constructing with no explicit signals should not panic.
	f := signal.New(proc)
	f.Start()
	f.Stop()
}
