// Package signal provides OS signal forwarding to child processes.
// It ensures that signals received by vaultpipe are propagated to the
// managed subprocess, enabling graceful shutdown and job control.
package signal

import (
	"os"
	"os/signal"
	"syscall"
)

// Forwarder listens for OS signals and forwards them to a target process.
type Forwarder struct {
	proc   *os.Process
	signals []os.Signal
	ch      chan os.Signal
	done    chan struct{}
}

// New creates a Forwarder that will relay the given signals to proc.
// If signals is empty, a default set is used (SIGINT, SIGTERM, SIGHUP).
func New(proc *os.Process, signals ...os.Signal) *Forwarder {
	if len(signals) == 0 {
		signals = []os.Signal{
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGHUP,
		}
	}
	return &Forwarder{
		proc:    proc,
		signals: signals,
		ch:      make(chan os.Signal, 8),
		done:    make(chan struct{}),
	}
}

// Start begins forwarding signals in a background goroutine.
// Call Stop to unregister and clean up.
func (f *Forwarder) Start() {
	signal.Notify(f.ch, f.signals...)
	go func() {
		defer close(f.done)
		for {
			select {
			case sig, ok := <-f.ch:
				if !ok {
					return
				}
				// Best-effort: ignore send errors (process may have already exited).
				_ = f.proc.Signal(sig)
			}
		}
	}()
}

// Stop unregisters signal notifications and waits for the forwarding
// goroutine to exit.
func (f *Forwarder) Stop() {
	signal.Stop(f.ch)
	close(f.ch)
	<-f.done
}
