package watch_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func TestJitterIntegration(t *testing.T) {
	policy := watch.DefaultJitterPolicy()
	j := watch.NewJitter(policy)

	base := 200 * time.Millisecond
	results := make([]time.Duration, 50)
	for i := range results {
		results[i] = j.Apply(base)
	}

	// All results must be within [base, base*(1+MaxFraction)].
	max := base + time.Duration(float64(base)*policy.MaxFraction)
	for _, r := range results {
		if r < base || r > max {
			t.Errorf("result %v out of range [%v, %v]", r, base, max)
		}
	}

	// Expect at least some variation across 50 samples.
	allSame := true
	for _, r := range results[1:] {
		if r != results[0] {
			allSame = false
			break
		}
	}
	if allSame {
		t.Error("expected variation in jitter results, all were identical")
	}
}
