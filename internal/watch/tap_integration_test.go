package watch_test

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/watch"
)

func TestTapRunnerIntegration(t *testing.T) {
	sentinel := errors.New("oops")
	calls := 0
	fn := func() error {
		calls++
		if calls%2 == 0 {
			return sentinel
		}
		return nil
	}

	r := watch.NewTapRunner(fn, 10)

	for i := 0; i < 6; i++ {
		r.Run() //nolint:errcheck
	}

	snap := r.Tap().Snapshot()
	if len(snap) != 6 {
		t.Fatalf("expected 6 recorded results, got %d", len(snap))
	}

	errCount := 0
	for _, e := range snap {
		if e != nil {
			errCount++
		}
	}
	if errCount != 3 {
		t.Fatalf("expected 3 errors, got %d", errCount)
	}
}
