package watch

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestPriorityRunnerStopsOnContextCancel(t *testing.T) {
	q := NewPriorityQueue[PriorityTask](0)
	r := NewPriorityRunner(q, 10*time.Millisecond, nil)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- r.Run(ctx) }()
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("runner did not stop")
	}
}

func TestPriorityRunnerExecutesTasksInOrder(t *testing.T) {
	q := NewPriorityQueue[PriorityTask](0)
	var mu sync.Mutex
	var order []string

	record := func(name string) func(context.Context) error {
		return func(_ context.Context) error {
			mu.Lock()
			order = append(order, name)
			mu.Unlock()
			return nil
		}
	}

	q.Push(PriorityTask{Name: "low", Priority: 1, Fn: record("low")}, 1)
	q.Push(PriorityTask{Name: "high", Priority: 10, Fn: record("high")}, 10)
	q.Push(PriorityTask{Name: "mid", Priority: 5, Fn: record("mid")}, 5)

	ctx, cancel := context.WithCancel(context.Background())
	r := NewPriorityRunner(q, 10*time.Millisecond, nil)
	go r.Run(ctx)

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		n := len(order)
		mu.Unlock()
		if n >= 3 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	cancel()

	mu.Lock()
	defer mu.Unlock()
	if len(order) < 3 {
		t.Fatalf("expected 3 tasks executed, got %d", len(order))
	}
	if order[0] != "high" || order[1] != "mid" || order[2] != "low" {
		t.Fatalf("unexpected order: %v", order)
	}
}

func TestPriorityRunnerNilLoggerDefaults(t *testing.T) {
	q := NewPriorityQueue[PriorityTask](0)
	r := NewPriorityRunner(q, 0, nil)
	if r.logger == nil {
		t.Fatal("expected default logger")
	}
	if r.pollInterval <= 0 {
		t.Fatal("expected positive poll interval")
	}
}
