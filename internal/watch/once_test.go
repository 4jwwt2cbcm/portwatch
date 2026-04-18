package watch

import (
	"errors"
	"sync"
	"testing"
)

func TestOnceRunnerExecutesOnce(t *testing.T) {
	calls := 0
	r := NewOnceRunner(func() error {
		calls++
		return nil
	})

	for i := 0; i < 5; i++ {
		if err := r.Run(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if calls != 1 {
		t.Errorf("expected fn called once, got %d", calls)
	}
}

func TestOnceRunnerReturnsSameError(t *testing.T) {
	sentinel := errors.New("boom")
	r := NewOnceRunner(func() error { return sentinel })

	for i := 0; i < 3; i++ {
		if err := r.Run(); !errors.Is(err, sentinel) {
			t.Errorf("expected sentinel error, got %v", err)
		}
	}
}

func TestOnceRunnerHasRun(t *testing.T) {
	r := NewOnceRunner(func() error { return nil })
	if r.HasRun() {
		t.Fatal("expected HasRun false before first call")
	}
	_ = r.Run()
	if !r.HasRun() {
		t.Fatal("expected HasRun true after first call")
	}
}

func TestOnceRunnerReset(t *testing.T) {
	calls := 0
	r := NewOnceRunner(func() error {
		calls++
		return nil
	})
	_ = r.Run()
	r.Reset()
	_ = r.Run()
	if calls != 2 {
		t.Errorf("expected 2 calls after reset, got %d", calls)
	}
}

func TestOnceRunnerNilFnDefaults(t *testing.T) {
	r := NewOnceRunner(nil)
	if err := r.Run(); err != nil {
		t.Errorf("expected nil error for nil fn, got %v", err)
	}
}

func TestOnceRunnerConcurrentSafe(t *testing.T) {
	calls := 0
	var mu sync.Mutex
	r := NewOnceRunner(func() error {
		mu.Lock()
		calls++
		mu.Unlock()
		return nil
	})

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = r.Run()
		}()
	}
	wg.Wait()

	if calls != 1 {
		t.Errorf("expected exactly 1 call under concurrency, got %d", calls)
	}
}
