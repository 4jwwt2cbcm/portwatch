package watch

import (
	"testing"
	"time"
)

func makeTickerPool() *TickerPool {
	return NewTickerPool(DefaultTickerPoolPolicy())
}

func TestDefaultTickerPoolPolicyValues(t *testing.T) {
	p := DefaultTickerPoolPolicy()
	if p.MaxTickers != 16 {
		t.Fatalf("expected MaxTickers=16, got %d", p.MaxTickers)
	}
	if p.DefaultInterval != time.Second {
		t.Fatalf("expected DefaultInterval=1s, got %v", p.DefaultInterval)
	}
}

func TestTickerPoolDefaultsOnZero(t *testing.T) {
	p := NewTickerPool(TickerPoolPolicy{})
	if p.policy.MaxTickers != 16 {
		t.Fatalf("expected default MaxTickers=16, got %d", p.policy.MaxTickers)
	}
	if p.policy.DefaultInterval != time.Second {
		t.Fatalf("expected default DefaultInterval=1s, got %v", p.policy.DefaultInterval)
	}
}

func TestTickerPoolAddAndLen(t *testing.T) {
	p := makeTickerPool()
	defer p.StopAll()

	ok := p.Add("scan", 50*time.Millisecond)
	if !ok {
		t.Fatal("expected Add to succeed")
	}
	if p.Len() != 1 {
		t.Fatalf("expected Len=1, got %d", p.Len())
	}
}

func TestTickerPoolDuplicateReturnsFalse(t *testing.T) {
	p := makeTickerPool()
	defer p.StopAll()

	p.Add("scan", 50*time.Millisecond)
	ok := p.Add("scan", 50*time.Millisecond)
	if ok {
		t.Fatal("expected duplicate Add to return false")
	}
	if p.Len() != 1 {
		t.Fatalf("expected Len=1, got %d", p.Len())
	}
}

func TestTickerPoolFullReturnsFalse(t *testing.T) {
	p := NewTickerPool(TickerPoolPolicy{MaxTickers: 2, DefaultInterval: time.Second})
	defer p.StopAll()

	p.Add("a", time.Second)
	p.Add("b", time.Second)
	ok := p.Add("c", time.Second)
	if ok {
		t.Fatal("expected Add to return false when pool is full")
	}
}

func TestTickerPoolCReturnsChannel(t *testing.T) {
	p := makeTickerPool()
	defer p.StopAll()

	p.Add("probe", 20*time.Millisecond)
	ch := p.C("probe")
	if ch == nil {
		t.Fatal("expected non-nil channel for registered ticker")
	}

	select {
	case <-ch:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for ticker tick")
	}
}

func TestTickerPoolCUnknownReturnsNil(t *testing.T) {
	p := makeTickerPool()
	if ch := p.C("missing"); ch != nil {
		t.Fatal("expected nil channel for unknown ticker")
	}
}

func TestTickerPoolRemoveDecreasesLen(t *testing.T) {
	p := makeTickerPool()

	p.Add("x", 50*time.Millisecond)
	p.Remove("x")
	if p.Len() != 0 {
		t.Fatalf("expected Len=0 after Remove, got %d", p.Len())
	}
	if ch := p.C("x"); ch != nil {
		t.Fatal("expected nil channel after Remove")
	}
}

func TestTickerPoolStopAllClearsPool(t *testing.T) {
	p := makeTickerPool()

	p.Add("a", 50*time.Millisecond)
	p.Add("b", 50*time.Millisecond)
	p.StopAll()
	if p.Len() != 0 {
		t.Fatalf("expected Len=0 after StopAll, got %d", p.Len())
	}
}

func TestTickerPoolDefaultIntervalUsedOnZero(t *testing.T) {
	// Just verify Add succeeds when interval=0 (uses DefaultInterval).
	p := makeTickerPool()
	defer p.StopAll()

	ok := p.Add("default-interval", 0)
	if !ok {
		t.Fatal("expected Add with zero interval to succeed")
	}
}
