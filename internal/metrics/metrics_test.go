package metrics

import (
	"testing"
	"time"
)

func TestNewCollectorZeroValues(t *testing.T) {
	c := NewCollector()
	s := c.Snapshot()
	if s.TotalScans != 0 || s.TotalPortsSeen != 0 || s.PortsAdded != 0 || s.PortsRemoved != 0 {
		t.Errorf("expected zero values, got %+v", s)
	}
	if !s.LastScanAt.IsZero() {
		t.Error("expected zero LastScanAt")
	}
}

func TestRecordScanIncrementsTotalScans(t *testing.T) {
	c := NewCollector()
	c.RecordScan(5, 2, 0, time.Millisecond)
	c.RecordScan(5, 0, 1, time.Millisecond)
	s := c.Snapshot()
	if s.TotalScans != 2 {
		t.Errorf("expected TotalScans=2, got %d", s.TotalScans)
	}
}

func TestRecordScanAccumulatesPorts(t *testing.T) {
	c := NewCollector()
	c.RecordScan(3, 3, 0, time.Millisecond)
	c.RecordScan(4, 1, 0, time.Millisecond)
	s := c.Snapshot()
	if s.TotalPortsSeen != 7 {
		t.Errorf("expected TotalPortsSeen=7, got %d", s.TotalPortsSeen)
	}
	if s.PortsAdded != 4 {
		t.Errorf("expected PortsAdded=4, got %d", s.PortsAdded)
	}
}

func TestRecordScanSetsLastScanFields(t *testing.T) {
	c := NewCollector()
	before := time.Now()
	c.RecordScan(7, 0, 2, 42*time.Millisecond)
	s := c.Snapshot()
	if s.LastPortCount != 7 {
		t.Errorf("expected LastPortCount=7, got %d", s.LastPortCount)
	}
	if s.LastScanDuration != 42*time.Millisecond {
		t.Errorf("expected duration 42ms, got %s", s.LastScanDuration)
	}
	if s.LastScanAt.Before(before) {
		t.Error("LastScanAt should be after test start")
	}
}

func TestSnapshotIsConsistentCopy(t *testing.T) {
	c := NewCollector()
	c.RecordScan(2, 1, 0, time.Millisecond)
	s1 := c.Snapshot()
	c.RecordScan(3, 1, 0, time.Millisecond)
	s2 := c.Snapshot()
	if s1.TotalScans == s2.TotalScans {
		t.Error("snapshots should differ after second record")
	}
}

func TestPortsRemovedAccumulates(t *testing.T) {
	c := NewCollector()
	c.RecordScan(5, 0, 2, time.Millisecond)
	c.RecordScan(3, 0, 1, time.Millisecond)
	s := c.Snapshot()
	if s.PortsRemoved != 3 {
		t.Errorf("expected PortsRemoved=3, got %d", s.PortsRemoved)
	}
}
