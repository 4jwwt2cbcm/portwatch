package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestEvictRunnerReturnsCachedResult(t *testing.T) {
	calls := 0
	fn := func(_ context.Context, _ string) (int, error) {
		calls++
		return 42, nil
	}
	r := NewEvictRunner(NewEvictCache[int](EvictPolicy{TTL: time.Minute, Capacity: 10}), fn)
	v1, _ := r.Run(context.Background(), "k")
	v2, _ := r.Run(context.Background(), "k")
	if v1 != 42 || v2 != 42 {
		t.Fatalf("expected 42, got %d %d", v1, v2)
	}
	if calls != 1 {
		t.Fatalf("expected fn called once, got %d", calls)
	}
}

func TestEvictRunnerCallsFnAfterExpiry(t *testing.T) {
	calls := 0
	fn := func(_ context.Context, _ string) (int, error) {
		calls++
		return calls, nil
	}
	r := NewEvictRunner(NewEvictCache[int](EvictPolicy{TTL: time.Millisecond, Capacity: 10}), fn)
	r.Run(context.Background(), "k")
	time.Sleep(5 * time.Millisecond)
	r.Run(context.Background(), "k")
	if calls != 2 {
		t.Fatalf("expected fn called twice after expiry, got %d", calls)
	}
}

func TestEvictRunnerPropagatesError(t *testing.T) {
	sentinel := errors.New("boom")
	fn := func(_ context.Context, _ string) (string, error) {
		return "", sentinel
	}
	r := NewEvictRunner[string](nil, fn)
	_, err := r.Run(context.Background(), "k")
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestEvictRunnerDoesNotCacheOnError(t *testing.T) {
	calls := 0
	fn := func(_ context.Context, _ string) (int, error) {
		calls++
		return 0, errors.New("err")
	}
	r := NewEvictRunner(NewEvictCache[int](EvictPolicy{TTL: time.Minute, Capacity: 10}), fn)
	r.Run(context.Background(), "k")
	r.Run(context.Background(), "k")
	if calls != 2 {
		t.Fatalf("expected fn called twice (no cache on error), got %d", calls)
	}
}

func TestEvictRunnerNilArgDefaults(t *testing.T) {
	r := NewEvictRunner[string](nil, nil)
	v, err := r.Run(context.Background(), "x")
	if err != nil || v != "" {
		t.Fatalf("expected zero string and nil error, got %q %v", v, err)
	}
}
