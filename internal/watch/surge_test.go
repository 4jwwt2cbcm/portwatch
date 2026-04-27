package watch

import (
	"testing"
	"time"
)

func makeSurge(window time.Duration, threshold int, now func() time.Time) *Surge {
	return NewSurge(SurgePolicy{Window: window, Threshold: threshold}, now)
}

func TestDefaultSurgePolicyValues(t *testing.T) {
	p := DefaultSurgePolicy()
	if p.Window != 5*time.Second {
		t.Errorf("expected window 5s, got %v", p.Window)
	}
	if p.Threshold != 10 {
		t.Errorf("expected threshold 10, got %d", p.Threshold)
	}
}

func TestSurgeNoSurgeOnInit(t *testing.T) {
	s := makeSurge(time.Second, 3, nil)
	if s.Count() != 0 {
		t.Errorf("expected 0 events on init, got %d", s.Count())
	}
}

func TestSurgeRecordReturnsFalseBeforeThreshold(t *testing.T) {
	s := makeSurge(time.Second, 3, nil)
	if s.Record() {
		t.Error("expected no surge on first record")
	}
	if s.Record() {
		t.Error("expected no surge on second record")
	}
}

func TestSurgeRecordReturnsTrueAtThreshold(t *testing.T) {
	s := makeSurge(time.Second, 3, nil)
	s.Record()
	s.Record()
	if !s.Record() {
		t.Error("expected surge at threshold")
	}
}

func TestSurgeEvictsExpiredEvents(t *testing.T) {
	now := time.Now()
	clock := &now
	s := makeSurge(time.Second, 3, func() time.Time { return *clock })
	s.Record()
	s.Record()
	// advance past window
	future := now.Add(2 * time.Second)
	clock = &future
	if s.Count() != 0 {
		t.Errorf("expected 0 active events after window, got %d", s.Count())
	}
}

func TestSurgeResetClearsState(t *testing.T) {
	s := makeSurge(time.Second, 3, nil)
	s.Record()
	s.Record()
	s.Reset()
	if s.Count() != 0 {
		t.Errorf("expected 0 after reset, got %d", s.Count())
	}
}

func TestSurgeDefaultsOnZeroPolicy(t *testing.T) {
	s := NewSurge(SurgePolicy{}, nil)
	if s.policy.Window != DefaultSurgePolicy().Window {
		t.Errorf("expected default window, got %v", s.policy.Window)
	}
	if s.policy.Threshold != DefaultSurgePolicy().Threshold {
		t.Errorf("expected default threshold, got %d", s.policy.Threshold)
	}
}
