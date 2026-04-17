package watch

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBatchWriterFlushesOnFull(t *testing.T) {
	var mu sync.Mutex
	var got [][]int
	bw := NewBatchWriter[int](3, time.Hour, func(items []int) {
		mu.Lock()
		got = append(got, items)
		mu.Unlock()
	})
	bw.Send(1)
	bw.Send(2)
	bw.Send(3) // triggers flush
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 || len(got[0]) != 3 {
		t.Fatalf("expected one batch of 3, got %v", got)
	}
}

func TestBatchWriterFlushesOnInterval(t *testing.T) {
	var mu sync.Mutex
	var got [][]string
	bw := NewBatchWriter[string](100, 20*time.Millisecond, func(items []string) {
		mu.Lock()
		got = append(got, items)
		mu.Unlock()
	})
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	bw.Send("x")
	bw.Run(ctx)
	mu.Lock()
	defer mu.Unlock()
	if len(got) == 0 {
		t.Fatal("expected at least one flush from ticker")
	}
}

func TestBatchWriterFlushesOnContextCancel(t *testing.T) {
	var mu sync.Mutex
	var got [][]int
	bw := NewBatchWriter[int](100, time.Hour, func(items []int) {
		mu.Lock()
		got = append(got, items)
		mu.Unlock()
	})
	bw.Send(42)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	bw.Run(ctx)
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 || got[0][0] != 42 {
		t.Fatalf("expected flush on cancel, got %v", got)
	}
}

func TestBatchWriterEmptyFlushSkipsHandler(t *testing.T) {
	called := false
	bw := NewBatchWriter[int](4, 10*time.Millisecond, func(_ []int) {
		called = true
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	bw.Run(ctx)
	if called {
		t.Fatal("handler should not be called for empty buffer")
	}
}
