package watch

import (
	"sync"
	"time"
)

// WindowCounter counts events within a sliding time window.
type WindowCounter struct {
	mu       sync.Mutex
	window   time.Duration
	timestamps []time.Time
}

// NewWindowCounter creates a WindowCounter with the given sliding window duration.
func NewWindowCounter(window time.Duration) *WindowCounter {
	if window <= 0 {
		window = time.Minute
	}
	return &WindowCounter{window: window}
}

// Add records a new event at the current time.
func (w *WindowCounter) Add() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.timestamps = append(w.timestamps, time.Now())
	w.evict()
}

// Count returns the number of events within the current window.
func (w *WindowCounter) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	return len(w.timestamps)
}

// Reset clears all recorded events.
func (w *WindowCounter) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.timestamps = nil
}

// evict removes timestamps outside the current window. Must be called with mu held.
func (w *WindowCounter) evict() {
	cutoff := time.Now().Add(-w.window)
	i := 0
	for i < len(w.timestamps) && w.timestamps[i].Before(cutoff) {
		i++
	}
	w.timestamps = w.timestamps[i:]
}
