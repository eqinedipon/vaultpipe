package process

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Runner executes a subprocess with an injected environment.
type Runner struct {
	env []string
}

// NewRunner creates a Runner with the provided environment slice.
func NewRunner(env []string) *Runner {
	return &Runner{env: env}
}

// Run executes the given command with args, forwarding signals and
// returning the process exit code.
func (r *Runner) Run(command string, args []string) (int, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = r.env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return 1, fmt.Errorf("starting process: %w", err)
	}

	// Forward OS signals to the child process.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sigCh {
			_ = cmd.Process.Signal(sig)
		}
	}()

	err := cmd.Wait()
	signal.Stop(sigCh)
	close(sigCh)

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus(), nil
			}
		}
		return 1, fmt.Errorf("process wait: %w", err)
	}
	return 0, nil
}
