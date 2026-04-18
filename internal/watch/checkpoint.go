package watch

import (
	"sync"
	"time"
)

// Checkpoint records the last successful run time and allows callers to
// determine whether enough time has elapsed since the previous checkpoint.
type Checkpoint struct {
	mu       sync.Mutex
	last     time.Time
	interval time.Duration
}

// NewCheckpoint creates a Checkpoint with the given minimum interval between
// successful runs. If interval is zero it defaults to one minute.
func NewCheckpoint(interval time.Duration) *Checkpoint {
	if interval <= 0 {
		interval = time.Minute
	}
	return &Checkpoint{interval: interval}
}

// Mark records now as the last successful checkpoint time.
func (c *Checkpoint) Mark() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.last = time.Now()
}

// Due reports whether the checkpoint interval has elapsed since the last Mark.
// It returns true if Mark has never been called.
func (c *Checkpoint) Due() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.last.IsZero() {
		return true
	}
	return time.Since(c.last) >= c.interval
}

// Last returns the time of the most recent Mark, or the zero time if Mark has
// never been called.
func (c *Checkpoint) Last() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.last
}

// Reset clears the last checkpoint, causing Due to return true again.
func (c *Checkpoint) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.last = time.Time{}
}
