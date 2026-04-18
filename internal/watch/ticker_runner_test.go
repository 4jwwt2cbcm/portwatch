package watch

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestTickerRunnerStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	var calls int64
	r := NewTickerRunner(50*time.Millisecond, func(_ context.Context) error {
		atomic.AddInt64(&calls, 1)
		return nil
	}, nil)

	done := make(chan error, 1)
	go func() { done <- r.Run(ctx) }()

	time.Sleep(160 * time.Millisecond)
	cancel()

	err := <-done
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if atomic.LoadInt64(&calls) < 2 {
		t.Fatalf("expected at least 2 calls, got %d", calls)
	}
}

func TestTickerRunnerStopsOnCallbackError(t *testing.T) {
	sentinel := errors.New("stop")
	r := NewTickerRunner(20*time.Millisecond, func(_ context.Context) error {
		return sentinel
	}, nil)

	ctx := context.Background()
	err := r.Run(ctx)
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestTickerRunnerClampsInterval(t *testing.T) {
	// interval of 0 should be clamped to minimum by ClampInterval
	r := NewTickerRunner(0, func(_ context.Context) error { return nil }, nil)
	if r.interval <= 0 {
		t.Fatalf("expected positive interval after clamp, got %v", r.interval)
	}
}

func TestTickerRunnerNilLoggerDefaults(t *testing.T) {
	r := NewTickerRunner(50*time.Millisecond, func(_ context.Context) error { return nil }, nil)
	if r.logger == nil {
		t.Fatal("expected non-nil logger")
	}
}
