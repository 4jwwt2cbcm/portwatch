package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestStaggerRunnerRunsFunction(t *testing.T) {
	s := NewStagger(StaggerPolicy{Delay: 1 * time.Millisecond, MaxItems: 5})
	called := false
	r := NewStaggerRunner(s, func(_ context.Context) error {
		called = true
		return nil
	})
	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected fn to be called")
	}
}

func TestStaggerRunnerPropagatesError(t *testing.T) {
	s := NewStagger(StaggerPolicy{Delay: 1 * time.Millisecond, MaxItems: 5})
	sentinel := errors.New("boom")
	r := NewStaggerRunner(s, func(_ context.Context) error { return sentinel })
	if err := r.Run(context.Background()); !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestStaggerRunnerContextCancelReturnsErr(t *testing.T) {
	s := NewStagger(StaggerPolicy{Delay: 10 * time.Second, MaxItems: 5})
	// Consume the first slot so the next call will block for 10 s
	s.Next()
	r := NewStaggerRunner(s, func(_ context.Context) error { return nil })
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	if err := r.Run(ctx); !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestStaggerRunnerNilArgDefaults(t *testing.T) {
	r := NewStaggerRunner(nil, nil)
	if r.stagger == nil {
		t.Error("expected non-nil stagger after defaulting")
	}
	if err := r.Run(context.Background()); err != nil {
		t.Errorf("nil fn should be a no-op, got %v", err)
	}
}

func TestStaggerRunnerIncrementsStaggerCount(t *testing.T) {
	s := NewStagger(StaggerPolicy{Delay: 1 * time.Millisecond, MaxItems: 10})
	r := NewStaggerRunner(s, func(_ context.Context) error { return nil })
	for i := 0; i < 3; i++ {
		if err := r.Run(context.Background()); err != nil {
			t.Fatalf("run %d failed: %v", i, err)
		}
	}
	if s.Count() != 3 {
		t.Errorf("expected stagger count=3, got %d", s.Count())
	}
}
