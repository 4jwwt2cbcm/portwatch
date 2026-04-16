package watch

import (
	"testing"
)

func TestEventLogEmptyOnInit(t *testing.T) {
	l := NewEventLog(10)
	if l.Len() != 0 {
		t.Fatalf("expected 0 events, got %d", l.Len())
	}
}

func TestEventLogAppendAndSnapshot(t *testing.T) {
	l := NewEventLog(10)
	l.Append(EventScanOK, "all good")
	l.Append(EventPortAdded, "port 80")

	snap := l.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 events, got %d", len(snap))
	}
	if snap[0].Kind != EventScanOK {
		t.Errorf("expected scan_ok, got %s", snap[0].Kind)
	}
	if snap[1].Message != "port 80" {
		t.Errorf("unexpected message: %s", snap[1].Message)
	}
}

func TestEventLogEvictsOldestWhenFull(t *testing.T) {
	l := NewEventLog(3)
	l.Append(EventScanOK, "first")
	l.Append(EventScanOK, "second")
	l.Append(EventScanOK, "third")
	l.Append(EventPortAdded, "fourth")

	snap := l.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 events, got %d", len(snap))
	}
	if snap[0].Message != "second" {
		t.Errorf("expected oldest evicted, got %s", snap[0].Message)
	}
	if snap[2].Message != "fourth" {
		t.Errorf("expected newest last, got %s", snap[2].Message)
	}
}

func TestEventLogClear(t *testing.T) {
	l := NewEventLog(10)
	l.Append(EventScanError, "boom")
	l.Clear()
	if l.Len() != 0 {
		t.Fatalf("expected empty after clear, got %d", l.Len())
	}
}

func TestEventLogDefaultCapOnZero(t *testing.T) {
	l := NewEventLog(0)
	if l.cap != 100 {
		t.Errorf("expected default cap 100, got %d", l.cap)
	}
}

func TestSnapshotIsIndependentCopy(t *testing.T) {
	l := NewEventLog(10)
	l.Append(EventScanOK, "a")
	snap := l.Snapshot()
	l.Append(EventPortGone, "b")
	if len(snap) != 1 {
		t.Errorf("snapshot should not reflect later appends")
	}
}
