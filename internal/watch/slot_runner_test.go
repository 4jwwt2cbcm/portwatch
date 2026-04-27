package watch

import (
	"errors"
	"testing"
	"time"
)

func makeSlotRunner(max int) *SlotRunner {
	s := NewSlot(SlotPolicy{Duration: time.Second, MaxSlots: 10})
	return NewSlotRunner(s, max)
}

func TestSlotRunnerAllowsUpToMax(t *testing.T) {
	r := makeSlotRunner(3)
	for i := 0; i < 3; i++ {
		if err := r.Run(nil); err != nil {
			t.Fatalf("expected nil on call %d, got %v", i+1, err)
		}
	}
}

func TestSlotRunnerBlocksWhenExhausted(t *testing.T) {
	r := makeSlotRunner(2)
	r.Run(nil)
	r.Run(nil)
	if err := r.Run(nil); !errors.Is(err, ErrSlotFull) {
		t.Fatalf("expected ErrSlotFull, got %v", err)
	}
}

func TestSlotRunnerPropagatesError(t *testing.T) {
	r := makeSlotRunner(5)
	sentinel := errors.New("boom")
	if err := r.Run(func() error { return sentinel }); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel, got %v", err)
	}
}

func TestSlotRunnerNilArgDefaults(t *testing.T) {
	r := NewSlotRunner(nil, 0)
	if r.slot == nil {
		t.Fatal("expected non-nil slot")
	}
	if r.max != 1 {
		t.Fatalf("expected max 1, got %d", r.max)
	}
}

func TestSlotRunnerNilFnSucceeds(t *testing.T) {
	r := makeSlotRunner(10)
	if err := r.Run(nil); err != nil {
		t.Fatalf("expected nil for nil fn, got %v", err)
	}
}
