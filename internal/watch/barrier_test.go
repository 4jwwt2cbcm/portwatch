package watch

import (
	"context"
	"sync"
	"testing"
	"time"
)

func makeBarrier(n int) *Barrier {
	return NewBarrier(n)
}

func TestBarrierDefaultsToOneOnZero(t *testing.T) {
	b := NewBarrier(0)
	ctx := context.Background()
	done := make(chan error, 1)
	go func() { done <- b.Wait(ctx) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for barrier")
	}
}

func TestBarrierReleasesWhenTargetMet(t *testing.T) {
	b := makeBarrier(3)
	ctx := context.Background()
	var wg sync.WaitGroup
	errs := make([]error, 3)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			errs[idx] = b.Wait(ctx)
		}(i)
	}
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("barrier did not release in time")
	}
	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d: unexpected error %v", i, err)
		}
	}
}

func TestBarrierCancelledContextReturnsError(t *testing.T) {
	b := makeBarrier(2)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- b.Wait(ctx) }()
	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected context error, got nil")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out")
	}
}

func TestBarrierArrivedCount(t *testing.T) {
	b := makeBarrier(3)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			_ = b.Wait(ctx)
		}()
	}
	wg.Wait()
	time.Sleep(20 * time.Millisecond)
	if got := b.Arrived(); got < 1 {
		t.Errorf("expected at least 1 arrived, got %d", got)
	}
}

func TestBarrierResetUnblocksWaiters(t *testing.T) {
	b := makeBarrier(5)
	ctx := context.Background()
	done := make(chan error, 1)
	go func() { done <- b.Wait(ctx) }()
	time.Sleep(20 * time.Millisecond)
	b.Reset()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("reset did not unblock waiter")
	}
}
