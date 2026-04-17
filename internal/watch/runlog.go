package watch

import (
	"sync"
	"time"
)

// RunEntry records the outcome of a single watcher cycle.
type RunEntry struct {
	StartedAt  time.Time
	FinishedAt time.Time
	Err        error
	PortsFound int
}

// Duration returns how long the run took.
func (e RunEntry) Duration() time.Duration {
	return e.FinishedAt.Sub(e.StartedAt)
}

// RunLog keeps a bounded history of recent run entries.
type RunLog struct {
	mu      sync.Mutex
	entries []RunEntry
	cap     int
}

// NewRunLog creates a RunLog with the given capacity. Zero defaults to 64.
func NewRunLog(cap int) *RunLog {
	if cap <= 0 {
		cap = 64
	}
	return &RunLog{cap: cap, entries: make([]RunEntry, 0, cap)}
}

// Append adds a new entry, evicting the oldest if at capacity.
func (r *RunLog) Append(e RunEntry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.entries) >= r.cap {
		r.entries = r.entries[1:]
	}
	r.entries = append(r.entries, e)
}

// Snapshot returns a copy of all current entries.
func (r *RunLog) Snapshot() []RunEntry {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]RunEntry, len(r.entries))
	copy(out, r.entries)
	return out
}

// Len returns the current number of entries.
func (r *RunLog) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries)
}

// Clear removes all entries.
func (r *RunLog) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = r.entries[:0]
}
