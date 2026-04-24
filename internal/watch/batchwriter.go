package watch

import (
	"context"
	"time"
)

// BatchHandler is called with a slice of flushed items.
type BatchHandler[T any] func(items []T)

// BatchWriter collects items into a Buffer and flushes them either when
// the buffer is full or a flush interval elapses.
type BatchWriter[T any] struct {
	buffer   *Buffer[T]
	interval time.Duration
	handler  BatchHandler[T]
}

// NewBatchWriter creates a BatchWriter with the given capacity, flush interval,
// and handler. Capacity and interval are clamped to sane minimums.
func NewBatchWriter[T any](capacity int, interval time.Duration, handler BatchHandler[T]) *BatchWriter[T] {
	if interval < time.Millisecond {
		interval = time.Millisecond
	}
	return &BatchWriter[T]{
		buffer:   NewBuffer[T](capacity),
		interval: interval,
		handler:  handler,
	}
}

// Send adds an item. If the buffer becomes full the batch is flushed immediately.
func (bw *BatchWriter[T]) Send(item T) {
	if bw.buffer.Add(item) {
		bw.flush()
	}
}

// SendAll adds multiple items, flushing immediately whenever the buffer fills up.
func (bw *BatchWriter[T]) SendAll(items []T) {
	for _, item := range items {
		bw.Send(item)
	}
}

// Run starts the periodic flush loop, blocking until ctx is cancelled.
func (bw *BatchWriter[T]) Run(ctx context.Context) {
	ticker := time.NewTicker(bw.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			bw.flush()
			return
		case <-ticker.C:
			bw.flush()
		}
	}
}

func (bw *BatchWriter[T]) flush() {
	items := bw.buffer.Flush()
	if len(items) > 0 {
		bw.handler(items)
	}
}
