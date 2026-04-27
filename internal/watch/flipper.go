package watch

import "sync"

// Flipper is a boolean toggle that can be flipped atomically.
// It tracks how many times it has been toggled and supports
// an optional callback fired on each flip.
type Flipper struct {
	mu       sync.Mutex
	state    bool
	count    int
	onFlip   func(bool)
}

// NewFlipper returns a Flipper with the given initial state.
// If onFlip is non-nil it is called with the new state after each flip.
func NewFlipper(initial bool, onFlip func(bool)) *Flipper {
	return &Flipper{
		state:  initial,
		onFlip: onFlip,
	}
}

// Flip toggles the current state and increments the flip counter.
// Returns the new state.
func (f *Flipper) Flip() bool {
	f.mu.Lock()
	f.state = !f.state
	newState := f.state
	f.count++
	cb := f.onFlip
	f.mu.Unlock()

	if cb != nil {
		cb(newState)
	}
	return newState
}

// State returns the current boolean state without modifying it.
func (f *Flipper) State() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.state
}

// Count returns the total number of times Flip has been called.
func (f *Flipper) Count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.count
}

// Reset sets the state to the given value and clears the flip counter.
func (f *Flipper) Reset(state bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.state = state
	f.count = 0
}
