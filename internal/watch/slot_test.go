package watch

import (
	"testing"
	"time"
)

func makeSlot() *Slot {
	return NewSlot(SlotPolicy{
		Duration: time.Second,
		MaxSlots: 5,
	})
}

func TestDefaultSlotPolicyValues(t *testing.T) {
	p := DefaultSlotPolicy()
	if p.Duration != time.Minute {
		t.Fatalf("expected 1m, got %v", p.Duration)
	}
	if p.MaxSlots != 60 {
		t.Fatalf("expected 60, got %d", p.MaxSlots)
	}
}

func TestSlotDefaultsOnZero(t *testing.T) {
	s := NewSlot(SlotPolicy{})
	if s.policy.Duration != time.Minute {
		t.Fatalf("expected default duration")
	}
	if s.policy.MaxSlots != 60 {
		t.Fatalf("expected default max slots")
	}
}

func TestSlotCountZeroOnInit(t *testing.T) {
	s := makeSlot()
	if s.Count(time.Now()) != 0 {
		t.Fatal("expected zero count on init")
	}
}

func TestSlotRecordIncrements(t *testing.T) {
	s := makeSlot()
	now := time.Now()
	s.Record(now)
	s.Record(now)
	if got := s.Count(now); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestSlotDifferentSlotsAreIndependent(t *testing.T) {
	s := makeSlot()
	t1 := time.Now()
	t2 := t1.Add(2 * time.Second)
	s.Record(t1)
	s.Record(t2)
	s.Record(t2)
	if s.Count(t1) != 1 {
		t.Fatalf("expected 1 in slot1, got %d", s.Count(t1))
	}
	if s.Count(t2) != 2 {
		t.Fatalf("expected 2 in slot2, got %d", s.Count(t2))
	}
}

func TestSlotResetClearsAll(t *testing.T) {
	s := makeSlot()
	now := time.Now()
	s.Record(now)
	s.Reset()
	if s.Count(now) != 0 {
		t.Fatal("expected zero after reset")
	}
}

func TestSlotEvictsOldestWhenFull(t *testing.T) {
	s := NewSlot(SlotPolicy{Duration: time.Second, MaxSlots: 3})
	base := time.Now()
	// fill 3 slots
	for i := 0; i < 3; i++ {
		s.Record(base.Add(time.Duration(i) * time.Second))
	}
	// adding a 4th should evict the oldest
	s.Record(base.Add(3 * time.Second))
	if len(s.buckets) > 3 {
		t.Fatalf("expected at most 3 buckets, got %d", len(s.buckets))
	}
}
