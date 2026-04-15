package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func makeRetrier(maxAttempts int) *Retrier {
	policy := RetryPolicy{
		MaxAttempts: maxAttempts,
		Backoff: NewBackoff(BackoffPolicy{
			Initial: 1 * time.Millisecond,
			Max:     5 * time.Millisecond,
			Factor:  2.0,
		}),
	}
	return NewRetrier(policy, nil)
}

func TestRetrierSucceedsOnFirstAttempt(t *testing.T) {
	r := makeRetrier(3)
	calls := 0
	err := r.Run(context.Background(), func(_ context.Context) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetrierRetriesOnError(t *testing.T) {
	r := makeRetrier(3)
	calls := 0
	sentinel := errors.New("temporary failure")
	err := r.Run(context.Background(), func(_ context.Context) error {
		calls++
		if calls < 3 {
			return sentinel
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error after retry, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRetrierExhaustsAttempts(t *testing.T) {
	r := makeRetrier(3)
	sentinel := errors.New("persistent error")
	calls := 0
	err := r.Run(context.Background(), func(_ context.Context) error {
		calls++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRetrierRespectsContextCancel(t *testing.T) {
	r := makeRetrier(10)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := r.Run(ctx, func(_ context.Context) error {
		return errors.New("should not retry")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDefaultRetryPolicyValues(t *testing.T) {
	p := DefaultRetryPolicy()
	if p.MaxAttempts != 5 {
		t.Errorf("expected MaxAttempts=5, got %d", p.MaxAttempts)
	}
	if p.Backoff == nil {
		t.Error("expected non-nil Backoff")
	}
}
