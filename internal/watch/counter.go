package watch

import "sync"

// Counter is a thread-safe monotonic counter with optional reset support.
type Counter struct {
	mu    sync.Mutex
	value int64
}

// NewCounter returns a new Counter starting at zero.
func NewCounter() *Counter {
	return &Counter{}
}

// Inc increments the counter by 1 and returns the new value.
func (c *Counter) Inc() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
	return c.value
}

// Add increments the counter by n and returns the new value.
func (c *Counter) Add(n int64) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += n
	return c.value
}

// Value returns the current counter value.
func (c *Counter) Value() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// Reset sets the counter back to zero and returns the previous value.
func (c *Counter) Reset() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	prev := c.value
	c.value = 0
	return prev
}

// Snapshot atomically returns the current value and resets the counter to zero.
// This is useful for periodic reporting where you want to capture and clear
// the accumulated count in a single operation without a gap between reads.
func (c *Counter) Snapshot() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	snap := c.value
	c.value = 0
	return snap
}
