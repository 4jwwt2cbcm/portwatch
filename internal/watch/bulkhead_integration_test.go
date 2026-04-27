package watch_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"portwatch/internal/watch"
)

func TestBulkheadIntegration(t *testing.T) {
	const (
		maxConcurrent = 3
		queueDepth    = 5
		totalWorkers  = 10
	)

	b := watch.NewBulkhead(watch.BulkheadPolicy{
		MaxConcurrent: maxConcurrent,
		QueueDepth:    queueDepth,
	})

	var (
		ran     atomic.Int64
		shed    atomic.Int64
		wg      sync.WaitGroup
		unblock = make(chan struct{})
	)

	// Fill all concurrent slots so queuing / shedding is exercised.
	for i := 0; i < maxConcurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = b.Do(context.Background(), func() error {
				<-unblock
				return nil
			})
		}()
	}
	time.Sleep(20 * time.Millisecond)

	ctx := context.Background()
	for i := 0; i < totalWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := b.Do(ctx, func() error {
				ran.Add(1)
				return nil
			})
			if err == watch.ErrBulkheadFull {
				shed.Add(1)
			}
		}()
	}
	time.Sleep(20 * time.Millisecond)
	close(unblock)
	wg.Wait()

	total := ran.Load() + shed.Load()
	if total != totalWorkers {
		t.Errorf("ran+shed = %d, want %d", total, totalWorkers)
	}
	if shed.Load() < 1 {
		t.Errorf("expected at least 1 shed call, got %d", shed.Load())
	}
	t.Logf("ran=%d shed=%d", ran.Load(), shed.Load())
}
