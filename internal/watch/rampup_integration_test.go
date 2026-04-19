package watch_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func TestRampUpIntegration(t *testing.T) {
	p := watch.RampUpPolicy{
		Steps:   3,
		Initial: 10 * time.Millisecond,
		Target:  30 * time.Millisecond,
	}
	r := watch.NewRampUp(p)

	var intervals []time.Duration
	for !r.Done() {
		intervals = append(intervals, r.Next())
	}
	// consume final step
	intervals = append(intervals, r.Next())

	if len(intervals) == 0 {
		t.Fatal("expected at least one interval")
	}
	if intervals[0] != 10*time.Millisecond {
		t.Errorf("first interval should be initial, got %v", intervals[0])
	}
	last := intervals[len(intervals)-1]
	if last != 30*time.Millisecond {
		t.Errorf("last interval should be target 30ms, got %v", last)
	}
	for i := 1; i < len(intervals); i++ {
		if intervals[i] < intervals[i-1] {
			t.Errorf("intervals not monotonically increasing at index %d", i)
		}
	}
}
