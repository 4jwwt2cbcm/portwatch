package watch

import "context"

// Semaphore limits the number of concurrent operations.
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore creates a Semaphore with the given concurrency limit.
// n must be greater than zero.
func NewSemaphore(n int) *Semaphore {
	if n <= 0 {
		n = 1
	}
	return &Semaphore{ch: make(chan struct{}, n)}
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Returns ctx.Err() if the context is cancelled before a slot is acquired.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release frees a previously acquired slot.
// It panics if Release is called more times than Acquire.
func (s *Semaphore) Release() {
	select {
	case <-s.ch:
	default:
		panic("watch: semaphore Release called without matching Acquire")
	}
}

// TryAcquire attempts to acquire a slot without blocking.
// Returns true if successful, false if no slot is available.
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// Available returns the number of free slots.
func (s *Semaphore) Available() int {
	return cap(s.ch) - len(s.ch)
}
