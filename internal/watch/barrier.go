package watch

import (
	"context"
	"sync"
)

// Barrier blocks goroutines until a required number have arrived,
// then releases all of them simultaneously.
type Barrier struct {
	mu      sync.Mutex
	cond    *sync.Cond
	target  int
	arrived int
	gen     int
}

// NewBarrier creates a Barrier that releases when n goroutines have arrived.
// If n < 1, it defaults to 1.
func NewBarrier(n int) *Barrier {
	if n < 1 {
		n = 1
	}
	b := &Barrier{target: n}
	b.cond = sync.NewCond(&b.mu)
	return b
}

// Wait blocks until the required number of goroutines have called Wait,
// or until ctx is cancelled. Returns ctx.Err() if context is done.
func (b *Barrier) Wait(ctx context.Context) error {
	b.mu.Lock()
	gen := b.gen
	b.arrived++
	if b.arrived >= b.target {
		b.arrived = 0
		b.gen++
		b.cond.Broadcast()
		b.mu.Unlock()
		return nil
	}

	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			b.mu.Lock()
			b.cond.Broadcast()
			b.mu.Unlock()
		case <-done:
		}
	}()

	for b.gen == gen && ctx.Err() == nil {
		b.cond.Wait()
	}
	close(done)

	if b.gen == gen {
		b.arrived--
		b.mu.Unlock()
		return ctx.Err()
	}
	b.mu.Unlock()
	return nil
}

// Reset resets the barrier to its initial state.
func (b *Barrier) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.arrived = 0
	b.gen++
	b.cond.Broadcast()
}

// Arrived returns the number of goroutines currently waiting.
func (b *Barrier) Arrived() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.arrived
}
