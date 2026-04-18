package watch_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func TestTimeoutIntegration(t *testing.T) {
	policy := watch.TimeoutPolicy{Duration: 50 * time.Millisecond}

	t.Run("completes within deadline", func(t *testing.T) {
		runner := watch.NewTimeoutRunner(policy, func(ctx context.Context) error {
			return nil
		})
		if err := runner.Run(context.Background()); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("times out slow operation", func(t *testing.T) {
		runner := watch.NewTimeoutRunner(policy, func(ctx context.Context) error {
			time.Sleep(500 * time.Millisecond)
			return nil
		})
		err := runner.Run(context.Background())
		if !errors.Is(err, watch.ErrTimeout) {
			t.Errorf("expected ErrTimeout, got %v", err)
		}
	})
}
