package watch

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestBulkheadRunnerRunsFunction(t *testing.T) {
	b := makeBulkhead(2, 0)
	ran := false
	r := NewBulkheadRunner(b, func(_ context.Context) error {
		ran = true
		return nil
	})
	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ran {
		t.Error("function was not called")
	}
}

func TestBulkheadRunnerPropagatesError(t *testing.T) {
	sentinel := errors.New("fn error")
	b := makeBulkhead(2, 0)
	r := NewBulkheadRunner(b, func(_ context.Context) error { return sentinel })
	if err := r.Run(context.Background()); !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestBulkheadRunnerBlocksWhenExhausted(t *testing.T) {
	b := makeBulkhead(1, 0)
	unblock := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = b.Do(context.Background(), func() error { <-unblock; return nil })
	}()
	time.Sleep(10 * time.Millisecond)
	r := NewBulkheadRunner(b, func(_ context.Context) error { return nil })
	err := r.Run(context.Background())
	if err != ErrBulkheadFull {
		t.Errorf("expected ErrBulkheadFull, got %v", err)
	}
	close(unblock)
	wg.Wait()
}

func TestBulkheadRunnerNilBulkheadDefaults(t *testing.T) {
	r := NewBulkheadRunner(nil, func(_ context.Context) error { return nil })
	if r.bulkhead == nil {
		t.Error("expected non-nil bulkhead")
	}
}

func TestBulkheadRunnerNilFnDefaults(t *testing.T) {
	b := makeBulkhead(2, 0)
	r := NewBulkheadRunner(b, nil)
	if err := r.Run(context.Background()); err != nil {
		t.Errorf("nil fn should be noop, got %v", err)
	}
}
