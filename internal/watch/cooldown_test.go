package watch

import (
	"testing"
	"time"
)

func makeCooldown(d time.Duration) *Cooldown {
	c := NewCooldown(CooldownPolicy{Duration: d})
	return c
}

func TestDefaultCooldownPolicyValues(t *testing.T) {
	p := DefaultCooldownPolicy()
	if p.Duration != 30*time.Second {
		t.Errorf("expected 30s, got %v", p.Duration)
	}
}

func TestCooldownFirstCallAllowed(t *testing.T) {
	c := makeCooldown(time.Second)
	if !c.Allow("scan") {
		t.Error("expected first call to be allowed")
	}
}

func TestCooldownSecondCallSuppressed(t *testing.T) {
	c := makeCooldown(time.Second)
	c.Allow("scan")
	if c.Allow("scan") {
		t.Error("expected second call within cooldown to be suppressed")
	}
}

func TestCooldownAllowedAfterExpiry(t *testing.T) {
	c := makeCooldown(10 * time.Millisecond)
	now := time.Now()
	c.now = func() time.Time { return now }
	c.Allow("scan")
	c.now = func() time.Time { return now.Add(20 * time.Millisecond) }
	if !c.Allow("scan") {
		t.Error("expected call after expiry to be allowed")
	}
}

func TestCooldownResetAllowsImmediately(t *testing.T) {
	c := makeCooldown(time.Minute)
	c.Allow("scan")
	c.Reset("scan")
	if !c.Allow("scan") {
		t.Error("expected call after reset to be allowed")
	}
}

func TestCooldownActiveReturnsTrueWhileCooling(t *testing.T) {
	c := makeCooldown(time.Minute)
	c.Allow("scan")
	if !c.Active("scan") {
		t.Error("expected Active to return true during cooldown")
	}
}

func TestCooldownActiveFalseAfterExpiry(t *testing.T) {
	c := makeCooldown(10 * time.Millisecond)
	now := time.Now()
	c.now = func() time.Time { return now }
	c.Allow("scan")
	c.now = func() time.Time { return now.Add(20 * time.Millisecond) }
	if c.Active("scan") {
		t.Error("expected Active to return false after cooldown expiry")
	}
}

func TestCooldownDifferentKeysIndependent(t *testing.T) {
	c := makeCooldown(time.Minute)
	c.Allow("a")
	if !c.Allow("b") {
		t.Error("expected different key to be allowed independently")
	}
}

func TestCooldownZeroDurationDefaulted(t *testing.T) {
	c := NewCooldown(CooldownPolicy{Duration: 0})
	if c.policy.Duration != 30*time.Second {
		t.Errorf("expected default duration, got %v", c.policy.Duration)
	}
}
