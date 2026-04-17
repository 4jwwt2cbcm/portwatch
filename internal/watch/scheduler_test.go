package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func makeScheduler(interval time.Duration, cb func(ctx context.Context) error) *Scheduler {
	policy := SchedulerPolicy{
		InitialDelay: 0,
		Interval:     interval,
		Jitter:       0,
	}
	return NewScheduler(policy, cb)
}

func TestDefaultSchedulerPolicyValues(t *testing.T) {
	p := DefaultSchedulerPolicy()
	if p.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", p.Interval)
	}
	if p.Jitter != 2*time.Second {
		t.Errorf("expected 2s jitter, got %v", p.Jitter)
	}
}

func TestSchedulerRunStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	called := 0
	s := makeScheduler(10*time.Millisecond, func(ctx context.Context) error {
		called++
		if called >= 2 {
			cancel()
		}
		return nil
	})

	err := s.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
	if s.Fires() < 2 {
		t.Errorf("expected at least 2 fires, got %d", s.Fires())
	}
}

func TestSchedulerCallbackErrorStopsRun(t *testing.T) {
	ctx := context.Background()
	sentinel := errors.New("stop")
	s := makeScheduler(time.Millisecond, func(ctx context.Context) error {
		return sentinel
	})
	err := s.Run(ctx)
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
	if s.Fires() != 0 {
		t.Errorf("fires should be 0 on first-call error, got %d", s.Fires())
	}
}

func TestSchedulerInitialDelayRespectsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	policy := SchedulerPolicy{
		InitialDelay: 5 * time.Second,
		Interval:     time.Second,
	}
	s := NewScheduler(policy, func(ctx context.Context) error { return nil })
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err := s.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled during initial delay, got %v", err)
	}
}
