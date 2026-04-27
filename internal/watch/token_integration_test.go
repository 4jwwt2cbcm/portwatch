package watch_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

// TestTokenPoolIntegration verifies that the refill loop restores
// tokens so a runner that was previously blocked can proceed.
func TestTokenPoolIntegration(t *testing.T) {
	pool := watch.NewTokenPool(watch.TokenPolicy{
		Capacity:     2,
		RefillEvery:  30 * time.Millisecond,
		RefillAmount: 2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go pool.Run(ctx) //nolint:errcheck

	// Drain the pool.
	pool.Take()
	pool.Take()

	runner := watch.NewTokenRunner(pool, func(_ context.Context) error { return nil })

	// Should be rejected immediately.
	if err := runner.Run(ctx); !errors.Is(err, watch.ErrNoToken) {
		t.Fatalf("expected ErrNoToken before refill, got %v", err)
	}

	// Wait for at least one refill cycle.
	time.Sleep(60 * time.Millisecond)

	// Should succeed after refill.
	if err := runner.Run(ctx); err != nil {
		t.Errorf("expected success after refill, got %v", err)
	}
}
