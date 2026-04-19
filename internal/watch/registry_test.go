package watch

import (
	"errors"
	"testing"
)

func TestRegistryEmptyOnInit(t *testing.T) {
	r := NewRegistry()
	if r.Count() != 0 {
		t.Fatalf("expected 0, got %d", r.Count())
	}
}

func TestRegistryRegisterAndRun(t *testing.T) {
	r := NewRegistry()
	called := false
	err := r.Register("test", func() error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := r.Run("test"); err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	if !called {
		t.Fatal("expected runner to be called")
	}
}

func TestRegistryDuplicateReturnsError(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("dup", func() error { return nil })
	err := r.Register("dup", func() error { return nil })
	if err == nil {
		t.Fatal("expected error for duplicate registration")
	}
}

func TestRegistryRunUnknownReturnsError(t *testing.T) {
	r := NewRegistry()
	err := r.Run("missing")
	if err == nil {
		t.Fatal("expected error for unknown runner")
	}
}

func TestRegistryRunPropagatesError(t *testing.T) {
	r := NewRegistry()
	expected := errors.New("boom")
	_ = r.Register("fail", func() error { return expected })
	if err := r.Run("fail"); !errors.Is(err, expected) {
		t.Fatalf("expected %v, got %v", expected, err)
	}
}

func TestRegistryUnregister(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("gone", func() error { return nil })
	r.Unregister("gone")
	if r.Count() != 0 {
		t.Fatalf("expected 0 after unregister, got %d", r.Count())
	}
	if err := r.Run("gone"); err == nil {
		t.Fatal("expected error after unregister")
	}
}

func TestRegistryNamesReflectsEntries(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("a", func() error { return nil })
	_ = r.Register("b", func() error { return nil })
	names := r.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}
