package watch

import (
	"testing"
	"time"
)

func makeStagger() *Stagger {
	return NewStagger(DefaultStaggerPolicy())
}

func TestDefaultStaggerPolicyValues(t *testing.T) {
	p := DefaultStaggerPolicy()
	if p.Delay != 100*time.Millisecond {
		t.Errorf("expected 100ms delay, got %v", p.Delay)
	}
	if p.MaxItems != 10 {
		t.Errorf("expected MaxItems=10, got %d", p.MaxItems)
	}
}

func TestStaggerDefaultsOnZero(t *testing.T) {
	s := NewStagger(StaggerPolicy{})
	if s.policy.Delay <= 0 {
		t.Error("expected non-zero delay after defaulting")
	}
	if s.policy.MaxItems <= 0 {
		t.Error("expected non-zero MaxItems after defaulting")
	}
}

func TestStaggerCountZeroOnInit(t *testing.T) {
	s := makeStagger()
	if s.Count() != 0 {
		t.Errorf("expected count=0, got %d", s.Count())
	}
}

func TestStaggerNextIncrementsCount(t *testing.T) {
	s := makeStagger()
	s.Next()
	s.Next()
	s.Next()
	if s.Count() != 3 {
		t.Errorf("expected count=3, got %d", s.Count())
	}
}

func TestStaggerNextSlotsAreOrdered(t *testing.T) {
	s := NewStagger(StaggerPolicy{Delay: 10 * time.Millisecond, MaxItems: 5})
	var slots []time.Time
	for i := 0; i < 4; i++ {
		slots = append(slots, s.Next())
	}
	for i := 1; i < len(slots); i++ {
		if !slots[i].After(slots[i-1]) {
			t.Errorf("slot[%d] (%v) not after slot[%d] (%v)", i, slots[i], i-1, slots[i-1])
		}
	}
}

func TestStaggerResetClearsCount(t *testing.T) {
	s := makeStagger()
	s.Next()
	s.Next()
	s.Reset()
	if s.Count() != 0 {
		t.Errorf("expected count=0 after reset, got %d", s.Count())
	}
}

func TestStaggerResetAllowsReuse(t *testing.T) {
	s := NewStagger(StaggerPolicy{Delay: 5 * time.Millisecond, MaxItems: 5})
	first := s.Next()
	s.Reset()
	second := s.Next()
	if second.Before(first) {
		t.Errorf("second slot (%v) should not precede first (%v) after reset", second, first)
	}
}
