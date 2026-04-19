package watch

import (
	"context"
	"sync"
)

// Drain collects items from a channel until it is closed or the context is
// cancelled, then returns everything it received.
type Drain[T any] struct {
	mu    sync.Mutex
	items []T
}

// NewDrain creates a new Drain instance.
func NewDrain[T any]() *Drain[T] {
	return &Drain[T]{}
}

// Run reads from ch until it is closed or ctx is done.
func (d *Drain[T]) Run(ctx context.Context, ch <-chan T) {
	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-ch:
			if !ok {
				return
			}
			d.mu.Lock()
			d.items = append(d.items, v)
			d.mu.Unlock()
		}
	}
}

// Snapshot returns a copy of all collected items.
func (d *Drain[T]) Snapshot() []T {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([]T, len(d.items))
	copy(out, d.items)
	return out
}

// Len returns the number of items collected so far.
func (d *Drain[T]) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.items)
}

// Clear discards all collected items.
func (d *Drain[T]) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.items = nil
}
