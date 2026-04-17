package watch

import (
	"errors"
	"testing"
)

func TestHookedRunnerFiresBeforeAndAfter(t *testing.T) {
	reg := NewHookRegistry()
	var fired []HookEvent
	for _, e := range []HookEvent{HookBeforeScan, HookAfterScan, HookOnError} {
		e := e
		reg.Register(e, func(ev HookEvent, _ map[string]any) { fired = append(fired, ev) })
	}

	r := NewHookedRunner(func() error { return nil }, reg)
	if err := r.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fired) != 2 || fired[0] != HookBeforeScan || fired[1] != HookAfterScan {
		t.Fatalf("unexpected hook sequence: %v", fired)
	}
}

func TestHookedRunnerFiresOnError(t *testing.T) {
	reg := NewHookRegistry()
	var fired []HookEvent
	for _, e := range []HookEvent{HookBeforeScan, HookAfterScan, HookOnError} {
		e := e
		reg.Register(e, func(ev HookEvent, _ map[string]any) { fired = append(fired, ev) })
	}

	expected := errors.New("scan failed")
	r := NewHookedRunner(func() error { return expected }, reg)
	if err := r.Run(); err != expected {
		t.Fatalf("expected scan error, got %v", err)
	}
	if len(fired) != 2 || fired[0] != HookBeforeScan || fired[1] != HookOnError {
		t.Fatalf("unexpected hook sequence: %v", fired)
	}
}

func TestHookedRunnerNilRegistryDefaults(t *testing.T) {
	r := NewHookedRunner(func() error { return nil }, nil)
	if err := r.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHookedRunnerErrorMetaContainsMessage(t *testing.T) {
	reg := NewHookRegistry()
	var gotMeta map[string]any
	reg.Register(HookOnError, func(_ HookEvent, meta map[string]any) { gotMeta = meta })

	r := NewHookedRunner(func() error { return errors.New("timeout") }, reg)
	_ = r.Run()
	if gotMeta["err"] != "timeout" {
		t.Fatalf("expected err=timeout in meta, got %v", gotMeta)
	}
}
