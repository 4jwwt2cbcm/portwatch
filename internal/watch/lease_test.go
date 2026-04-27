package watch

import (
	"testing"
	"time"
)

func makeLease(ttl time.Duration) *Lease {
	l := NewLease(LeasePolicy{TTL: ttl, RenewAt: 0.5})
	return l
}

func TestDefaultLeasePolicyValues(t *testing.T) {
	p := DefaultLeasePolicy()
	if p.TTL != 30*time.Second {
		t.Errorf("expected 30s TTL, got %v", p.TTL)
	}
	if p.RenewAt != 0.75 {
		t.Errorf("expected RenewAt=0.75, got %v", p.RenewAt)
	}
}

func TestLeaseNotHeldOnInit(t *testing.T) {
	l := makeLease(time.Second)
	if l.Held() {
		t.Error("expected lease not held on init")
	}
}

func TestLeaseAcquireSucceeds(t *testing.T) {
	l := makeLease(time.Second)
	if !l.Acquire() {
		t.Error("expected first Acquire to succeed")
	}
	if !l.Held() {
		t.Error("expected lease to be held after Acquire")
	}
}

func TestLeaseAcquireBlocksWhenHeld(t *testing.T) {
	l := makeLease(time.Second)
	l.Acquire()
	if l.Acquire() {
		t.Error("expected second Acquire to fail while held")
	}
}

func TestLeaseAcquireSucceedsAfterExpiry(t *testing.T) {
	l := makeLease(10 * time.Millisecond)
	l.Acquire()
	time.Sleep(20 * time.Millisecond)
	if !l.Acquire() {
		t.Error("expected Acquire to succeed after expiry")
	}
}

func TestLeaseRenewExtendsTTL(t *testing.T) {
	l := makeLease(100 * time.Millisecond)
	l.Acquire()
	first := l.ExpiresAt()
	time.Sleep(10 * time.Millisecond)
	if !l.Renew() {
		t.Error("expected Renew to succeed")
	}
	if !l.ExpiresAt().After(first) {
		t.Error("expected expiry to be extended after Renew")
	}
}

func TestLeaseRenewFailsWhenNotHeld(t *testing.T) {
	l := makeLease(time.Second)
	if l.Renew() {
		t.Error("expected Renew to fail when not held")
	}
}

func TestLeaseReleaseRelinquishesLease(t *testing.T) {
	l := makeLease(time.Second)
	l.Acquire()
	l.Release()
	if l.Held() {
		t.Error("expected lease not held after Release")
	}
	if !l.Acquire() {
		t.Error("expected Acquire to succeed after Release")
	}
}

func TestLeaseShouldRenewFalseWhenNotHeld(t *testing.T) {
	l := makeLease(time.Second)
	if l.ShouldRenew() {
		t.Error("expected ShouldRenew false when not held")
	}
}

func TestLeaseShouldRenewTrueNearExpiry(t *testing.T) {
	l := NewLease(LeasePolicy{TTL: 40 * time.Millisecond, RenewAt: 0.5})
	l.Acquire()
	time.Sleep(25 * time.Millisecond) // past 50% threshold
	if !l.ShouldRenew() {
		t.Error("expected ShouldRenew true near expiry")
	}
}
