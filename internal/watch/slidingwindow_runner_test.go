package watch

import (
	"errors"
	"testing"
	"time"
)

func TestSlidingWindowRunnerRunsFunction(t *testing.T) {
	called := false
	sw := makeSlidingWindow(5, time.Minute)
	r := NewSlidingWindowRunner(sw, func() error {
		called = true
		return nil
	})
	if err := r.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected fn to be called")
	}
}

func TestSlidingWindowRunnerBlocksWhenFull(t *testing.T) {
	sw := makeSlidingWindow(2, time.Minute)
	r := NewSlidingWindowRunner(sw, func() error { return nil })
	// first two calls fill window (size=2 means record returns true on 2nd)
	r.Run() // count=1, not full yet
	err := r.Run() // count=2, full — returns ErrWindowFull
	if !errors.Is(err, ErrWindowFull) {
		t.Errorf("expected ErrWindowFull, got %v", err)
	}
}

func TestSlidingWindowRunnerPropagatesError(t *testing.T) {
	sentinel := errors.New("fn error")
	sw := makeSlidingWindow(5, time.Minute)
	r := NewSlidingWindowRunner(sw, func() error { return sentinel })
	if err := r.Run(); !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestSlidingWindowRunnerNilArgDefaults(t *testing.T) {
	r := NewSlidingWindowRunner(nil, nil)
	if r.win == nil {
		t.Error("expected non-nil window")
	}
	if err := r.Run(); err != nil {
		t.Errorf("unexpected error with nil fn: %v", err)
	}
}
