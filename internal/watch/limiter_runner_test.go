package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestLimiterRunnerRunsFunction(t *testing.T) {
	called := false
	r := NewLimiterRunner(makeLimiter(5, time.Second), func(ctx context.Context) error {
		called = true
		return nil
	})
	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected function to be called")
	}
}

func TestLimiterRunnerPropagatesError(t *testing.T) {
	sentinel := errors.New("fn error")
	r := NewLimiterRunner(makeLimiter(5, time.Second), func(ctx context.Context) error {
		return sentinel
	})
	if err := r.Run(context.Background()); !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestLimiterRunnerBlocksWhenExhausted(t *testing.T) {
	l := makeLimiter(1, time.Hour)
	l.Allow() // exhaust
	r := NewLimiterRunner(l, func(ctx context.Context) error { return nil })
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := r.Run(ctx); err == nil {
		t.Error("expected error when context cancelled while blocked")
	}
}

func TestLimiterRunnerNilLimiterDefaults(t *testing.T) {
	r := NewLimiterRunner(nil, nil)
	if err := r.Run(context.Background()); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
