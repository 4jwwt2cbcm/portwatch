package watch

import (
	"testing"
	"time"
)

func makeTrigger(minInterval time.Duration) *Trigger {
	return NewTrigger(TriggerPolicy{MinInterval: minInterval})
}

func TestDefaultTriggerPolicyValues(t *testing.T) {
	p := DefaultTriggerPolicy()
	if p.MinInterval != 0 {
		t.Fatalf("expected MinInterval 0, got %v", p.MinInterval)
	}
}

func TestTriggerFireSendsSignal(t *testing.T) {
	tr := makeTrigger(0)
	if !tr.Fire() {
		t.Fatal("expected Fire to return true")
	}
	select {
	case <-tr.C():
	default:
		t.Fatal("expected signal on channel")
	}
}

func TestTriggerSecondFireSuppressedWhilePending(t *testing.T) {
	tr := makeTrigger(0)
	tr.Fire()
	if tr.Fire() {
		t.Fatal("expected second Fire to be suppressed")
	}
}

func TestTriggerCountIncrementsOnFire(t *testing.T) {
	tr := makeTrigger(0)
	tr.Fire()
	<-tr.C()
	tr.Fire()
	<-tr.C()
	if tr.Count() != 2 {
		t.Fatalf("expected count 2, got %d", tr.Count())
	}
}

func TestTriggerMinIntervalSuppressesFastFire(t *testing.T) {
	tr := makeTrigger(500 * time.Millisecond)
	tr.Fire()
	<-tr.C()
	if tr.Fire() {
		t.Fatal("expected Fire to be suppressed within min interval")
	}
}

func TestTriggerAllowedAfterMinInterval(t *testing.T) {
	tr := makeTrigger(10 * time.Millisecond)
	tr.Fire()
	<-tr.C()
	time.Sleep(20 * time.Millisecond)
	if !tr.Fire() {
		t.Fatal("expected Fire to succeed after min interval")
	}
}

func TestTriggerResetClearsState(t *testing.T) {
	tr := makeTrigger(500 * time.Millisecond)
	tr.Fire()
	tr.Reset()
	if tr.Count() != 0 {
		t.Fatalf("expected count 0 after reset, got %d", tr.Count())
	}
	// Should be fireable again immediately after reset
	if !tr.Fire() {
		t.Fatal("expected Fire to succeed after reset")
	}
}
