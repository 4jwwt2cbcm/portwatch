package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestTimeoutRunnerSucceeds(t *testing.T) {
	r := NewTimeoutRunner(TimeoutPolicy{Duration: 100 * time.Millisecond}, func(ctx context.Context) error {
		return nil
	})
	if err := r.Run(context.Background()); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTimeoutRunnerTimesOut(t *testing.T) {
	r := NewTimeoutRunner(TimeoutPolicy{Duration: 20 * time.Millisecond}, func(ctx context.Context) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})
	err := r.Run(context.Background())
	if !errors.Is(err, ErrTimeout) {
		t.Errorf("expected ErrTimeout, got %v", err)
	}
}

func TestTimeoutRunnerPropagatesError(t *testing.T) {
	sentinel := errors.New("inner")
	r := NewTimeoutRunner(TimeoutPolicy{Duration: 100 * time.Millisecond}, func(ctx context.Context) error {
		return sentinel
	})
	err := r.Run(context.Background())
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel, got %v", err)
	}
}

func TestTimeoutRunnerDefaultsOnZero(t *testing.T) {
	r := NewTimeoutRunner(TimeoutPolicy{}, func(ctx context.Context) error { return nil })
	if r.timeout.policy.Duration != 30*time.Second {
		t.Errorf("expected 30s default, got %v", r.timeout.policy.Duration)
	}
}
