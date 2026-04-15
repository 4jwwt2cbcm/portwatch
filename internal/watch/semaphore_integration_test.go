package watch_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

// TestSemaphoreIntegration verifies that a semaphore correctly gates a
// burst of goroutines so that no more than N run concurrently.
func TestSemaphoreIntegration(t *testing.T) {
	const (
		limit   = 2
		workers = 8
		workDur = 10 * time.Millisecond
	)

	s := watch.NewSemaphore(limit)
	ctx := context.Background()

	var active int64
	errCh := make(chan error, workers)

	for i := 0; i < workers; i++ {
		go func() {
			if err := s.Acquire(ctx); err != nil {
				errCh <- err
				return
			}
			defer s.Release()

			cur := atomic.AddInt64(&active, 1)
			if cur > int64(limit) {
				t.Errorf("active goroutines %d exceeded semaphore limit %d", cur, limit)
			}
			time.Sleep(workDur)
			atomic.AddInt64(&active, -1)
			errCh <- nil
		}()
	}

	for i := 0; i < workers; i++ {
		if err := <-errCh; err != nil {
			t.Fatalf("worker returned unexpected error: %v", err)
		}
	}
}
