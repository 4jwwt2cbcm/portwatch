package watch

import (
	"testing"
	"time"
)

func makeHoldOff(d time.Duration) *HoldOff {
	h := NewHoldOff(HoldOffPolicy{Duration: d})
	return h
}

func TestDefaultHoldOffPolicyValues(t *testing.T) {
	p := DefaultHoldOffPolicy()
	if p.Duration != 5*time.Second {
		t.Fatalf("expected 5s, got %v", p.Duration)
	}
}

func TestHoldOffClearBeforeSignal(t *testing.T) {
	h := makeHoldOff(100 * time.Millisecond)
	if !h.Clear() {
		t.Fatal("expected clear before any signal")
	}
}

func TestHoldOffNotClearImmediatelyAfterSignal(t *testing.T) {
	h := makeHoldOff(100 * time.Millisecond)
	h.Signal()
	if h.Clear() {
		t.Fatal("expected not clear immediately after signal")
	}
}

func TestHoldOffClearAfterDurationElapsed(t *testing.T) {
	h := makeHoldOff(10 * time.Millisecond)
	h.Signal()
	time.Sleep(20 * time.Millisecond)
	if !h.Clear() {
		t.Fatal("expected clear after hold-off duration elapsed")
	}
}

func TestHoldOffSignalResetsQuietPeriod(t *testing.T) {
	h := makeHoldOff(50 * time.Millisecond)
	h.Signal()
	time.Sleep(30 * time.Millisecond)
	h.Signal() // re-arm
	if h.Clear() {
		t.Fatal("expected not clear; quiet period should have reset")
	}
}

func TestHoldOffResetAllowsClearImmediately(t *testing.T) {
	h := makeHoldOff(100 * time.Millisecond)
	h.Signal()
	h.Reset()
	if !h.Clear() {
		t.Fatal("expected clear after Reset")
	}
}

func TestHoldOffLastSeenZeroBeforeSignal(t *testing.T) {
	h := makeHoldOff(100 * time.Millisecond)
	if !h.LastSeen().IsZero() {
		t.Fatal("expected zero LastSeen before any signal")
	}
}

func TestHoldOffLastSeenUpdatedOnSignal(t *testing.T) {
	h := makeHoldOff(100 * time.Millisecond)
	before := time.Now()
	h.Signal()
	after := time.Now()
	ls := h.LastSeen()
	if ls.Before(before) || ls.After(after) {
		t.Fatalf("LastSeen %v not between %v and %v", ls, before, after)
	}
}

func TestHoldOffDefaultsOnZeroDuration(t *testing.T) {
	h := NewHoldOff(HoldOffPolicy{Duration: 0})
	if h.policy.Duration != DefaultHoldOffPolicy().Duration {
		t.Fatalf("expected default duration, got %v", h.policy.Duration)
	}
}
