package watch

import (
	"testing"
	"time"
)

func makeTokenBucket(cap, rate int) *TokenBucket {
	tb := NewTokenBucket(TokenBucketPolicy{Capacity: cap, RefillRate: rate})
	return tb
}

func TestDefaultTokenBucketPolicyValues(t *testing.T) {
	p := DefaultTokenBucketPolicy()
	if p.Capacity != 10 {
		t.Errorf("expected capacity 10, got %d", p.Capacity)
	}
	if p.RefillRate != 1 {
		t.Errorf("expected refill rate 1, got %d", p.RefillRate)
	}
}

func TestTokenBucketFullOnInit(t *testing.T) {
	tb := makeTokenBucket(5, 1)
	if tb.Tokens() != 5 {
		t.Errorf("expected 5 tokens, got %d", tb.Tokens())
	}
}

func TestTokenBucketAllowConsumesToken(t *testing.T) {
	tb := makeTokenBucket(3, 1)
	if !tb.Allow() {
		t.Fatal("expected Allow to return true")
	}
	if tb.Tokens() != 2 {
		t.Errorf("expected 2 tokens after allow, got %d", tb.Tokens())
	}
}

func TestTokenBucketDeniesWhenEmpty(t *testing.T) {
	tb := makeTokenBucket(2, 1)
	tb.Allow()
	tb.Allow()
	if tb.Allow() {
		t.Error("expected Allow to return false when empty")
	}
}

func TestTokenBucketRefillsOverTime(t *testing.T) {
	tb := makeTokenBucket(5, 2)
	tb.tokens = 0
	past := time.Now().Add(-2 * time.Second)
	tb.lastRefil = past
	if !tb.Allow() {
		t.Error("expected Allow after refill")
	}
}

func TestTokenBucketCapsAtCapacity(t *testing.T) {
	tb := makeTokenBucket(3, 10)
	tb.tokens = 0
	tb.lastRefil = time.Now().Add(-10 * time.Second)
	tb.refill()
	if tb.tokens > 3 {
		t.Errorf("tokens should not exceed capacity, got %d", tb.tokens)
	}
}

func TestNewTokenBucketDefaultsOnZeroPolicy(t *testing.T) {
	tb := NewTokenBucket(TokenBucketPolicy{})
	if tb.policy.Capacity != 10 {
		t.Errorf("expected default capacity 10, got %d", tb.policy.Capacity)
	}
	if tb.policy.RefillRate != 1 {
		t.Errorf("expected default refill rate 1, got %d", tb.policy.RefillRate)
	}
}
