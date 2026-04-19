package watch_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func TestBarrierIntegration(t *testing.T) {
	const workers = 5
	b := watch.NewBarrier(workers)
	ctx := context.Background()

	var readyCount atomic.Int32
	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			readyCount.Add(1)
			if err := b.Wait(ctx); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	close(start)

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	select {
	case <-done:
		if n := readyCount.Load(); n != workers {
			t.Errorf("expected %d ready, got %d", workers, n)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("integration test timed out")
	}
}
