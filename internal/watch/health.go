package watch

import (
	"sync"
	"time"
)

// HealthStatus represents the current health of the watcher loop.
type HealthStatus struct {
	Healthy      bool
	LastSuccess  time.Time
	LastError    error
	ConsecErrors int
}

// HealthTracker tracks the liveness of the watch loop and exposes a status
// snapshot that can be queried by an operator or a future HTTP endpoint.
type HealthTracker struct {
	mu           sync.RWMutex
	lastSuccess  time.Time
	lastError    error
	consecErrors int
	maxErrors    int
}

// NewHealthTracker returns a HealthTracker that considers the loop unhealthy
// after maxErrors consecutive failures.
func NewHealthTracker(maxErrors int) *HealthTracker {
	if maxErrors <= 0 {
		maxErrors = 3
	}
	return &HealthTracker{maxErrors: maxErrors}
}

// RecordSuccess resets the consecutive error counter and stamps the last
// successful cycle time.
func (h *HealthTracker) RecordSuccess() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastSuccess = time.Now()
	h.consecErrors = 0
	h.lastError = nil
}

// RecordError increments the consecutive error counter and stores the error.
func (h *HealthTracker) RecordError(err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastError = err
	h.consecErrors++
}

// Status returns a point-in-time snapshot of the current health.
func (h *HealthTracker) Status() HealthStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return HealthStatus{
		Healthy:      h.consecErrors < h.maxErrors,
		LastSuccess:  h.lastSuccess,
		LastError:    h.lastError,
		ConsecErrors: h.consecErrors,
	}
}
