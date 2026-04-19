package watch

import (
	"testing"
	"time"
)

func makeRotator(max int, age time.Duration) *Rotator[string] {
	r := NewRotator[string](RotatorPolicy{MaxEntries: max, MaxAge: age})
	return r
}

func TestDefaultRotatorPolicyValues(t *testing.T) {
	p := DefaultRotatorPolicy()
	if p.MaxEntries != 1000 {
		t.Fatalf("expected 1000, got %d", p.MaxEntries)
	}
	if p.MaxAge != 24*time.Hour {
		t.Fatalf("expected 24h, got %v", p.MaxAge)
	}
}

func TestRotatorEmptyOnInit(t *testing.T) {
	r := makeRotator(10, time.Minute)
	if r.Len() != 0 {
		t.Fatal("expected empty rotator")
	}
}

func TestRotatorAddAndSnapshot(t *testing.T) {
	r := makeRotator(10, time.Minute)
	r.Add("a")
	r.Add("b")
	snap := r.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
	if snap[0].Value != "a" || snap[1].Value != "b" {
		t.Fatal("unexpected values in snapshot")
	}
}

func TestRotatorEvictsWhenFull(t *testing.T) {
	r := makeRotator(3, time.Hour)
	r.Add("a")
	r.Add("b")
	r.Add("c")
	r.Add("d")
	snap := r.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(snap))
	}
	if snap[0].Value != "b" {
		t.Fatalf("expected oldest evicted, got %s", snap[0].Value)
	}
}

func TestRotatorEvictsStaleByAge(t *testing.T) {
	r := makeRotator(100, 50*time.Millisecond)
	fixedOld := time.Now().Add(-100 * time.Millisecond)
	r.mu.Lock()
	r.entries = append(r.entries, RotatorEntry[string]{Value: "old", Timestamp: fixedOld})
	r.mu.Unlock()
	r.Add("new")
	snap := r.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry after eviction, got %d", len(snap))
	}
	if snap[0].Value != "new" {
		t.Fatal("expected 'new' to remain")
	}
}

func TestRotatorClear(t *testing.T) {
	r := makeRotator(10, time.Minute)
	r.Add("x")
	r.Clear()
	if r.Len() != 0 {
		t.Fatal("expected empty after clear")
	}
}

func TestRotatorDefaultsOnZeroPolicy(t *testing.T) {
	r := NewRotator[int](RotatorPolicy{})
	if r.policy.MaxEntries != 1000 {
		t.Fatalf("expected default MaxEntries, got %d", r.policy.MaxEntries)
	}
}
