package watch_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func TestLeakyBucketIntegration(t *testing.T) {
	p := watch.LeakyBucketPolicy{
		Capacity:  4,
		LeakRate:  2,
		LeakEvery: 20 * time.Millisecond,
	}
	b := watch.NewLeakyBucket(p)

	// Fill to capacity.
	for i := 0; i < 4; i++ {
		if !b.Allow() {
			t.Fatalf("expected allow on fill step %d", i)
		}
	}

	// Bucket full — next should be denied.
	if b.Allow() {
		t.Error("expected deny when bucket is full")
	}

	// Wait for two leak periods (should drain 4 units).
	time.Sleep(50 * time.Millisecond)

	// Bucket should be empty now — allow again.
	if !b.Allow() {
		t.Error("expected allow after bucket drained")
	}

	if b.Level() != 1 {
		t.Errorf("expected level 1, got %d", b.Level())
	}
}
