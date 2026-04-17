package cache_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/cache"
)

func TestGet_MissWhenEmpty(t *testing.T) {
	c := cache.New(time.Minute)
	_, ok := c.Get("secret/foo")
	if ok {
		t.Fatal("expected cache miss on empty cache")
	}
}

func TestSet_ThenGet(t *testing.T) {
	c := cache.New(time.Minute)
	vals := map[string]string{"API_KEY": "abc123"}
	c.Set("secret/foo", vals)
	got, ok := c.Get("secret/foo")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got["API_KEY"] != "abc123" {
		t.Errorf("unexpected value: %s", got["API_KEY"])
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("secret/foo", map[string]string{"K": "v"})
	got, _ := c.Get("secret/foo")
	got["K"] = "mutated"
	again, _ := c.Get("secret/foo")
	if again["K"] == "mutated" {
		t.Error("cache returned reference, not copy")
	}
}

func TestGet_ExpiredEntry(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("secret/foo", map[string]string{"K": "v"})
	time.Sleep(20 * time.Millisecond)
	_, ok := c.Get("secret/foo")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestZeroTTL_DisablesCache(t *testing.T) {
	c := cache.New(0)
	c.Set("secret/foo", map[string]string{"K": "v"})
	_, ok := c.Get("secret/foo")
	if ok {
		t.Fatal("zero TTL should disable caching")
	}
}

func TestInvalidate(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("secret/foo", map[string]string{"K": "v"})
	c.Invalidate("secret/foo")
	_, ok := c.Get("secret/foo")
	if ok {
		t.Fatal("expected miss after invalidation")
	}
}

func TestFlush(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("a", map[string]string{"x": "1"})
	c.Set("b", map[string]string{"y": "2"})
	c.Flush()
	for _, k := range []string{"a", "b"} {
		if _, ok := c.Get(k); ok {
			t.Errorf("expected miss for %s after flush", k)
		}
	}
}
