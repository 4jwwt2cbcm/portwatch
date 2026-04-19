package watch

import (
	"context"
	"errors"
)

// ErrBudgetExhausted is returned when the budget limiter denies a call.
var ErrBudgetExhausted = errors.New("budget exhausted")

// BudgetRunner wraps a function and gates execution through a BudgetLimiter.
type BudgetRunner struct {
	limiter *BudgetLimiter
	fn      func(ctx context.Context) error
}

// NewBudgetRunner creates a BudgetRunner. If limiter is nil, a default is used.
func NewBudgetRunner(limiter *BudgetLimiter, fn func(ctx context.Context) error) *BudgetRunner {
	if limiter == nil {
		limiter = NewBudgetLimiter(DefaultBudgetPolicy())
	}
	if fn == nil {
		fn = func(ctx context.Context) error { return nil }
	}
	return &BudgetRunner{limiter: limiter, fn: fn}
}

// Run executes fn if the budget allows, otherwise returns ErrBudgetExhausted.
func (r *BudgetRunner) Run(ctx context.Context) error {
	if !r.limiter.Allow() {
		return ErrBudgetExhausted
	}
	return r.fn(ctx)
}
