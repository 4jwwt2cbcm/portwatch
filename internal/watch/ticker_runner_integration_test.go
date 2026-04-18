package watch_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func TestTickerRunnerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	var count int64
	r := watch.NewTickerRunner(30*time.Millisecond, func(_ context.Context) error {
		atomic.AddInt64(&count, 1)
		return nil
	}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_ = r.Run(ctx)

	got := atomic.LoadInt64(&count)
	if got < 4 {
		t.Fatalf("expected at least 4 ticks in 200ms with 30ms interval, got %d", got)
	}
}
