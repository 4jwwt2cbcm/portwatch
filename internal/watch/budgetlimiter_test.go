package watch

import (
	"testing"
	"time"
)

func makeBudgetLimiter(max int, window time.Duration) *BudgetLimiter {
	return NewBudgetLimiter(BudgetPolicy{
		Max:      max,
		Window:   window,
		CostFunc: func() int { return 1 },
	})
}

func TestDefaultBudgetPolicyValues(t *testing.T) {
	p := DefaultBudgetPolicy()
	if p.Max != 100 {
		t.Errorf("expected Max=100, got %d", p.Max)
	}
	if p.Window != time.Minute {
		t.Errorf("expected Window=1m, got %v", p.Window)
	}
}

func TestBudgetLimiterAllowsUpToMax(t *testing.T) {
	b := makeBudgetLimiter(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !b.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
	if b.Allow() {
		t.Error("expected Allow()=false after budget exhausted")
	}
}

func TestBudgetLimiterRemaining(t *testing.T) {
	b := makeBudgetLimiter(5, time.Minute)
	b.Allow()
	b.Allow()
	if got := b.Remaining(); got != 3 {
		t.Errorf("expected Remaining=3, got %d", got)
	}
}

func TestBudgetLimiterResetClearsState(t *testing.T) {
	b := makeBudgetLimiter(2, time.Minute)
	b.Allow()
	b.Allow()
	b.Reset()
	if !b.Allow() {
		t.Error("expected Allow()=true after Reset")
	}
}

func TestBudgetLimiterEvictsAfterWindow(t *testing.T) {
	b := NewBudgetLimiter(BudgetPolicy{
		Max:      2,
		Window:   20 * time.Millisecond,
		CostFunc: func() int { return 1 },
	})
	b.Allow()
	b.Allow()
	if b.Allow() {
		t.Error("expected budget exhausted")
	}
	time.Sleep(30 * time.Millisecond)
	if !b.Allow() {
		t.Error("expected Allow()=true after window expired")
	}
}

func TestBudgetLimiterCustomCostFunc(t *testing.T) {
	b := NewBudgetLimiter(BudgetPolicy{
		Max:      10,
		Window:   time.Minute,
		CostFunc: func() int { return 5 },
	})
	if !b.Allow() {
		t.Fatal("expected first Allow")
	}
	if !b.Allow() {
		t.Fatal("expected second Allow")
	}
	if b.Allow() {
		t.Error("expected third Allow to be denied (cost=5, total=10)")
	}
}
