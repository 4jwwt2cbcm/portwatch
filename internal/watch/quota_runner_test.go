package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestQuotaRunnerAllowsWhenBudgetAvailable(t *testing.T) {
	q := makeQuota(5, time.Minute)
	r := NewQuotaRunner(q, func(_ context.Context) error { return nil })
	if err := r.Run(context.Background()); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestQuotaRunnerBlocksWhenExhausted(t *testing.T) {
	q := makeQuota(1, time.Minute)
	r := NewQuotaRunner(q, func(_ context.Context) error { return nil })
	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("first run should succeed: %v", err)
	}
	if err := r.Run(context.Background()); !errors.Is(err, ErrQuotaExceeded) {
		t.Errorf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestQuotaRunnerPropagatesError(t *testing.T) {
	sentinel := errors.New("fn error")
	q := makeQuota(5, time.Minute)
	r := NewQuotaRunner(q, func(_ context.Context) error { return sentinel })
	if err := r.Run(context.Background()); !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestQuotaRunnerNilQuotaDefaults(t *testing.T) {
	r := NewQuotaRunner(nil, func(_ context.Context) error { return nil })
	if r.quota == nil {
		t.Error("expected non-nil default quota")
	}
}

func TestQuotaRunnerNilFnDefaults(t *testing.T) {
	q := makeQuota(5, time.Minute)
	r := NewQuotaRunner(q, nil)
	if err := r.Run(context.Background()); err != nil {
		t.Errorf("expected nil from default fn, got %v", err)
	}
}
