package watch

import (
	"testing"
	"time"
)

func makeSuppressor(window time.Duration) *Suppressor {
	return NewSuppressor(SuppressPolicy{Window: window, MaxKeys: 16})
}

func TestDefaultSuppressPolicyValues(t *testing.T) {
	p := DefaultSuppressPolicy()
	if p.Window != 30*time.Second {
		t.Errorf("expected 30s window, got %v", p.Window)
	}
	if p.MaxKeys != 1024 {
		t.Errorf("expected 1024 max keys, got %d", p.MaxKeys)
	}
}

func TestSuppressorFirstCallAllowed(t *testing.T) {
	s := makeSuppressor(time.Second)
	if !s.Allow("port:tcp:8080") {
		t.Error("expected first call to be allowed")
	}
}

func TestSuppressorSecondCallSuppressed(t *testing.T) {
	s := makeSuppressor(time.Second)
	s.Allow("port:tcp:8080")
	if s.Allow("port:tcp:8080") {
		t.Error("expected second call within window to be suppressed")
	}
}

func TestSuppressorAllowedAfterWindowExpires(t *testing.T) {
	s := makeSuppressor(20 * time.Millisecond)
	s.Allow("key")
	time.Sleep(30 * time.Millisecond)
	if !s.Allow("key") {
		t.Error("expected call after window expiry to be allowed")
	}
}

func TestSuppressorDifferentKeysIndependent(t *testing.T) {
	s := makeSuppressor(time.Second)
	if !s.Allow("a") {
		t.Error("expected key a to be allowed")
	}
	if !s.Allow("b") {
		t.Error("expected key b to be allowed")
	}
}

func TestSuppressorCountIncrements(t *testing.T) {
	s := makeSuppressor(time.Second)
	s.Allow("k")
	s.Allow("k")
	s.Allow("k")
	if c := s.Count("k"); c != 3 {
		t.Errorf("expected count 3, got %d", c)
	}
}

func TestSuppressorCountZeroForUnknownKey(t *testing.T) {
	s := makeSuppressor(time.Second)
	if c := s.Count("missing"); c != 0 {
		t.Errorf("expected 0, got %d", c)
	}
}

func TestSuppressorResetAllowsImmediately(t *testing.T) {
	s := makeSuppressor(time.Second)
	s.Allow("k")
	s.Reset()
	if !s.Allow("k") {
		t.Error("expected allow after reset")
	}
}

func TestSuppressorMaxKeysEvicts(t *testing.T) {
	s := NewSuppressor(SuppressPolicy{Window: time.Second, MaxKeys: 2})
	s.Allow("a")
	s.Allow("b")
	// third key exceeds MaxKeys, triggers eviction
	if !s.Allow("c") {
		t.Error("expected c to be allowed after eviction")
	}
}
