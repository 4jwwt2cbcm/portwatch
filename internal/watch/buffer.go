package watch

import "sync"

// Buffer accumulates items up to a capacity, then drains them as a batch.
type Buffer[T any] struct {
	mu       sync.Mutex
	items    []T
	capacity int
}

// NewBuffer creates a Buffer with the given capacity (minimum 1).
func NewBuffer[T any](capacity int) *Buffer[T] {
	if capacity < 1 {
		capacity = 1
	}
	return &Buffer[T]{
		items:    make([]T, 0, capacity),
		capacity: capacity,
	}
}

// Add appends an item. Returns true if the buffer is now full.
func (b *Buffer[T]) Add(item T) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.items = append(b.items, item)
	return len(b.items) >= b.capacity
}

// Flush returns all buffered items and resets the buffer.
func (b *Buffer[T]) Flush() []T {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]T, len(b.items))
	copy(out, b.items)
	b.items = b.items[:0]
	return out
}

// Len returns the current number of buffered items.
func (b *Buffer[T]) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.items)
}

// Cap returns the configured capacity.
func (b *Buffer[T]) Cap() int {
	return b.capacity
}
