package watch

import (
	"sync"
	"testing"
)

func TestHookRegistryRegisterAndFire(t *testing.T) {
	r := NewHookRegistry()
	called := false
	r.Register(HookBeforeScan, func(event HookEvent, meta map[string]any) {
		called = true
	})
	r.Fire(HookBeforeScan, nil)
	if !called {
		t.Fatal("expected hook to be called")
	}
}

func TestHookRegistryMultipleHooks(t *testing.T) {
	r := NewHookRegistry()
	count := 0
	for i := 0; i < 3; i++ {
		r.Register(HookAfterScan, func(event HookEvent, meta map[string]any) {
			count++
		})
	}
	r.Fire(HookAfterScan, nil)
	if count != 3 {
		t.Fatalf("expected 3 calls, got %d", count)
	}
}

func TestHookRegistryMetaPassed(t *testing.T) {
	r := NewHookRegistry()
	var got map[string]any
	r.Register(HookOnError, func(event HookEvent, meta map[string]any) {
		got = meta
	})
	r.Fire(HookOnError, map[string]any{"err": "timeout"})
	if got["err"] != "timeout" {
		t.Fatalf("unexpected meta: %v", got)
	}
}

func TestHookRegistryCount(t *testing.T) {
	r := NewHookRegistry()
	if r.Count(HookBeforeScan) != 0 {
		t.Fatal("expected zero hooks")
	}
	r.Register(HookBeforeScan, func(HookEvent, map[string]any) {})
	if r.Count(HookBeforeScan) != 1 {
		t.Fatal("expected one hook")
	}
}

func TestHookRegistryClear(t *testing.T) {
	r := NewHookRegistry()
	r.Register(HookAfterScan, func(HookEvent, map[string]any) {})
	r.Clear(HookAfterScan)
	if r.Count(HookAfterScan) != 0 {
		t.Fatal("expected hooks to be cleared")
	}
}

func TestHookRegistryConcurrentFire(t *testing.T) {
	r := NewHookRegistry()
	var mu sync.Mutex
	count := 0
	r.Register(HookBeforeScan, func(HookEvent, map[string]any) {
		mu.Lock()
		count++
		mu.Unlock()
	})
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.Fire(HookBeforeScan, nil)
		}()
	}
	wg.Wait()
	if count != 10 {
		t.Fatalf("expected 10, got %d", count)
	}
}
