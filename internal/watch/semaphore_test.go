package watch

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewSemaphoreDefaultsToOne(t *testing.T) {
	s := NewSemaphore(0)
	if s.Available() != 1 {
		t.Fatalf("expected 1 available slot, got %d", s.Available())
	}
}

func TestSemaphoreAcquireAndRelease(t *testing.T) {
	s := NewSemaphore(2)
	ctx := context.Background()

	if err := s.Acquire(ctx); err != nil {
		t.Fatalf("unexpected error on first Acquire: %v", err)
	}
	if s.Available() != 1 {
		t.Fatalf("expected 1 available after acquire, got %d", s.Available())
	}
	s.Release()
	if s.Available() != 2 {
		t.Fatalf("expected 2 available after release, got %d", s.Available())
	}
}

func TestSemaphoreTryAcquire(t *testing.T) {
	s := NewSemaphore(1)

	if !s.TryAcquire() {
		t.Fatal("expected TryAcquire to succeed on empty semaphore")
	}
	if s.TryAcquire() {
		t.Fatal("expected TryAcquire to fail when semaphore is full")
	}
	s.Release()
	if !s.TryAcquire() {
		t.Fatal("expected TryAcquire to succeed after release")
	}
	s.Release()
}

func TestSemaphoreAcquireRespectsContextCancel(t *testing.T) {
	s := NewSemaphore(1)
	_ = s.TryAcquire() // fill the slot

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := s.Acquire(ctx)
	if err == nil {
		t.Fatal("expected error when context is cancelled, got nil")
	}
}

func TestSemaphoreLimitsConcurrency(t *testing.T) {
	const limit = 3
	s := NewSemaphore(limit)
	ctx := context.Background()

	var mu sync.Mutex
	var peak, current int
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Acquire(ctx)
			defer s.Release()
			mu.Lock()
			current++
			if current > peak {
				peak = current
			}
			mu.Unlock()
			time.Sleep(5 * time.Millisecond)
			mu.Lock()
			current--
			mu.Unlock()
		}()
	}
	wg.Wait()

	if peak > limit {
		t.Fatalf("peak concurrency %d exceeded limit %d", peak, limit)
	}
}
