package watch

import (
	"testing"
	"time"
)

func makeShedCounter(window time.Duration) *ShedCounter {
	return NewShedCounter(ShedCounterPolicy{
		Window:   window,
		Capacity: 100,
	})
}

func TestDefaultShedCounterPolicyValues(t *testing.T) {
	p := DefaultShedCounterPolicy()
	if p.Window != 30*time.Second {
		t.Errorf("expected 30s window, got %v", p.Window)
	}
	if p.Capacity != 1000 {
		t.Errorf("expected capacity 1000, got %d", p.Capacity)
	}
}

func TestShedCounterZeroOnInit(t *testing.T) {
	s := makeShedCounter(time.Minute)
	if got := s.Count(); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestShedCounterRecordIncrements(t *testing.T) {
	s := makeShedCounter(time.Minute)
	s.Record()
	s.Record()
	if got := s.Count(); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}

func TestShedCounterEvictsExpiredEntries(t *testing.T) {
	s := makeShedCounter(50 * time.Millisecond)
	s.Record()
	s.Record()
	time.Sleep(80 * time.Millisecond)
	if got := s.Count(); got != 0 {
		t.Errorf("expected 0 after window expiry, got %d", got)
	}
}

func TestShedCounterResetClearsEntries(t *testing.T) {
	s := makeShedCounter(time.Minute)
	s.Record()
	s.Record()
	s.Reset()
	if got := s.Count(); got != 0 {
		t.Errorf("expected 0 after reset, got %d", got)
	}
}

func TestShedCounterRespectsCapacity(t *testing.T) {
	s := NewShedCounter(ShedCounterPolicy{Window: time.Minute, Capacity: 3})
	for i := 0; i < 10; i++ {
		s.Record()
	}
	if got := s.Count(); got != 3 {
		t.Errorf("expected count capped at 3, got %d", got)
	}
}

func TestShedCounterDefaultsOnZeroPolicy(t *testing.T) {
	s := NewShedCounter(ShedCounterPolicy{})
	if s.policy.Window != 30*time.Second {
		t.Errorf("expected default window, got %v", s.policy.Window)
	}
	if s.policy.Capacity != 1000 {
		t.Errorf("expected default capacity, got %d", s.policy.Capacity)
	}
}
