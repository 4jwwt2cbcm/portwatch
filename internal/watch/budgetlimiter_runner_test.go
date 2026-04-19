package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestBudgetRunnerAllowsWhenBudgetAvailable(t *testing.T) {
	b := makeBudgetLimiter(5, time.Minute)
	called := false
	r := NewBudgetRunner(b, func(ctx context.Context) error {
		called = true
		return nil
	})
	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected fn to be called")
	}
}

func TestBudgetRunnerBlocksWhenExhausted(t *testing.T) {
	b := makeBudgetLimiter(1, time.Minute)
	r := NewBudgetRunner(b, func(ctx context.Context) error { return nil })
	r.Run(context.Background())
	err := r.Run(context.Background())
	if !errors.Is(err, ErrBudgetExhausted) {
		t.Errorf("expected ErrBudgetExhausted, got %v", err)
	}
}

func TestBudgetRunnerPropagatesError(t *testing.T) {
	b := makeBudgetLimiter(5, time.Minute)
	sentinel := errors.New("fn error")
	r := NewBudgetRunner(b, func(ctx context.Context) error { return sentinel })
	if err := r.Run(context.Background()); !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestBudgetRunnerNilLimiterDefaults(t *testing.T) {
	r := NewBudgetRunner(nil, nil)
	if err := r.Run(context.Background()); err != nil {
		t.Errorf("expected nil error with defaults, got %v", err)
	}
}
