package watch_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"portwatch/internal/watch"
)

func TestLimiterIntegration(t *testing.T) {
	policy := watch.LimiterPolicy{MaxCalls: 3, Window: 200 * time.Millisecond}
	l := watch.NewLimiter(policy)
	runner := watch.NewLimiterRunner(l, nil)

	var count atomic.Int32
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Fire 3 calls quickly — all should succeed.
	for i := 0; i < 3; i++ {
		if err := runner.Run(ctx); err != nil {
			t.Fatalf("call %d failed: %v", i+1, err)
		}
		count.Add(1)
	}
	if count.Load() != 3 {
		t.Errorf("expected 3 calls, got %d", count.Load())
	}

	// 4th call should block until window expires.
	start := time.Now()
	if err := runner.Run(ctx); err != nil {
		t.Fatalf("4th call failed: %v", err)
	}
	elapsed := time.Since(start)
	if elapsed < 150*time.Millisecond {
		t.Errorf("expected wait for window, got %v", elapsed)
	}
}
