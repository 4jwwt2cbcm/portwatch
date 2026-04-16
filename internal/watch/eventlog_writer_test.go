package watch

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestEventLogWriterFlushesEvents(t *testing.T) {
	log := NewEventLog(10)
	log.Append(EventScanOK, "scan done")
	log.Append(EventPortAdded, "port 443")

	var buf bytes.Buffer
	w := NewEventLogWriter(log, 10*time.Millisecond, &buf)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	go w.Run(ctx)

	<-ctx.Done()
	out := buf.String()
	if !strings.Contains(out, "scan done") {
		t.Errorf("expected 'scan done' in output, got: %s", out)
	}
	if !strings.Contains(out, "port 443") {
		t.Errorf("expected 'port 443' in output, got: %s", out)
	}
}

func TestEventLogWriterEmptyLogWritesNothing(t *testing.T) {
	log := NewEventLog(10)
	var buf bytes.Buffer
	w := NewEventLogWriter(log, 10*time.Millisecond, &buf)

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()
	go w.Run(ctx)

	<-ctx.Done()
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty log, got: %s", buf.String())
	}
}

func TestEventLogWriterDefaultsToStderr(t *testing.T) {
	log := NewEventLog(10)
	w := NewEventLogWriter(log, time.Second, nil)
	if w.out == nil {
		t.Error("expected non-nil writer")
	}
}

func TestEventLogWriterStopsOnCancel(t *testing.T) {
	log := NewEventLog(10)
	var buf bytes.Buffer
	w := NewEventLogWriter(log, 5*time.Millisecond, &buf)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		w.Run(ctx)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Error("writer did not stop after context cancel")
	}
}
