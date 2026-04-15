package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time summary of scanner activity.
type Snapshot struct {
	TotalScans   int64
	PortsAdded   int64
	PortsRemoved int64
	LastScanAt   time.Time
	LastScanDur  time.Duration
}

// Collector tracks runtime metrics for portwatch.
type Collector struct {
	mu           sync.RWMutex
	totalScans   int64
	portsAdded   int64
	portsRemoved int64
	lastScanAt   time.Time
	lastScanDur  time.Duration
}

// NewCollector returns an initialised Collector.
func NewCollector() *Collector {
	return &Collector{}
}

// RecordScan records the completion of a single scan cycle.
func (c *Collector) RecordScan(dur time.Duration, added, removed int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.totalScans++
	c.portsAdded += int64(added)
	c.portsRemoved += int64(removed)
	c.lastScanAt = time.Now()
	c.lastScanDur = dur
}

// Snapshot returns a consistent copy of current metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Snapshot{
		TotalScans:   c.totalScans,
		PortsAdded:   c.portsAdded,
		PortsRemoved: c.portsRemoved,
		LastScanAt:   c.lastScanAt,
		LastScanDur:  c.lastScanDur,
	}
}

// Reset clears all counters. Intended for testing.
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	*c = Collector{}
}
