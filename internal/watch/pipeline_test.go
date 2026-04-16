package watch

import (
	"context"
	"errors"
	"testing"
)

func addOne(_ context.Context, in int) (int, error) {
	return in + 1, nil
}

func double(_ context.Context, in int) (int, error) {
	return in * 2, nil
}

func failStage(_ context.Context, in int) (int, error) {
	return in, errors.New("stage failed")
}

func TestPipelineEmptyReturnsInput(t *testing.T) {
	p := NewPipeline[int]()
	out, err := p.Run(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != 5 {
		t.Errorf("expected 5, got %d", out)
	}
}

func TestPipelineRunsStagesInOrder(t *testing.T) {
	// (3 + 1) * 2 = 8
	p := NewPipeline(addOne, double)
	out, err := p.Run(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != 8 {
		t.Errorf("expected 8, got %d", out)
	}
}

func TestPipelineStopsOnError(t *testing.T) {
	called := false
	after := func(_ context.Context, in int) (int, error) {
		called = true
		return in, nil
	}
	p := NewPipeline(addOne, failStage, after)
	_, err := p.Run(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
	if called {
		t.Error("stage after failure should not have been called")
	}
}

func TestPipelineRespectsContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p := NewPipeline(addOne)
	_, err := p.Run(ctx, 0)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestPipelineLen(t *testing.T) {
	p := NewPipeline(addOne, double)
	if p.Len() != 2 {
		t.Errorf("expected 2, got %d", p.Len())
	}
}
