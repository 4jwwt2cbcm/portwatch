package watch

import (
	"context"
	"errors"
	"testing"
)

func TestGateRunnerAllowsWhenOpen(t *testing.T) {
	g := NewGate(true)
	called := false
	r := NewGateRunner(g, func(_ context.Context) error {
		called = true
		return nil
	})
	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected fn to be called when gate is open")
	}
}

func TestGateRunnerSkipsWhenClosed(t *testing.T) {
	g := NewGate(false)
	called := false
	r := NewGateRunner(g, func(_ context.Context) error {
		called = true
		return nil
	})
	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected fn to be skipped when gate is closed")
	}
}

func TestGateRunnerPropagatesError(t *testing.T) {
	g := NewGate(true)
	sentinel := errors.New("boom")
	r := NewGateRunner(g, func(_ context.Context) error {
		return sentinel
	})
	if err := r.Run(context.Background()); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestGateRunnerNilGateDefaultsOpen(t *testing.T) {
	called := false
	r := NewGateRunner(nil, func(_ context.Context) error {
		called = true
		return nil
	})
	_ = r.Run(context.Background())
	if !called {
		t.Fatal("expected fn to be called when nil gate defaults to open")
	}
}

func TestGateRunnerRespectsDynamicToggle(t *testing.T) {
	g := NewGate(true)
	count := 0
	r := NewGateRunner(g, func(_ context.Context) error {
		count++
		return nil
	})
	_ = r.Run(context.Background())
	g.Close()
	_ = r.Run(context.Background())
	g.Open()
	_ = r.Run(context.Background())
	if count != 2 {
		t.Fatalf("expected fn called 2 times, got %d", count)
	}
}
