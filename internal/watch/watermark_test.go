package watch

import (
	"testing"
	"time"
)

func makeWatermark() *Watermark {
	return NewWatermark(DefaultWatermarkPolicy())
}

func TestDefaultWatermarkPolicyValues(t *testing.T) {
	p := DefaultWatermarkPolicy()
	if p.HighWater != 0.80 {
		t.Errorf("expected HighWater 0.80, got %v", p.HighWater)
	}
	if p.LowWater != 0.50 {
		t.Errorf("expected LowWater 0.50, got %v", p.LowWater)
	}
}

func TestWatermarkNotAboveOnInit(t *testing.T) {
	w := makeWatermark()
	if w.Above() {
		t.Error("expected Above() == false on init")
	}
}

func TestWatermarkLevelZeroOnInit(t *testing.T) {
	w := makeWatermark()
	if w.Level() != 0 {
		t.Errorf("expected level 0, got %v", w.Level())
	}
}

func TestWatermarkAboveAfterHighWaterCrossed(t *testing.T) {
	w := makeWatermark()
	w.Set(0.85)
	if !w.Above() {
		t.Error("expected Above() == true after crossing high watermark")
	}
}

func TestWatermarkRemainsAboveUntilLowWaterCrossed(t *testing.T) {
	w := makeWatermark()
	w.Set(0.90)
	w.Set(0.65) // between low and high — should stay above
	if !w.Above() {
		t.Error("expected Above() == true while level is between watermarks")
	}
}

func TestWatermarkDropsBelowAfterLowWaterCrossed(t *testing.T) {
	w := makeWatermark()
	w.Set(0.90)
	w.Set(0.40) // below low watermark
	if w.Above() {
		t.Error("expected Above() == false after dropping below low watermark")
	}
}

func TestWatermarkLastUpdatedSetOnSet(t *testing.T) {
	w := makeWatermark()
	before := time.Now()
	w.Set(0.5)
	if w.LastUpdated().Before(before) {
		t.Error("expected LastUpdated to be after Set call")
	}
}

func TestWatermarkResetClearsState(t *testing.T) {
	w := makeWatermark()
	w.Set(0.95)
	w.Reset()
	if w.Above() {
		t.Error("expected Above() == false after Reset")
	}
	if w.Level() != 0 {
		t.Errorf("expected level 0 after Reset, got %v", w.Level())
	}
	if !w.LastUpdated().IsZero() {
		t.Error("expected LastUpdated to be zero after Reset")
	}
}

func TestWatermarkDefaultsOnZeroPolicy(t *testing.T) {
	w := NewWatermark(WatermarkPolicy{})
	w.Set(0.85)
	if !w.Above() {
		t.Error("expected defaults applied when policy is zero")
	}
}
