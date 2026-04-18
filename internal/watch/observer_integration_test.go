package watch_test

import (
	"testing"

	"github.com/user/portwatch/internal/watch"
)

func TestObserverIntegration(t *testing.T) {
	obs := watch.NewObserver()

	var scanEvents []any
	var alertEvents []any

	unsub1 := obs.Subscribe("scan", func(_ string, p any) { scanEvents = append(scanEvents, p) })
	obs.Subscribe("alert", func(_ string, p any) { alertEvents = append(alertEvents, p) })

	obs.Publish("scan", "result-1")
	obs.Publish("scan", "result-2")
	obs.Publish("alert", "alert-1")

	if len(scanEvents) != 2 {
		t.Fatalf("expected 2 scan events, got %d", len(scanEvents))
	}
	if len(alertEvents) != 1 {
		t.Fatalf("expected 1 alert event, got %d", len(alertEvents))
	}

	unsub1()
	obs.Publish("scan", "result-3")
	if len(scanEvents) != 2 {
		t.Fatalf("expected still 2 scan events after unsub, got %d", len(scanEvents))
	}

	obs.Clear("alert")
	obs.Publish("alert", "alert-2")
	if len(alertEvents) != 1 {
		t.Fatalf("expected still 1 alert event after clear, got %d", len(alertEvents))
	}
}
