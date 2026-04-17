package vault_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/cache"
	"github.com/yourusername/vaultpipe/internal/vault"
)

func TestCachedClient_HitAvoidsFetch(t *testing.T) {
	var calls atomic.Int32
	srv := mockKVv2Server(t, map[string]interface{}{"TOKEN": "s3cr3t"}, &calls)

	client, err := vault.NewClient(vault.Config{Address: srv.URL, Token: "test"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	ca := cache.New(time.Minute)
	cc := vault.NewCachedClient(client, ca)

	for i := 0; i < 3; i++ {
		vals, err := cc.GetSecrets("secret", "myapp")
		if err != nil {
			t.Fatalf("GetSecrets call %d: %v", i, err)
		}
		if vals["TOKEN"] != "s3cr3t" {
			t.Errorf("unexpected value: %s", vals["TOKEN"])
		}
	}
	if calls.Load() != 1 {
		t.Errorf("expected 1 upstream call, got %d", calls.Load())
	}
}

func TestCachedClient_ZeroTTL_AlwaysFetches(t *testing.T) {
	var calls atomic.Int32
	srv := mockKVv2Server(t, map[string]interface{}{"TOKEN": "s3cr3t"}, &calls)

	client, err := vault.NewClient(vault.Config{Address: srv.URL, Token: "test"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	ca := cache.New(0)
	cc := vault.NewCachedClient(client, ca)

	for i := 0; i < 3; i++ {
		if _, err := cc.GetSecrets("secret", "myapp"); err != nil {
			t.Fatalf("GetSecrets: %v", err)
		}
	}
	if calls.Load() != 3 {
		t.Errorf("expected 3 upstream calls, got %d", calls.Load())
	}
}

func TestCachedClient_InvalidatePath(t *testing.T) {
	var calls atomic.Int32
	srv := mockKVv2Server(t, map[string]interface{}{"K": "v"}, &calls)

	client, err := vault.NewClient(vault.Config{Address: srv.URL, Token: "test"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	ca := cache.New(time.Minute)
	cc := vault.NewCachedClient(client, ca)

	if _, err := cc.GetSecrets("secret", "myapp"); err != nil {
		t.Fatal(err)
	}
	cc.InvalidatePath("secret", "myapp")
	if _, err := cc.GetSecrets("secret", "myapp"); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 2 {
		t.Errorf("expected 2 upstream calls after invalidation, got %d", calls.Load())
	}
}
