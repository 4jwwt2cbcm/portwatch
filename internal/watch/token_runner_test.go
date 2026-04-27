package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestTokenRunnerSucceedsWhenTokenAvailable(t *testing.T) {
	pool := makeTokenPool(1, 1, time.Second)
	r := NewTokenRunner(pool, func(_ context.Context) error { return nil })
	if err := r.Run(context.Background()); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestTokenRunnerReturnsErrNoTokenWhenEmpty(t *testing.T) {
	pool := makeTokenPool(1, 1, time.Hour)
	pool.Take() // drain
	r := NewTokenRunner(pool, func(_ context.Context) error { return nil })
	if err := r.Run(context.Background()); !errors.Is(err, ErrNoToken) {
		t.Errorf("expected ErrNoToken, got %v", err)
	}
}

func TestTokenRunnerPropagatesError(t *testing.T) {
	pool := makeTokenPool(5, 1, time.Second)
	sentinel := errors.New("boom")
	r := NewTokenRunner(pool, func(_ context.Context) error { return sentinel })
	if err := r.Run(context.Background()); !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestTokenRunnerNilPoolDefaults(t *testing.T) {
	r := NewTokenRunner(nil, nil)
	if r.pool == nil {
		t.Error("expected default pool to be set")
	}
	if err := r.Run(context.Background()); err != nil {
		t.Errorf("expected nil error from default noop fn, got %v", err)
	}
}

func TestTokenRunnerConsumesToken(t *testing.T) {
	pool := makeTokenPool(3, 1, time.Hour)
	r := NewTokenRunner(pool, func(_ context.Context) error { return nil })
	r.Run(context.Background()) //nolint:errcheck
	if got := pool.Available(); got != 2 {
		t.Errorf("expected 2 tokens remaining, got %d", got)
	}
}
