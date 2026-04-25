// Package env provides utilities for environment variable management.
package env

import (
	"context"
	"sync"
	"time"
)

// WatchFunc is called whenever a change is detected between two snapshots.
// changed contains only the keys that were added or modified.
type WatchFunc func(changed map[string]string)

// Watcher polls a snapshot source at a fixed interval and calls a callback
// when the environment changes relative to the previous snapshot.
type Watcher struct {
	mu       sync.Mutex
	last     *Snapshot
	interval time.Duration
	fetch    func() (*Snapshot, error)
	onChange WatchFunc
}

// WatcherOption configures a Watcher.
type WatcherOption func(*Watcher)

// WithInterval overrides the default polling interval (default: 30s).
func WithInterval(d time.Duration) WatcherOption {
	return func(w *Watcher) {
		w.interval = d
	}
}

// NewWatcher creates a Watcher that calls fetch to obtain fresh snapshots and
// invokes onChange whenever the result differs from the previous snapshot.
func NewWatcher(fetch func() (*Snapshot, error), onChange WatchFunc, opts ...WatcherOption) *Watcher {
	w := &Watcher{
		interval: 30 * time.Second,
		fetch:    fetch,
		onChange: onChange,
	}
	for _, o := range opts {
		o(w)
	}
	return w
}

// Run starts the polling loop. It blocks until ctx is cancelled.
// An initial fetch is performed immediately before the first tick.
func (w *Watcher) Run(ctx context.Context) error {
	if err := w.poll(); err != nil {
		return err
	}
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// best-effort; errors are silently ignored after initial fetch
			_ = w.poll()
		}
	}
}

// poll fetches a new snapshot and fires onChange if anything changed.
func (w *Watcher) poll() error {
	snap, err := w.fetch()
	if err != nil {
		return err
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.last == nil {
		w.last = snap
		return nil
	}
	diff := w.last.Diff(snap)
	if len(diff) > 0 {
		w.onChange(diff)
		w.last = snap
	}
	return nil
}
