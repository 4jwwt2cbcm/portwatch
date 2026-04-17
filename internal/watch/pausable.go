package watch

import "sync"

// Pausable wraps a run loop so it can be paused and resumed dynamically.
type Pausable struct {
	mu     sync.Mutex
	cond   *sync.Cond
	paused bool
}

// NewPausable returns a ready-to-use Pausable (initially running).
func NewPausable() *Pausable {
	p := &Pausable{}
	p.cond = sync.NewCond(&p.mu)
	return p
}

// Pause suspends callers of Wait until Resume is called.
func (p *Pausable) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = true
}

// Resume unblocks any callers waiting in Wait.
func (p *Pausable) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = false
	p.cond.Broadcast()
}

// IsPaused reports whether the Pausable is currently paused.
func (p *Pausable) IsPaused() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.paused
}

// Wait blocks the caller until the Pausable is not paused.
// It returns immediately if not paused.
func (p *Pausable) Wait() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.paused {
		p.cond.Wait()
	}
}
