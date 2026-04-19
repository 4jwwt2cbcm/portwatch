package watch

import (
	"testing"
	"time"
)

func makePressureTracker() *PressureTracker {
	return NewPressureTracker(DefaultPressurePolicy())
}

func TestDefaultPressurePolicyValues(t *testing.T) {
	p := DefaultPressurePolicy()
	if p.HighWatermark != 0.80 {
		t.Fatalf("expected HighWatermark 0.80, got %v", p.HighWatermark)
	}
	if p.LowWatermark != 0.50 {
		t.Fatalf("expected LowWatermark 0.50, got %v", p.LowWatermark)
	}
	if p.Window != 30*time.Second {
		t.Fatalf("expected Window 30s, got %v", p.Window)
	}
}

func TestPressureNotHighOnInit(t *testing.T) {
	pt := makePressureTracker()
	if pt.High() {
		t.Fatal("expected not high on init")
	}
}

func TestPressureAverageZeroOnInit(t *testing.T) {
	pt := makePressureTracker()
	if pt.Average() != 0 {
		t.Fatalf("expected 0 average, got %v", pt.Average())
	}
}

func TestPressureHighAfterHighLoad(t *testing.T) {
	pt := makePressureTracker()
	pt.Record(0.9)
	pt.Record(0.95)
	if !pt.High() {
		t.Fatal("expected high pressure after high load samples")
	}
}

func TestPressureNotHighBelowWatermark(t *testing.T) {
	pt := makePressureTracker()
	pt.Record(0.3)
	pt.Record(0.4)
	if pt.High() {
		t.Fatal("expected not high below watermark")
	}
}

func TestPressureRelievesAfterLowLoad(t *testing.T) {
	pt := NewPressureTracker(PressurePolicy{
		HighWatermark: 0.80,
		LowWatermark:  0.50,
		Window:        30 * time.Second,
	})
	pt.Record(0.9)
	if !pt.High() {
		t.Fatal("expected high after 0.9 sample")
	}
	// replace with low samples
	pt.Record(0.2)
	pt.Record(0.2)
	pt.Record(0.2)
	if pt.High() {
		t.Fatal("expected pressure relieved after low samples")
	}
}

func TestPressureDefaultsWindowOnZero(t *testing.T) {
	pt := NewPressureTracker(PressurePolicy{})
	if pt.policy.Window != DefaultPressurePolicy().Window {
		t.Fatalf("expected default window, got %v", pt.policy.Window)
	}
}
