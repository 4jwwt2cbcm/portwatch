package watch

import (
	"context"
	"testing"
	"time"
)

func makeTokenPool(cap, refillAmount int, refillEvery time.Duration) *TokenPool {
	return NewTokenPool(TokenPolicy{
		Capacity:     cap,
		RefillEvery:  refillEvery,
		RefillAmount: refillAmount,
	})
}

func TestDefaultTokenPolicyValues(t *testing.T) {
	p := DefaultTokenPolicy()
	if p.Capacity != 10 {
		t.Errorf("expected capacity 10, got %d", p.Capacity)
	}
	if p.RefillEvery != time.Second {
		t.Errorf("expected refill every 1s, got %v", p.RefillEvery)
	}
	if p.RefillAmount != 1 {
		t.Errorf("expected refill amount 1, got %d", p.RefillAmount)
	}
}

func TestTokenPoolFullOnInit(t *testing.T) {
	tp := makeTokenPool(5, 1, time.Second)
	if got := tp.Available(); got != 5 {
		t.Errorf("expected 5 available, got %d", got)
	}
}

func TestTokenPoolTakeReducesAvailable(t *testing.T) {
	tp := makeTokenPool(3, 1, time.Second)
	if !tp.Take() {
		t.Fatal("expected Take to succeed")
	}
	if got := tp.Available(); got != 2 {
		t.Errorf("expected 2 available, got %d", got)
	}
}

func TestTokenPoolTakeReturnsFalseWhenEmpty(t *testing.T) {
	tp := makeTokenPool(1, 1, time.Second)
	tp.Take()
	if tp.Take() {
		t.Error("expected Take to fail on empty pool")
	}
}

func TestTokenPoolRefillCapsAtCapacity(t *testing.T) {
	tp := makeTokenPool(3, 10, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	go tp.Run(ctx) //nolint:errcheck
	time.Sleep(50 * time.Millisecond)
	if got := tp.Available(); got > 3 {
		t.Errorf("expected available <= 3, got %d", got)
	}
}

func TestTokenPoolRunStopsOnContextCancel(t *testing.T) {
	tp := makeTokenPool(5, 1, 10*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := tp.Run(ctx)
	if err == nil {
		t.Error("expected error on cancelled context")
	}
}

func TestTokenPoolDefaultsOnZeroPolicy(t *testing.T) {
	tp := NewTokenPool(TokenPolicy{})
	if tp.policy.Capacity != 10 {
		t.Errorf("expected default capacity 10, got %d", tp.policy.Capacity)
	}
}
