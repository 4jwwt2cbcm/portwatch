package metrics

import (
	"testing"
	"time"
)

func TestNewCollectorZeroValues(t *testing.T) {
	c := NewCollector()
	s := c.Snapshot()
	if s.TotalScans != 0 || s.PortsAdded != 0 || s.PortsRemoved != 0 {
		t.Fatalf("expected zero values, got %+v", s)
	}
	if !s.LastScanAt.IsZero() {
		t.Fatalf("expected zero time, got %v", s.LastScanAt)
	}
}

func TestRecordScanIncrementsTotalScans(t *testing.T) {
	c := NewCollector()
	c.RecordScan(10*time.Millisecond, 0, 0)
	c.RecordScan(20*time.Millisecond, 0, 0)
	s := c.Snapshot()
	if s.TotalScans != 2 {
		t.Fatalf("expected TotalScans=2, got %d", s.TotalScans)
	}
}

func TestRecordScanAccumulatesPorts(t *testing.T) {
	c := NewCollector()
	c.RecordScan(5*time.Millisecond, 3, 1)
	c.RecordScan(5*time.Millisecond, 2, 4)
	s := c.Snapshot()
	if s.PortsAdded != 5 {
		t.Fatalf("expected PortsAdded=5, got %d", s.PortsAdded)
	}
	if s.PortsRemoved != 5 {
		t.Fatalf("expected PortsRemoved=5, got %d", s.PortsRemoved)
	}
}

func TestRecordScanSetsLastScanFields(t *testing.T) {
	c := NewCollector()
	before := time.Now()
	c.RecordScan(42*time.Millisecond, 1, 0)
	after := time.Now()
	s := c.Snapshot()
	if s.LastScanAt.Before(before) || s.LastScanAt.After(after) {
		t.Fatalf("LastScanAt %v outside expected range [%v, %v]", s.LastScanAt, before, after)
	}
	if s.LastScanDur != 42*time.Millisecond {
		t.Fatalf("expected LastScanDur=42ms, got %v", s.LastScanDur)
	}
}

func TestSnapshotIsConsistentCopy(t *testing.T) {
	c := NewCollector()
	c.RecordScan(1*time.Millisecond, 1, 0)
	s1 := c.Snapshot()
	c.RecordScan(1*time.Millisecond, 1, 0)
	s2 := c.Snapshot()
	if s1.TotalScans == s2.TotalScans {
		t.Fatalf("expected snapshots to differ after second record")
	}
}

func TestReset(t *testing.T) {
	c := NewCollector()
	c.RecordScan(10*time.Millisecond, 5, 3)
	c.Reset()
	s := c.Snapshot()
	if s.TotalScans != 0 || s.PortsAdded != 0 || s.PortsRemoved != 0 {
		t.Fatalf("expected zero values after reset, got %+v", s)
	}
}
