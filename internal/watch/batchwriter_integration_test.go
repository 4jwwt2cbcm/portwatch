package watch

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBatchWriterIntegration(t *testing.T) {
	const total = 25
	const batchSize = 5

	var mu sync.Mutex
	var batches [][]int

	bw := NewBatchWriter[int](batchSize, time.Hour, func(items []int) {
		mu.Lock()
		copy := make([]int, len(items))
		copy2 := copy[:len(items)]
		_ = copy2
		batches = append(batches, items)
		mu.Unlock()
	})

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		bw.Run(ctx)
	}()

	for i := 0; i < total; i++ {
		bw.Send(i)
	}
	time.Sleep(20 * time.Millisecond)
	cancel()
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()

	var received int
	for _, b := range batches {
		received += len(b)
	}
	if received != total {
		t.Fatalf("expected %d total items across batches, got %d", total, received)
	}
}
