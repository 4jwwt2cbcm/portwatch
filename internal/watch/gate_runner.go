package watch

import "context"

// GateRunner wraps a function and only executes it when a Gate is open.
// If the gate is closed the call is skipped and no error is returned.
type GateRunner struct {
	gate *Gate
	fn   func(ctx context.Context) error
}

// NewGateRunner returns a GateRunner that guards fn with gate.
func NewGateRunner(gate *Gate, fn func(ctx context.Context) error) *GateRunner {
	if gate == nil {
		gate = NewGate(true)
	}
	return &GateRunner{gate: gate, fn: fn}
}

// Run executes the wrapped function only if the gate is open.
// Returns nil without calling fn when the gate is closed.
func (r *GateRunner) Run(ctx context.Context) error {
	if !r.gate.Allow() {
		return nil
	}
	return r.fn(ctx)
}
