package watch

import (
	"sync"
	"time"
)

// DedupWindow suppresses duplicate events within a sliding time window.
type DedupWindow struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	window  time.Duration
	nowFunc func() time.Time
}

// NewDedupWindow creates a DedupWindow that suppresses repeated keys within window.
func NewDedupWindow(window time.Duration) *DedupWindow {
	return &DedupWindow{
		seen:    make(map[string]time.Time),
		window:  window,
		nowFunc: time.Now,
	}
}

// IsDuplicate returns true if key was seen within the dedup window.
// If not a duplicate, the key is recorded and false is returned.
func (d *DedupWindow) IsDuplicate(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.nowFunc()
	if last, ok := d.seen[key]; ok && now.Sub(last) < d.window {
		return true
	}
	d.seen[key] = now
	return false
}

// Evict removes expired entries from the seen map.
func (d *DedupWindow) Evict() {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.nowFunc()
	for k, t := range d.seen {
		if now.Sub(t) >= d.window {
			delete(d.seen, k)
		}
	}
}

// Len returns the number of keys currently tracked in the dedup window.
func (d *DedupWindow) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}

// Reset clears all seen entries.
func (d *DedupWindow) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
}
