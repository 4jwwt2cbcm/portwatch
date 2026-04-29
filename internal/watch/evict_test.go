package watch

import (
	"testing"
	"time"
)

func makeEvictCache[T any](ttl time.Duration, cap int) *EvictCache[T] {
	return NewEvictCache[T](EvictPolicy{TTL: ttl, Capacity: cap})
}

func TestDefaultEvictPolicyValues(t *testing.T) {
	p := DefaultEvictPolicy()
	if p.TTL != 30*time.Second {
		t.Fatalf("expected 30s TTL, got %v", p.TTL)
	}
	if p.Capacity != 256 {
		t.Fatalf("expected capacity 256, got %d", p.Capacity)
	}
}

func TestEvictCacheSetAndGet(t *testing.T) {
	c := makeEvictCache[string](time.Minute, 10)
	c.Set("key", "value")
	v, ok := c.Get("key")
	if !ok || v != "value" {
		t.Fatalf("expected value 'value', got %q ok=%v", v, ok)
	}
}

func TestEvictCacheMissReturnsZero(t *testing.T) {
	c := makeEvictCache[int](time.Minute, 10)
	v, ok := c.Get("missing")
	if ok || v != 0 {
		t.Fatalf("expected zero and false, got %d %v", v, ok)
	}
}

func TestEvictCacheExpiredEntryReturnsZero(t *testing.T) {
	c := makeEvictCache[string](time.Millisecond, 10)
	c.Set("k", "v")
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("k")
	if ok {
		t.Fatal("expected expired entry to return false")
	}
}

func TestEvictCacheLenExcludesExpired(t *testing.T) {
	c := makeEvictCache[string](time.Millisecond, 10)
	c.Set("a", "1")
	c.Set("b", "2")
	time.Sleep(5 * time.Millisecond)
	if c.Len() != 0 {
		t.Fatalf("expected 0 after expiry, got %d", c.Len())
	}
}

func TestEvictCacheEvictsOldestWhenFull(t *testing.T) {
	c := makeEvictCache[int](time.Minute, 3)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	c.Set("d", 4) // should evict "a"
	_, ok := c.Get("a")
	if ok {
		t.Fatal("expected 'a' to be evicted")
	}
	_, ok = c.Get("d")
	if !ok {
		t.Fatal("expected 'd' to be present")
	}
}

func TestEvictCacheDefaultsOnZero(t *testing.T) {
	c := NewEvictCache[string](EvictPolicy{})
	if c.policy.TTL != DefaultEvictPolicy().TTL {
		t.Fatalf("expected default TTL")
	}
	if c.policy.Capacity != DefaultEvictPolicy().Capacity {
		t.Fatalf("expected default capacity")
	}
}
