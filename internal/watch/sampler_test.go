package watch

import (
	"testing"
	"time"
)

func makeSampler(rate int, window time.Duration) *Sampler[string] {
	return NewSampler[string](SamplerPolicy{Rate: rate, Window: window})
}

func TestDefaultSamplerPolicyValues(t *testing.T) {
	p := DefaultSamplerPolicy()
	if p.Rate != 10 {
		t.Fatalf("expected rate 10, got %d", p.Rate)
	}
	if p.Window != time.Minute {
		t.Fatalf("expected window 1m, got %v", p.Window)
	}
}

func TestSamplerDefaultsOnZero(t *testing.T) {
	s := NewSampler[int](SamplerPolicy{})
	if s.policy.Rate != 10 {
		t.Fatalf("expected default rate 10, got %d", s.policy.Rate)
	}
}

func TestSamplerRecordAccepts(t *testing.T) {
	s := makeSampler(3, time.Minute)
	if !s.Record("a") {
		t.Fatal("expected first record to be accepted")
	}
}

func TestSamplerRateLimitRejects(t *testing.T) {
	s := makeSampler(2, time.Minute)
	s.Record("a")
	s.Record("b")
	if s.Record("c") {
		t.Fatal("expected third record to be rejected")
	}
}

func TestSamplerSnapshotReturnsValues(t *testing.T) {
	s := makeSampler(5, time.Minute)
	s.Record("x")
	s.Record("y")
	snap := s.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 samples, got %d", len(snap))
	}
	if snap[0] != "x" || snap[1] != "y" {
		t.Fatalf("unexpected snapshot values: %v", snap)
	}
}

func TestSamplerEvictsExpiredEntries(t *testing.T) {
	s := makeSampler(10, 50*time.Millisecond)
	s.Record("old")
	time.Sleep(80 * time.Millisecond)
	s.Record("new")
	snap := s.Snapshot()
	if len(snap) != 1 || snap[0] != "new" {
		t.Fatalf("expected only 'new', got %v", snap)
	}
}

func TestSamplerLenReflectsWindow(t *testing.T) {
	s := makeSampler(10, 50*time.Millisecond)
	s.Record("a")
	s.Record("b")
	if s.Len() != 2 {
		t.Fatalf("expected len 2, got %d", s.Len())
	}
	time.Sleep(80 * time.Millisecond)
	if s.Len() != 0 {
		t.Fatalf("expected len 0 after expiry, got %d", s.Len())
	}
}
