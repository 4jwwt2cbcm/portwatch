package watch

import (
	"sync"
	"time"
)

// FusePolicy configures the Fuse behaviour.
type FusePolicy struct {
	MaxErrors int
	ResetAfter time.Duration
}

// DefaultFusePolicy returns sensible defaults.
func DefaultFusePolicy() FusePolicy {
	return FusePolicy{
		MaxErrors:  5,
		ResetAfter: 30 * time.Second,
	}
}

// Fuse is a one-shot trip mechanism that blows after MaxErrors failures
// within ResetAfter and can be manually reset.
type Fuse struct {
	mu      sync.Mutex
	policy  FusePolicy
	errors  int
	blown   bool
	lastErr time.Time
}

// NewFuse returns a new Fuse with the given policy.
func NewFuse(p FusePolicy) *Fuse {
	if p.MaxErrors <= 0 {
		p.MaxErrors = DefaultFusePolicy().MaxErrors
	}
	if p.ResetAfter <= 0 {
		p.ResetAfter = DefaultFusePolicy().ResetAfter
	}
	return &Fuse{policy: p}
}

// Record registers an error and trips the fuse if the threshold is reached.
func (f *Fuse) Record() {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := time.Now()
	if !f.lastErr.IsZero() && now.Sub(f.lastErr) > f.policy.ResetAfter {
		f.errors = 0
		f.blown = false
	}
	f.errors++
	f.lastErr = now
	if f.errors >= f.policy.MaxErrors {
		f.blown = true
	}
}

// Blown returns true if the fuse has tripped.
func (f *Fuse) Blown() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.blown
}

// Reset manually resets the fuse.
func (f *Fuse) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.errors = 0
	f.blown = false
	f.lastErr = time.Time{}
}

// Errors returns the current error count.
func (f *Fuse) Errors() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.errors
}
