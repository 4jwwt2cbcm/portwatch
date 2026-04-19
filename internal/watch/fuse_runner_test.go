package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestFuseRunnerSucceedsNormally(t *testing.T) {
	r := NewFuseRunner(nil, func(context.Context) error { return nil })
	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFuseRunnerRecordsError(t *testing.T) {
	sentinel := errors.New("boom")
	r := NewFuseRunner(
		NewFuse(FusePolicy{MaxErrors: 5, ResetAfter: time.Second}),
		func(context.Context) error { return sentinel },
	)
	err := r.Run(context.Background())
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if r.Fuse().Errors() != 1 {
		t.Fatalf("expected 1 recorded error, got %d", r.Fuse().Errors())
	}
}

func TestFuseRunnerBlocksWhenBlown(t *testing.T) {
	f := NewFuse(FusePolicy{MaxErrors: 2, ResetAfter: time.Second})
	calls := 0
	r := NewFuseRunner(f, func(context.Context) error {
		calls++
		return errors.New("err")
	})
	r.Run(context.Background())
	r.Run(context.Background())
	if !f.Blown() {
		t.Fatal("expected fuse blown")
	}
	err := r.Run(context.Background())
	if !errors.Is(err, ErrFuseBlown) {
		t.Fatalf("expected ErrFuseBlown, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected fn called 2 times, got %d", calls)
	}
}

func TestFuseRunnerNilArgDefaults(t *testing.T) {
	r := NewFuseRunner(nil, nil)
	if r.Fuse() == nil {
		t.Fatal("expected default fuse")
	}
	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
