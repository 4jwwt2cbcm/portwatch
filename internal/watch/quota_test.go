package watch

import (
	"testing"
	"time"
)

func makeQuota(max int, window time.Duration) *Quota {
	return NewQuota(QuotaPolicy{Max: max, Window: window})
}

func TestDefaultQuotaPolicyValues(t *testing.T) {
	p := DefaultQuotaPolicy()
	if p.Max != 100 {
		t.Errorf("expected Max=100, got %d", p.Max)
	}
	if p.Window != time.Minute {
		t.Errorf("expected Window=1m, got %v", p.Window)
	}
}

func TestQuotaAllowsUpToMax(t *testing.T) {
	q := makeQuota(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !q.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
	if q.Allow() {
		t.Error("expected Allow()=false after max reached")
	}
}

func TestQuotaRemainingDecrementsOnAllow(t *testing.T) {
	q := makeQuota(5, time.Minute)
	if q.Remaining() != 5 {
		t.Fatalf("expected Remaining=5, got %d", q.Remaining())
	}
	q.Allow()
	q.Allow()
	if q.Remaining() != 3 {
		t.Errorf("expected Remaining=3, got %d", q.Remaining())
	}
}

func TestQuotaResetClearsUsage(t *testing.T) {
	q := makeQuota(2, time.Minute)
	q.Allow()
	q.Allow()
	if q.Allow() {
		t.Fatal("expected quota exhausted before reset")
	}
	q.Reset()
	if !q.Allow() {
		t.Error("expected Allow()=true after reset")
	}
}

func TestQuotaWindowExpiryClearsUsage(t *testing.T) {
	q := makeQuota(2, time.Millisecond)
	q.Allow()
	q.Allow()
	if q.Allow() {
		t.Fatal("expected quota exhausted")
	}
	time.Sleep(5 * time.Millisecond)
	if !q.Allow() {
		t.Error("expected Allow()=true after window expiry")
	}
}

func TestQuotaDefaultsOnZeroPolicy(t *testing.T) {
	q := NewQuota(QuotaPolicy{})
	if q.policy.Max != 100 {
		t.Errorf("expected default Max=100, got %d", q.policy.Max)
	}
	if q.policy.Window != time.Minute {
		t.Errorf("expected default Window=1m, got %v", q.policy.Window)
	}
}
