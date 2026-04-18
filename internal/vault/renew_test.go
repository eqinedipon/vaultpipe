package vault

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

func newRenewMockServer(t *testing.T, ttlSeconds int, renewCalls *int32) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth/token/lookup-self":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"ttl": float64(ttlSeconds)},
			})
		case "/v1/auth/token/renew-self":
			atomic.AddInt32(renewCalls, 1)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"auth": map[string]interface{}{}})
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestRenewer_RenewsBeforeExpiry(t *testing.T) {
	var calls int32
	srv := newRenewMockServer(t, 2, &calls) // 2s TTL → renew after ~1s
	defer srv.Close()

	cfg := vaultapi.DefaultConfig()
	cfg.Address = srv.URL
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	client.SetToken("test-token")

	logger := log.New(os.Stderr, "", 0)
	r := NewRenewer(client, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	r.Start(ctx)
	<-ctx.Done()

	if atomic.LoadInt32(&calls) < 1 {
		t.Errorf("expected at least 1 renew call, got %d", calls)
	}
}

func TestRenewer_ZeroTTL_NoRenew(t *testing.T) {
	var calls int32
	srv := newRenewMockServer(t, 0, &calls)
	defer srv.Close()

	cfg := vaultapi.DefaultConfig()
	cfg.Address = srv.URL
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	client.SetToken("root")

	logger := log.New(os.Stderr, "", 0)
	r := NewRenewer(client, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	r.Start(ctx)
	<-ctx.Done()

	if atomic.LoadInt32(&calls) != 0 {
		t.Errorf("expected 0 renew calls for non-expiring token, got %d", calls)
	}
}
