package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func makeTimeout(d time.Duration) *Timeout {
	return NewTimeout(TimeoutPolicy{Duration: d})
}

func TestDefaultTimeoutPolicyValues(t *testing.T) {
	p := DefaultTimeoutPolicy()
	if p.Duration != 30*time.Second {
		t.Errorf("expected 30s, got %v", p.Duration)
	}
}

func TestTimeoutSucceedsWithinDeadline(t *testing.T) {
	to := makeTimeout(100 * time.Millisecond)
	err := to.Run(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTimeoutReturnsErrTimeoutOnExceed(t *testing.T) {
	to := makeTimeout(20 * time.Millisecond)
	err := to.Run(context.Background(), func(ctx context.Context) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})
	if !errors.Is(err, ErrTimeout) {
		t.Errorf("expected ErrTimeout, got %v", err)
	}
}

func TestTimeoutPropagatesFnError(t *testing.T) {
	to := makeTimeout(100 * time.Millisecond)
	sentinel := errors.New("fn error")
	err := to.Run(context.Background(), func(ctx context.Context) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestTimeoutRespectsParentCancel(t *testing.T) {
	to := makeTimeout(5 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	err := to.Run(ctx, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	if err == nil {
		t.Error("expected error from cancelled context")
	}
}

func TestNewTimeoutDefaultsOnZeroDuration(t *testing.T) {
	to := NewTimeout(TimeoutPolicy{Duration: 0})
	if to.policy.Duration != 30*time.Second {
		t.Errorf("expected default 30s, got %v", to.policy.Duration)
	}
}
