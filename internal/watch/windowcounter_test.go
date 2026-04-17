package watch

import (
	"testing"
	"time"
)

func TestWindowCounterZeroOnInit(t *testing.T) {
	wc := NewWindowCounter(time.Second)
	if got := wc.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestWindowCounterDefaultsWindowOnZero(t *testing.T) {
	wc := NewWindowCounter(0)
	if wc.window != time.Minute {
		t.Fatalf("expected default window of 1m, got %v", wc.window)
	}
}

func TestWindowCounterAddIncrementsCount(t *testing.T) {
	wc := NewWindowCounter(time.Second)
	wc.Add()
	wc.Add()
	if got := wc.Count(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestWindowCounterResetClearsCount(t *testing.T) {
	wc := NewWindowCounter(time.Second)
	wc.Add()
	wc.Add()
	wc.Reset()
	if got := wc.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestWindowCounterEvictsExpiredEntries(t *testing.T) {
	wc := NewWindowCounter(50 * time.Millisecond)
	wc.Add()
	wc.Add()
	time.Sleep(80 * time.Millisecond)
	wc.Add()
	if got := wc.Count(); got != 1 {
		t.Fatalf("expected 1 after eviction, got %d", got)
	}
}

func TestWindowCounterCountDoesNotMutateOnEmpty(t *testing.T) {
	wc := NewWindowCounter(time.Second)
	// call Count multiple times on empty counter — should not panic
	for i := 0; i < 5; i++ {
		wc.Count()
	}
}
