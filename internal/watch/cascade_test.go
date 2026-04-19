package watch

import (
	"context"
	"errors"
	"testing"
)

func makeCascade() *Cascade {
	return NewCascade()
}

func TestCascadeEmptyRunSucceeds(t *testing.T) {
	c := makeCascade()
	if err := c.Run(context.Background()); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCascadeRunsStepsInOrder(t *testing.T) {
	c := makeCascade()
	var order []int
	c.Add("a", func(_ context.Context) error { order = append(order, 1); return nil })
	c.Add("b", func(_ context.Context) error { order = append(order, 2); return nil })
	c.Add("c", func(_ context.Context) error { order = append(order, 3); return nil })
	if err := c.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 3 || order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Fatalf("unexpected order: %v", order)
	}
}

func TestCascadeAbortsOnStepError(t *testing.T) {
	c := makeCascade()
	sentinel := errors.New("step failed")
	ran := false
	c.Add("fail", func(_ context.Context) error { return sentinel })
	c.Add("after", func(_ context.Context) error { ran = true; return nil })
	err := c.Run(context.Background())
	if !errors.Is(err, ErrCascadeAborted) {
		t.Fatalf("expected ErrCascadeAborted, got %v", err)
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel wrapped, got %v", err)
	}
	if ran {
		t.Fatal("step after failure should not have run")
	}
}

func TestCascadeCancelledContextStopsEarly(t *testing.T) {
	c := makeCascade()
	ran := false
	c.Add("step", func(_ context.Context) error { ran = true; return nil })
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := c.Run(ctx)
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
	if ran {
		t.Fatal("step should not run after cancel")
	}
}

func TestCascadeLenAndClear(t *testing.T) {
	c := makeCascade()
	c.Add("x", func(_ context.Context) error { return nil })
	c.Add("y", func(_ context.Context) error { return nil })
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
	c.Clear()
	if c.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", c.Len())
	}
}
