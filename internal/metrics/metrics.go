package metrics

import (
	"sync"
	"time"
)

// Snapshot is an immutable copy of collected metrics at a point in time.
type Snapshot struct {
	TotalScans       int64
	TotalPortsSeen   int64
	PortsAdded       int64
	PortsRemoved     int64
	LastScanAt       time.Time
	LastScanDuration time.Duration
	LastPortCount    int
}

// Collector accumulates runtime metrics for portwatch.
type Collector struct {
	mu               sync.Mutex
	totalScans       int64
	totalPortsSeen   int64
	portsAdded       int64
	portsRemoved     int64
	lastScanAt       time.Time
	lastScanDuration time.Duration
	lastPortCount    int
}

// NewCollector returns a zero-value Collector ready for use.
func NewCollector() *Collector {
	return &Collector{}
}

// RecordScan updates metrics after a single scan cycle completes.
func (c *Collector) RecordScan(portCount, added, removed int, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.totalScans++
	c.totalPortsSeen += int64(portCount)
	c.portsAdded += int64(added)
	c.portsRemoved += int64(removed)
	c.lastScanAt = time.Now()
	c.lastScanDuration = duration
	c.lastPortCount = portCount
}

// Snapshot returns a consistent copy of the current metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Snapshot{
		TotalScans:       c.totalScans,
		TotalPortsSeen:   c.totalPortsSeen,
		PortsAdded:       c.portsAdded,
		PortsRemoved:     c.portsRemoved,
		LastScanAt:       c.lastScanAt,
		LastScanDuration: c.lastScanDuration,
		LastPortCount:    c.lastPortCount,
	}
}

// Reset clears all accumulated metrics back to their zero values.
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.totalScans = 0
	c.totalPortsSeen = 0
	c.portsAdded = 0
	c.portsRemoved = 0
	c.lastScanAt = time.Time{}
	c.lastScanDuration = 0
	c.lastPortCount = 0
}
