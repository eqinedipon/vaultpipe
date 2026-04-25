package env

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestWatcher_CallsOnChange_WhenSnapshotDiffers(t *testing.T) {
	calls := 0
	var mu sync.Mutex

	generation := 0
	fetch := func() (*Snapshot, error) {
		mu.Lock()
		defer mu.Unlock()
		generation++
		if generation == 1 {
			return NewSnapshotFromMap(map[string]string{"KEY": "v1"}), nil
		}
		return NewSnapshotFromMap(map[string]string{"KEY": "v2"}), nil
	}

	var changed map[string]string
	w := NewWatcher(fetch, func(c map[string]string) {
		mu.Lock()
		defer mu.Unlock()
		calls++
		changed = c
	}, WithInterval(10*time.Millisecond))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	_ = w.Run(ctx)

	mu.Lock()
	defer mu.Unlock()
	if calls == 0 {
		t.Fatal("expected onChange to be called at least once")
	}
	if changed["KEY"] != "v2" {
		t.Errorf("expected changed KEY=v2, got %q", changed["KEY"])
	}
}

func TestWatcher_NoChange_OnChangeNotCalled(t *testing.T) {
	calls := 0
	var mu sync.Mutex

	fetch := func() (*Snapshot, error) {
		return NewSnapshotFromMap(map[string]string{"KEY": "stable"}), nil
	}

	w := NewWatcher(fetch, func(_ map[string]string) {
		mu.Lock()
		defer mu.Unlock()
		calls++
	}, WithInterval(10*time.Millisecond))

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_ = w.Run(ctx)

	mu.Lock()
	defer mu.Unlock()
	if calls != 0 {
		t.Errorf("expected no onChange calls, got %d", calls)
	}
}

func TestWatcher_FetchError_ReturnsError(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	fetch := func() (*Snapshot, error) {
		return nil, sentinel
	}

	w := NewWatcher(fetch, func(_ map[string]string) {})
	err := w.Run(context.Background())
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestWatcher_ContextCancel_StopsLoop(t *testing.T) {
	fetch := func() (*Snapshot, error) {
		return NewSnapshotFromMap(map[string]string{}), nil
	}
	w := NewWatcher(fetch, func(_ map[string]string) {}, WithInterval(5*time.Millisecond))

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- w.Run(ctx) }()

	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("watcher did not stop after context cancel")
	}
}
