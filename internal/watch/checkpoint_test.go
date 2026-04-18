package watch

import (
	"testing"
	"time"
)

func TestCheckpointDueOnInit(t *testing.T) {
	cp := NewCheckpoint(time.Second)
	if !cp.Due() {
		t.Fatal("expected Due() to be true before any Mark")
	}
}

func TestCheckpointNotDueAfterMark(t *testing.T) {
	cp := NewCheckpoint(time.Hour)
	cp.Mark()
	if cp.Due() {
		t.Fatal("expected Due() to be false immediately after Mark")
	}
}

func TestCheckpointDueAfterIntervalElapsed(t *testing.T) {
	cp := NewCheckpoint(10 * time.Millisecond)
	cp.Mark()
	time.Sleep(20 * time.Millisecond)
	if !cp.Due() {
		t.Fatal("expected Due() to be true after interval elapsed")
	}
}

func TestCheckpointLastIsZeroBeforeMark(t *testing.T) {
	cp := NewCheckpoint(time.Second)
	if !cp.Last().IsZero() {
		t.Fatal("expected Last() to be zero before any Mark")
	}
}

func TestCheckpointLastUpdatedAfterMark(t *testing.T) {
	cp := NewCheckpoint(time.Second)
	before := time.Now()
	cp.Mark()
	after := time.Now()
	l := cp.Last()
	if l.Before(before) || l.After(after) {
		t.Fatalf("Last() = %v not between %v and %v", l, before, after)
	}
}

func TestCheckpointResetMakesDueAgain(t *testing.T) {
	cp := NewCheckpoint(time.Hour)
	cp.Mark()
	cp.Reset()
	if !cp.Due() {
		t.Fatal("expected Due() to be true after Reset")
	}
	if !cp.Last().IsZero() {
		t.Fatal("expected Last() to be zero after Reset")
	}
}

func TestCheckpointDefaultsIntervalOnZero(t *testing.T) {
	cp := NewCheckpoint(0)
	if cp.interval != time.Minute {
		t.Fatalf("expected default interval of 1m, got %v", cp.interval)
	}
}
