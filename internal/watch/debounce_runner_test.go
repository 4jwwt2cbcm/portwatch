package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDebounceRunnerFiresOnce(t *testing.T) {
	d := makeDebounce(40 * time.Millisecond)
	var calls int
	r := NewDebounceRunner(d, func(_ context.Context) error {
		calls++
		return nil
	})
	ctx := context.Background()
	for i := 0; i < 4; i++ {
		r.Schedule(ctx)
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(80 * time.Millisecond)
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDebounceRunnerPropagatesError(t *testing.T) {
	d := makeDebounce(30 * time.Millisecond)
	sentinel := errors.New("boom")
	r := NewDebounceRunner(d, func(_ context.Context) error { return sentinel })
	ctx := context.Background()
	err := r.Wait(ctx, 100*time.Millisecond)
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestDebounceRunnerCancelledContext(t *testing.T) {
	d := makeDebounce(200 * time.Millisecond)
	r := NewDebounceRunner(d, func(_ context.Context) error { return nil })
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	err := r.Wait(ctx, 200*time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}

func TestDebounceRunnerNilArgDefaults(t *testing.T) {
	r := NewDebounceRunner(nil, nil)
	if r.debounce == nil {
		t.Fatal("expected non-nil debounce")
	}
	ctx := context.Background()
	err := r.Wait(ctx, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
