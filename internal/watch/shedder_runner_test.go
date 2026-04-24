package watch

import (
	"context"
	"errors"
	"testing"
)

func TestShedderRunnerRunsWhenBelowThreshold(t *testing.T) {
	s := makeShedder()
	s.Record(0.1)

	called := false
	r := NewShedderRunner(s, func(_ context.Context) error {
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

func TestShedderRunnerShedsAboveThreshold(t *testing.T) {
	s := makeShedder()
	s.Record(0.9)
	s.Record(0.9)

	r := NewShedderRunner(s, func(_ context.Context) error {
		t.Error("fn should not be called when shedding")
		return nil
	})
	err := r.Run(context.Background())
	if !errors.Is(err, ErrShed) {
		t.Errorf("expected ErrShed, got %v", err)
	}
}

func TestShedderRunnerPropagatesError(t *testing.T) {
	s := makeShedder()
	s.Record(0.1)

	sentinel := errors.New("fn error")
	r := NewShedderRunner(s, func(_ context.Context) error {
		return sentinel
	})
	err := r.Run(context.Background())
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestShedderRunnerNilShedderDefaults(t *testing.T) {
	r := NewShedderRunner(nil, func(_ context.Context) error { return nil })
	if err := r.Run(context.Background()); err != nil {
		t.Errorf("expected no error with nil shedder, got %v", err)
	}
}

func TestShedderRunnerNilFnDefaults(t *testing.T) {
	s := makeShedder()
	s.Record(0.1)
	r := NewShedderRunner(s, nil)
	if err := r.Run(context.Background()); err != nil {
		t.Errorf("expected no error with nil fn, got %v", err)
	}
}
