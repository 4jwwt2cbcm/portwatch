package watch

import (
	"testing"
	"time"
)

func makeShedder() *Shedder {
	return NewShedder(ShedderPolicy{
		MaxLoad:  0.75,
		Window:   500 * time.Millisecond,
		Cooldown: 100 * time.Millisecond,
	})
}

func TestDefaultShedderPolicyValues(t *testing.T) {
	p := DefaultShedderPolicy()
	if p.MaxLoad != 0.80 {
		t.Errorf("MaxLoad = %v, want 0.80", p.MaxLoad)
	}
	if p.Window != 5*time.Second {
		t.Errorf("Window = %v, want 5s", p.Window)
	}
	if p.Cooldown != 2*time.Second {
		t.Errorf("Cooldown = %v, want 2s", p.Cooldown)
	}
}

func TestShedderLoadZeroOnInit(t *testing.T) {
	s := makeShedder()
	if s.Load() != 0 {
		t.Errorf("expected zero load, got %v", s.Load())
	}
}

func TestShedderNoShedBelowThreshold(t *testing.T) {
	s := makeShedder()
	s.Record(0.5)
	s.Record(0.5)
	if s.ShouldShed() {
		t.Error("expected no shed below threshold")
	}
}

func TestShedderShedsAboveThreshold(t *testing.T) {
	s := makeShedder()
	s.Record(0.9)
	s.Record(0.9)
	if !s.ShouldShed() {
		t.Error("expected shed above threshold")
	}
}

func TestShedderCooldownSuppressesRepeatShed(t *testing.T) {
	s := makeShedder()
	s.Record(0.9)
	s.Record(0.9)
	if !s.ShouldShed() {
		t.Fatal("expected first shed")
	}
	// Second call within cooldown should be suppressed.
	s.Record(0.9)
	if s.ShouldShed() {
		t.Error("expected cooldown to suppress second shed")
	}
}

func TestShedderAllowsAfterCooldown(t *testing.T) {
	s := NewShedder(ShedderPolicy{
		MaxLoad:  0.75,
		Window:   500 * time.Millisecond,
		Cooldown: 10 * time.Millisecond,
	})
	s.Record(0.9)
	if !s.ShouldShed() {
		t.Fatal("expected first shed")
	}
	time.Sleep(20 * time.Millisecond)
	s.Record(0.9)
	if !s.ShouldShed() {
		t.Error("expected shed allowed after cooldown")
	}
}

func TestShedderEvictsOldSamples(t *testing.T) {
	s := NewShedder(ShedderPolicy{
		MaxLoad:  0.75,
		Window:   20 * time.Millisecond,
		Cooldown: 1 * time.Millisecond,
	})
	s.Record(0.9)
	time.Sleep(30 * time.Millisecond)
	// Old sample evicted; load should be zero.
	if s.Load() != 0 {
		t.Errorf("expected zero load after eviction, got %v", s.Load())
	}
	if s.ShouldShed() {
		t.Error("expected no shed after samples evicted")
	}
}
