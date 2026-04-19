package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDeadlineRunnerRunsWhenNotExpired(t *testing.T) {
	called := false
	r := NewDeadlineRunner(nil, func(ctx context.Context) error {
		called = true
		return nil
	})
	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected fn to be called")
	}
}

func TestDeadlineRunnerSkipsWhenExpired(t *testing.T) {
	d := makeDeadline(time.Now().Add(-time.Second))
	r := NewDeadlineRunner(d, func(ctx context.Context) error {
		t.Fatal("fn should not be called")
		return nil
	})
	err := r.Run(context.Background())
	if !errors.Is(err, ErrDeadlineExceeded) {
		t.Fatalf("expected ErrDeadlineExceeded, got %v", err)
	}
}

func TestDeadlineRunnerPropagatesError(t *testing.T) {
	sentinel := errors.New("boom")
	r := NewDeadlineRunner(nil, func(ctx context.Context) error {
		return sentinel
	})
	if err := r.Run(context.Background()); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel, got %v", err)
	}
}

func TestDeadlineRunnerNilArgDefaults(t *testing.T) {
	r := NewDeadlineRunner(nil, nil)
	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeadlineRunnerRunUntilStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	r := NewDeadlineRunner(nil, func(ctx context.Context) error { return nil })
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()
	err := r.RunUntil(ctx, 10*time.Millisecond)
	if err != context.Canceled {
		t.Fatalf("expected Canceled, got %v", err)
	}
}

func TestDeadlineRunnerRunUntilStopsOnExpiry(t *testing.T) {
	d := NewDeadline(DeadlinePolicy{At: time.Now().Add(40 * time.Millisecond)})
	r := NewDeadlineRunner(d, func(ctx context.Context) error { return nil })
	err := r.RunUntil(context.Background(), 10*time.Millisecond)
	if !errors.Is(err, ErrDeadlineExceeded) {
		t.Fatalf("expected ErrDeadlineExceeded, got %v", err)
	}
}
