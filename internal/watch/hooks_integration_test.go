package watch_test

import (
	"testing"

	"github.com/user/portwatch/internal/watch"
)

func TestHookRegistryIntegration(t *testing.T) {
	reg := watch.NewHookRegistry()

	events := []watch.HookEvent{}
	record := func(e watch.HookEvent, _ map[string]any) {
		events = append(events, e)
	}

	reg.Register(watch.HookBeforeScan, record)
	reg.Register(watch.HookAfterScan, record)
	reg.Register(watch.HookOnError, record)

	reg.Fire(watch.HookBeforeScan, nil)
	reg.Fire(watch.HookAfterScan, map[string]any{"ports": 5})
	reg.Fire(watch.HookOnError, map[string]any{"err": "connect refused"})

	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	if events[0] != watch.HookBeforeScan {
		t.Errorf("expected HookBeforeScan, got %s", events[0])
	}
	if events[1] != watch.HookAfterScan {
		t.Errorf("expected HookAfterScan, got %s", events[1])
	}
	if events[2] != watch.HookOnError {
		t.Errorf("expected HookOnError, got %s", events[2])
	}

	reg.Clear(watch.HookBeforeScan)
	if reg.Count(watch.HookBeforeScan) != 0 {
		t.Error("expected HookBeforeScan to be cleared")
	}
}
