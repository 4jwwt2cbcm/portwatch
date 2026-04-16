package watch

import (
	"testing"
	"time"
)

func makeRateLimiter(rate int, period time.Duration) *RateLimiter {
	rl := NewRateLimiter(RateLimitPolicy{Rate: rate, Period: period})
	return rl
}

func TestDefaultRateLimitPolicyValues(t *testing.T) {
	p := DefaultRateLimitPolicy()
	if p.Rate != 10 {
		t.Errorf("expected rate 10, got %d", p.Rate)
	}
	if p.Period != time.Minute {
		t.Errorf("expected period 1m, got %v", p.Period)
	}
}

func TestRateLimiterAllowsUpToRate(t *testing.T) {
	rl := makeRateLimiter(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !rl.Allow() {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
}

func TestRateLimiterBlocksWhenExhausted(t *testing.T) {
	rl := makeRateLimiter(2, time.Minute)
	rl.Allow()
	rl.Allow()
	if rl.Allow() {
		t.Error("expected deny after tokens exhausted")
	}
}

func TestRateLimiterRefillsAfterPeriod(t *testing.T) {
	base := time.Now()
	rl := makeRateLimiter(2, time.Second)
	rl.now = func() time.Time { return base }
	rl.Allow()
	rl.Allow()

	rl.now = func() time.Time { return base.Add(2 * time.Second) }
	if !rl.Allow() {
		t.Error("expected allow after period elapsed")
	}
}

func TestRateLimiterResetRestoresTokens(t *testing.T) {
	rl := makeRateLimiter(2, time.Minute)
	rl.Allow()
	rl.Allow()
	rl.Reset()
	if !rl.Allow() {
		t.Error("expected allow after reset")
	}
}

func TestRateLimiterDoesNotExceedMaxTokens(t *testing.T) {
	base := time.Now()
	rl := makeRateLimiter(3, time.Second)
	rl.now = func() time.Time { return base.Add(100 * time.Second) }
	// trigger a refill by calling Allow
	rl.Allow()
	if rl.tokens > 3 {
		t.Errorf("tokens %d exceed max rate %d", rl.tokens, 3)
	}
}
