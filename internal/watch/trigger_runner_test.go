package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestTriggerRunnerRunOnceExecutesOnFire(t *testing.T) {
	tr := makeTrigger(0)
	called := false
	r := NewTriggerRunner(tr, func(_ context.Context) error {
		called = true
		return nil
	})
	tr.Fire()
	if err := r.RunOnce(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected fn to be called")
	}
}

func TestTriggerRunnerRunOnceContextCancelReturnsErr(t *testing.T) {
	tr := makeTrigger(0)
	r := NewTriggerRunner(tr, func(_ context.Context) error { return nil })
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := r.RunOnce(ctx); !errors.Is(err, ErrTriggerNotFired) {
		t.Fatalf("expected ErrTriggerNotFired, got %v", err)
	}
}

func TestTriggerRunnerRunOncePropagatesError(t *testing.T) {
	tr := makeTrigger(0)
	sentinel := errors.New("fn error")
	r := NewTriggerRunner(tr, func(_ context.Context) error { return sentinel })
	tr.Fire()
	if err := r.RunOnce(context.Background()); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestTriggerRunnerNilArgDefaults(t *testing.T) {
	r := NewTriggerRunner(nil, nil)
	if r.trigger == nil {
		t.Fatal("expected non-nil trigger")
	}
	if r.fn == nil {
		t.Fatal("expected non-nil fn")
	}
}

func TestTriggerRunnerRunLoopStopsOnContextCancel(t *testing.T) {
	tr := makeTrigger(0)
	calls := 0
	r := NewTriggerRunner(tr, func(_ context.Context) error {
		calls++
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	// fire once so the loop executes at least one iteration
	tr.Fire()
	err := r.RunLoop(ctx)
	if err == nil {
		t.Fatal("expected non-nil error from RunLoop")
	}
	if calls == 0 {
		t.Fatal("expected fn to be called at least once")
	}
}
