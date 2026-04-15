package watch_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func TestRetrierIntegration(t *testing.T) {
	policy := watch.RetryPolicy{
		MaxAttempts: 4,
		Backoff: watch.NewBackoff(watch.BackoffPolicy{
			Initial: 2 * time.Millisecond,
			Max:     10 * time.Millisecond,
			Factor:  2.0,
		}),
	}
	logger := watch.NewLogger(nil)
	retrier := watch.NewRetrier(policy, logger)

	attempts := 0
	sentinel := errors.New("transient")

	start := time.Now()
	err := retrier.Run(context.Background(), func(_ context.Context) error {
		attempts++
		if attempts < 3 {
			return sentinel
		}
		return nil
	})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
	// At least 2ms (initial) + 4ms (doubled) = 6ms of backoff
	if elapsed < 6*time.Millisecond {
		t.Errorf("expected backoff delay, elapsed=%v", elapsed)
	}
}
