package watch

import (
	"context"
	"sync"
	"time"
)

// SchedulerPolicy controls scheduling behaviour.
type SchedulerPolicy struct {
	InitialDelay time.Duration
	Interval     time.Duration
	Jitter       time.Duration
}

// DefaultSchedulerPolicy returns sensible defaults.
func DefaultSchedulerPolicy() SchedulerPolicy {
	return SchedulerPolicy{
		InitialDelay: 0,
		Interval:     30 * time.Second,
		Jitter:       2 * time.Second,
	}
}

// Scheduler fires a callback on a configurable interval with optional jitter.
type Scheduler struct {
	policy  SchedulerPolicy
	callback func(ctx context.Context) error
	mu      sync.Mutex
	fires   int
}

// NewScheduler creates a Scheduler with the given policy and callback.
func NewScheduler(policy SchedulerPolicy, callback func(ctx context.Context) error) *Scheduler {
	return &Scheduler{policy: policy, callback: callback}
}

// Fires returns the number of times the callback has been invoked.
func (s *Scheduler) Fires() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.fires
}

// Run starts the scheduler loop, blocking until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) error {
	if s.policy.InitialDelay > 0 {
		select {
		case <-time.After(s.policy.InitialDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	for {
		if err := s.callback(ctx); err != nil {
			return err
		}
		s.mu.Lock()
		s.fires++
		s.mu.Unlock()

		interval := s.policy.Interval
		if s.policy.Jitter > 0 {
			interval += time.Duration(time.Now().UnixNano()%int64(s.policy.Jitter))
		}

		select {
		case <-time.After(interval):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
