package watch_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func TestEventLogWriterIntegration(t *testing.T) {
	log := watch.NewEventLog(50)

	var buf bytes.Buffer
	writer := watch.NewEventLogWriter(log, 20*time.Millisecond, &buf)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go writer.Run(ctx)

	log.Append(watch.EventScanOK, "initial scan")
	log.Append(watch.EventPortAdded, "port 8080 opened")
	log.Append(watch.EventPortGone, "port 22 closed")

	time.Sleep(60 * time.Millisecond)
	cancel()

	out := buf.String()
	for _, want := range []string{"initial scan", "port 8080 opened", "port 22 closed"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output\ngot:\n%s", want, out)
		}
	}

	if !strings.Contains(out, "scan_ok") {
		t.Errorf("expected event kind in output")
	}
}
