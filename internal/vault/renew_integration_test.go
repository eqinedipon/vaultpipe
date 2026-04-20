package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// TestRenewer_StopsCleanlyOnContextCancel verifies that the renewer goroutine
// exits promptly when the parent context is cancelled, without leaking.
func TestRenewer_StopsCleanlyOnContextCancel(t *testing.T) {
	t.Parallel()

	var renewCount int64

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth/token/renew-self":
			atomic.AddInt64(&renewCount, 1)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"auth": map[string]any{
					"client_token": "test-token",
					"lease_duration": 10,
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := newTestVaultClient(t, srv.URL)
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Use a short TTL so the renewer would fire quickly if not cancelled.
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		_ = client.Renew(ctx, 2*time.Second)
	}()

	// Cancel almost immediately — renewer must stop.
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case <-doneCh:
		// success: renewer exited
	case <-time.After(2 * time.Second):
		t.Fatal("renewer did not stop after context cancellation")
	}
}

// TestRenewer_MultipleRenewCycles checks that the renewer fires more than once
// when the TTL is very short, confirming the loop behaviour.
func TestRenewer_MultipleRenewCycles(t *testing.T) {
	t.Parallel()

	var renewCount int64

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/token/renew-self" {
			http.NotFound(w, r)
			return
		}
		atomic.AddInt64(&renewCount, 1)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"auth": map[string]any{
				"client_token": "test-token",
				// Return a 2-second TTL; renewer should fire at ~1 s intervals.
				"lease_duration": 2,
			},
		})
	}))
	t.Cleanup(srv.Close)

	client, err := newTestVaultClient(t, srv.URL)
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3500*time.Millisecond)
	t.Cleanup(cancel)

	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		_ = client.Renew(ctx, 2*time.Second)
	}()

	<-doneCh

	got := atomic.LoadInt64(&renewCount)
	if got < 2 {
		t.Errorf("expected at least 2 renew calls, got %d", got)
	}
}
