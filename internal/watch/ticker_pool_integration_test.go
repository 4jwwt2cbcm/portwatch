package watch_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

// TestTickerPoolIntegration verifies that multiple named tickers in a pool
// all fire within a reasonable time window and that StopAll cleanly shuts
// down every ticker without a race or panic.
func TestTickerPoolIntegration(t *testing.T) {
	pool := watch.NewTickerPool(watch.TickerPoolPolicy{
		MaxTickers:      4,
		DefaultInterval: 20 * time.Millisecond,
	})

	names := []string{"alpha", "beta", "gamma"}
	for _, n := range names {
		if ok := pool.Add(n, 20*time.Millisecond); !ok {
			t.Fatalf("failed to add ticker %q", n)
		}
	}

	if pool.Len() != len(names) {
		t.Fatalf("expected %d tickers, got %d", len(names), pool.Len())
	}

	// Drain at least one tick from each ticker.
	timeout := time.After(500 * time.Millisecond)
	for _, n := range names {
		ch := pool.C(n)
		if ch == nil {
			t.Fatalf("nil channel for ticker %q", n)
		}
		select {
		case <-ch:
		case <-timeout:
			t.Fatalf("timed out waiting for ticker %q to fire", n)
		}
	}

	pool.StopAll()

	if pool.Len() != 0 {
		t.Fatalf("expected empty pool after StopAll, got %d", pool.Len())
	}
}
