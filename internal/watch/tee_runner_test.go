package watch

import (
	"context"
	"errors"
	"testing"
)

func TestTeeRunnerEmitsOnSuccess(t *testing.T) {
	var got int
	tee := NewTee(func(v int, _ error) { got = v })
	r := NewTeeRunner(tee, func(_ context.Context) (int, error) { return 7, nil })
	v, err := r.Run(context.Background())
	if err != nil || v != 7 || got != 7 {
		t.Fatalf("expected 7, got v=%d got=%d err=%v", v, got, err)
	}
}

func TestTeeRunnerEmitsOnError(t *testing.T) {
	sentinel := errors.New("oops")
	var gotErr error
	tee := NewTee(func(_ int, err error) { gotErr = err })
	r := NewTeeRunner(tee, func(_ context.Context) (int, error) { return 0, sentinel })
	_, err := r.Run(context.Background())
	if err != sentinel || gotErr != sentinel {
		t.Fatalf("expected sentinel, got err=%v gotErr=%v", err, gotErr)
	}
}

func TestTeeRunnerNilTeeDefaults(t *testing.T) {
	r := NewTeeRunner[int](nil, func(_ context.Context) (int, error) { return 1, nil })
	v, err := r.Run(context.Background())
	if err != nil || v != 1 {
		t.Fatalf("expected 1, got v=%d err=%v", v, err)
	}
}

func TestTeeRunnerNilFnDefaults(t *testing.T) {
	tee := NewTee[int]()
	r := NewTeeRunner(tee, nil)
	v, err := r.Run(context.Background())
	if err != nil || v != 0 {
		t.Fatalf("expected zero, got v=%d err=%v", v, err)
	}
}

func TestTeeRunnerPropagatesContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	sentinel := context.Canceled
	r := NewTeeRunner[int](nil, func(ctx context.Context) (int, error) {
		return 0, ctx.Err()
	})
	_, err := r.Run(ctx)
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
