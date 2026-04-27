package watch

import (
	"context"
	"sync"
	"testing"
	"time"
)

func makeBulkhead(max, queue int) *Bulkhead {
	return NewBulkhead(BulkheadPolicy{MaxConcurrent: max, QueueDepth: queue})
}

func TestDefaultBulkheadPolicyValues(t *testing.T) {
	p := DefaultBulkheadPolicy()
	if p.MaxConcurrent != 8 {
		t.Errorf("MaxConcurrent: got %d, want 8", p.MaxConcurrent)
	}
	if p.QueueDepth != 16 {
		t.Errorf("QueueDepth: got %d, want 16", p.QueueDepth)
	}
}

func TestBulkheadAllowsUpToMax(t *testing.T) {
	b := makeBulkhead(2, 0)
	ctx := context.Background()
	var wg sync.WaitGroup
	ready := make(chan struct{})
	unblock := make(chan struct{})
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = b.Do(ctx, func() error {
				ready <- struct{}{}
				<-unblock
				return nil
			})
		}()
	}
	<-ready
	<-ready
	if b.Active() != 2 {
		t.Errorf("Active: got %d, want 2", b.Active())
	}
	close(unblock)
	wg.Wait()
}

func TestBulkheadRejectsWhenFullNoQueue(t *testing.T) {
	b := makeBulkhead(1, 0)
	ctx := context.Background()
	unblock := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = b.Do(ctx, func() error { <-unblock; return nil })
	}()
	time.Sleep(10 * time.Millisecond)
	err := b.Do(ctx, func() error { return nil })
	if err != ErrBulkheadFull {
		t.Errorf("expected ErrBulkheadFull, got %v", err)
	}
	if b.Shed() != 1 {
		t.Errorf("Shed: got %d, want 1", b.Shed())
	}
	close(unblock)
	wg.Wait()
}

func TestBulkheadQueuesAndRuns(t *testing.T) {
	b := makeBulkhead(1, 4)
	ctx := context.Background()
	unblock := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = b.Do(ctx, func() error { <-unblock; return nil })
	}()
	time.Sleep(10 * time.Millisecond)
	var ran atomic.Int64
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = b.Do(ctx, func() error { ran.Add(1); return nil })
	}()
	time.Sleep(5 * time.Millisecond)
	if b.Queued() != 1 {
		t.Errorf("Queued: got %d, want 1", b.Queued())
	}
	close(unblock)
	wg.Wait()
	if ran.Load() != 1 {
		t.Errorf("queued fn did not run")
	}
}

func TestBulkheadContextCancelRemovesFromQueue(t *testing.T) {
	b := makeBulkhead(1, 4)
	unblock := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = b.Do(context.Background(), func() error { <-unblock; return nil })
	}()
	time.Sleep(10 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := b.Do(ctx, func() error { return nil })
		if err != context.Canceled {
			t.Errorf("expected Canceled, got %v", err)
		}
	}()
	time.Sleep(5 * time.Millisecond)
	cancel()
	close(unblock)
	wg.Wait()
}

func TestBulkheadResetClearsShed(t *testing.T) {
	b := makeBulkhead(1, 0)
	ctx := context.Background()
	unblock := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = b.Do(ctx, func() error { <-unblock; return nil })
	}()
	time.Sleep(10 * time.Millisecond)
	_ = b.Do(ctx, func() error { return nil })
	if b.Shed() == 0 {
		t.Error("expected shed > 0")
	}
	b.Reset()
	if b.Shed() != 0 {
		t.Errorf("Reset: Shed should be 0, got %d", b.Shed())
	}
	close(unblock)
	wg.Wait()
}
