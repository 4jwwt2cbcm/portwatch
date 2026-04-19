package watch

import (
	"context"
	"errors"
	"sync"
)

// Cascade runs a chain of functions sequentially, stopping on the first error.
// Each step receives the context and may cancel downstream steps.
type Cascade struct {
	mu    sync.Mutex
	steps []CascadeStep
}

// CascadeStep is a named function in a Cascade chain.
type CascadeStep struct {
	Name string
	Fn   func(ctx context.Context) error
}

// ErrCascadeAborted is returned when a step fails and the chain is halted.
var ErrCascadeAborted = errors.New("cascade aborted")

// NewCascade returns an empty Cascade.
func NewCascade() *Cascade {
	return &Cascade{}
}

// Add appends a named step to the cascade.
func (c *Cascade) Add(name string, fn func(ctx context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.steps = append(c.steps, CascadeStep{Name: name, Fn: fn})
}

// Run executes each step in order. If any step returns an error, Run stops
// and returns that error wrapped with ErrCascadeAborted.
func (c *Cascade) Run(ctx context.Context) error {
	c.mu.Lock()
	steps := make([]CascadeStep, len(c.steps))
	copy(steps, c.steps)
	c.mu.Unlock()

	for _, s := range steps {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := s.Fn(ctx); err != nil {
			return errors.Join(ErrCascadeAborted, err)
		}
	}
	return nil
}

// Len returns the number of registered steps.
func (c *Cascade) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.steps)
}

// Clear removes all steps.
func (c *Cascade) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.steps = c.steps[:0]
}
