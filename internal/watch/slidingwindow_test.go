package watch

import (
	"testing"
	"time"
)

func makeSlidingWindow(size int, interval time.Duration) *SlidingWindow {
	sw := NewSlidingWindow(SlidingWindowPolicy{Size: size, Interval: interval})
	return sw
}

func TestDefaultSlidingWindowPolicyValues(t *testing.T) {
	p := DefaultSlidingWindowPolicy()
	if p.Size != 10 {
		t.Errorf("expected size 10, got %d", p.Size)
	}
	if p.Interval != 60*time.Second {
		t.Errorf("expected 60s interval, got %v", p.Interval)
	}
}

func TestSlidingWindowZeroOnInit(t *testing.T) {
	sw := makeSlidingWindow(5, time.Second)
	if sw.Count() != 0 {
		t.Errorf("expected 0, got %d", sw.Count())
	}
}

func TestSlidingWindowRecordReturnsFalseWhenNotFull(t *testing.T) {
	sw := makeSlidingWindow(3, time.Minute)
	if sw.Record() {
		t.Error("expected false on first record")
	}
}

func TestSlidingWindowRecordReturnsTrueWhenFull(t *testing.T) {
	sw := makeSlidingWindow(3, time.Minute)
	sw.Record()
	sw.Record()
	if !sw.Record() {
		t.Error("expected true when window full")
	}
}

func TestSlidingWindowEvictsExpiredEntries(t *testing.T) {
	sw := makeSlidingWindow(5, 100*time.Millisecond)
	base := time.Now()
	sw.now = func() time.Time { return base }
	sw.Record()
	sw.Record()
	// advance time past interval
	sw.now = func() time.Time { return base.Add(200 * time.Millisecond) }
	if sw.Count() != 0 {
		t.Errorf("expected 0 after eviction, got %d", sw.Count())
	}
}

func TestSlidingWindowCountWithinWindow(t *testing.T) {
	sw := makeSlidingWindow(5, time.Minute)
	sw.Record()
	sw.Record()
	if sw.Count() != 2 {
		t.Errorf("expected 2, got %d", sw.Count())
	}
}

func TestSlidingWindowResetClearsCount(t *testing.T) {
	sw := makeSlidingWindow(5, time.Minute)
	sw.Record()
	sw.Record()
	sw.Reset()
	if sw.Count() != 0 {
		t.Errorf("expected 0 after reset, got %d", sw.Count())
	}
}

func TestSlidingWindowDefaultsOnZeroPolicy(t *testing.T) {
	sw := NewSlidingWindow(SlidingWindowPolicy{})
	if sw.policy.Size != 10 {
		t.Errorf("expected default size 10, got %d", sw.policy.Size)
	}
}
