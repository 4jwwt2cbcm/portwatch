package watch

import (
	"sync"
	"time"
)

// Epoch tracks a monotonically incrementing generation counter that advances
// on each reset. It is useful for invalidating cached state when a new scan
// cycle begins.

// EpochSnapshot holds a point-in-time view of the epoch state.
type EpochSnapshot struct {
	Generation uint64
	StartedAt  time.Time
	ResetAt    time.Time
}

// Epoch is a thread-safe generation counter.
type Epoch struct {
	mu         sync.RWMutex
	generation uint64
	startedAt  time.Time
	resetAt    time.Time
	now        func() time.Time
}

// NewEpoch creates a new Epoch with generation 0.
func NewEpoch() *Epoch {
	now := time.Now()
	return &Epoch{
		generation: 0,
		startedAt:  now,
		resetAt:    now,
		now:        time.Now,
	}
}

// Generation returns the current generation number.
func (e *Epoch) Generation() uint64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.generation
}

// Advance increments the generation counter and records the reset time.
func (e *Epoch) Advance() uint64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.generation++
	e.resetAt = e.now()
	return e.generation
}

// Reset sets the generation back to zero and records the reset time.
func (e *Epoch) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.generation = 0
	e.resetAt = e.now()
}

// Snapshot returns a consistent point-in-time copy of the epoch state.
func (e *Epoch) Snapshot() EpochSnapshot {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return EpochSnapshot{
		Generation: e.generation,
		StartedAt:  e.startedAt,
		ResetAt:    e.resetAt,
	}
}

// Since returns true if the given generation is older than the current one.
func (e *Epoch) Since(gen uint64) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return gen < e.generation
}
