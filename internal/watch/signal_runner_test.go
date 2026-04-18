package watch

import (
	"context"
	"errors"
	"syscall"
	"testing"
	"time"
)

func TestSignalRunnerRunsFunction(t *testing.T) {
	r := NewSignalRunner(nil, syscall.SIGINT)
	called := false
	err := r.Run(context.Background(), func(ctx context.Context) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("fn was not called")
	}
}

func TestSignalRunnerPropagatesError(t *testing.T) {
	r := NewSignalRunner(nil, syscall.SIGINT)
	want := errors.New("boom")
	err := r.Run(context.Background(), func(ctx context.Context) error {
		return want
	})
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestSignalRunnerStopsOnContextCancel(t *testing.T) {
	r := NewSignalRunner(nil, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	err := r.Run(ctx, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	if err != nil {
		t.Fatalf("expected nil on context cancel, got %v", err)
	}
}

func TestSignalRunnerNilLoggerDefaults(t *testing.T) {
	r := NewSignalRunner(nil)
	if r.logger == nil {
		t.Fatal("logger should default to non-nil")
	}
}
