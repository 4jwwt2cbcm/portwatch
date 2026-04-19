package watch

import (
	"testing"
	"time"
)

func makeLeakyBucket(cap, rate int, every time.Duration) *LeakyBucket {
	return NewLeakyBucket(LeakyBucketPolicy{
		Capacity:  cap,
		LeakRate:  rate,
		LeakEvery: every,
	})
}

func TestDefaultLeakyBucketPolicyValues(t *testing.T) {
	p := DefaultLeakyBucketPolicy()
	if p.Capacity != 10 {
		t.Errorf("expected capacity 10, got %d", p.Capacity)
	}
	if p.LeakRate != 1 {
		t.Errorf("expected leak rate 1, got %d", p.LeakRate)
	}
	if p.LeakEvery != time.Second {
		t.Errorf("expected leak every 1s, got %v", p.LeakEvery)
	}
}

func TestLeakyBucketAllowsUpToCapacity(t *testing.T) {
	b := makeLeakyBucket(3, 1, time.Hour)
	for i := 0; i < 3; i++ {
		if !b.Allow() {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
	if b.Allow() {
		t.Error("expected deny when bucket full")
	}
}

func TestLeakyBucketLevelStartsAtZero(t *testing.T) {
	b := makeLeakyBucket(5, 1, time.Hour)
	if b.Level() != 0 {
		t.Errorf("expected level 0, got %d", b.Level())
	}
}

func TestLeakyBucketLeaksOverTime(t *testing.T) {
	b := makeLeakyBucket(5, 2, 10*time.Millisecond)
	b.Allow()
	b.Allow()
	if b.Level() != 2 {
		t.Fatalf("expected level 2, got %d", b.Level())
	}
	time.Sleep(25 * time.Millisecond)
	level := b.Level()
	if level != 0 {
		t.Errorf("expected level 0 after leak, got %d", level)
	}
}

func TestLeakyBucketResetDrains(t *testing.T) {
	b := makeLeakyBucket(5, 1, time.Hour)
	b.Allow()
	b.Allow()
	b.Reset()
	if b.Level() != 0 {
		t.Errorf("expected level 0 after reset, got %d", b.Level())
	}
}

func TestLeakyBucketDefaultsOnZeroPolicy(t *testing.T) {
	b := NewLeakyBucket(LeakyBucketPolicy{})
	if b.policy.Capacity != 10 {
		t.Errorf("expected default capacity 10, got %d", b.policy.Capacity)
	}
}
