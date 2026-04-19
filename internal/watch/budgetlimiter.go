package watch

import (
	"sync"
	"time"
)

// BudgetPolicy configures the BudgetLimiter.
type BudgetPolicy struct {
	Max      int
	Window   time.Duration
	CostFunc func() int
}

// DefaultBudgetPolicy returns sensible defaults.
func DefaultBudgetPolicy() BudgetPolicy {
	return BudgetPolicy{
		Max:      100,
		Window:   time.Minute,
		CostFunc: func() int { return 1 },
	}
}

// BudgetLimiter tracks a rolling cost budget over a time window.
type BudgetLimiter struct {
	mu      sync.Mutex
	policy  BudgetPolicy
	entries []budgetEntry
}

type budgetEntry struct {
	at   time.Time
	cost int
}

// NewBudgetLimiter creates a BudgetLimiter with the given policy.
func NewBudgetLimiter(p BudgetPolicy) *BudgetLimiter {
	if p.Max <= 0 {
		p.Max = DefaultBudgetPolicy().Max
	}
	if p.Window <= 0 {
		p.Window = DefaultBudgetPolicy().Window
	}
	if p.CostFunc == nil {
		p.CostFunc = DefaultBudgetPolicy().CostFunc
	}
	return &BudgetLimiter{policy: p}
}

// Allow returns true if the current cost fits within the budget.
func (b *BudgetLimiter) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	b.evict(now)
	cost := b.policy.CostFunc()
	if b.total()+cost > b.policy.Max {
		return false
	}
	b.entries = append(b.entries, budgetEntry{at: now, cost: cost})
	return true
}

// Remaining returns the remaining budget.
func (b *BudgetLimiter) Remaining() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.evict(time.Now())
	r := b.policy.Max - b.total()
	if r < 0 {
		return 0
	}
	return r
}

// Reset clears all recorded entries.
func (b *BudgetLimiter) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = nil
}

func (b *BudgetLimiter) evict(now time.Time) {
	cutoff := now.Add(-b.policy.Window)
	i := 0
	for i < len(b.entries) && b.entries[i].at.Before(cutoff) {
		i++
	}
	b.entries = b.entries[i:]
}

func (b *BudgetLimiter) total() int {
	t := 0
	for _, e := range b.entries {
		t += e.cost
	}
	return t
}
