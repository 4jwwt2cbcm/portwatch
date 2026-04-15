package metrics

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestReportWritesAllFields(t *testing.T) {
	c := NewCollector()
	c.RecordScan(3, 1, 0, 12*time.Millisecond)
	c.RecordScan(4, 0, 1, 8*time.Millisecond)

	var buf bytes.Buffer
	r := NewReporter(c, time.Minute, &buf)
	r.Report()

	out := buf.String()
	for _, want := range []string{
		"total_scans",
		"total_ports_seen",
		"ports_added",
		"ports_removed",
		"last_scan_at",
		"last_scan_duration",
		"last_port_count",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestReportNoScansOmitsTimeFields(t *testing.T) {
	c := NewCollector()
	var buf bytes.Buffer
	r := NewReporter(c, time.Minute, &buf)
	r.Report()

	out := buf.String()
	if strings.Contains(out, "last_scan_at") {
		t.Errorf("expected no last_scan_at when no scans recorded, got:\n%s", out)
	}
}

func TestNewReporterDefaultsToStdout(t *testing.T) {
	c := NewCollector()
	r := NewReporter(c, time.Minute, nil)
	if r.out == nil {
		t.Error("expected non-nil writer when nil passed to NewReporter")
	}
}

func TestReporterRunStopsOnContextCancel(t *testing.T) {
	c := NewCollector()
	var buf bytes.Buffer
	r := NewReporter(c, 10*time.Millisecond, &buf)

	ctx := &cancelCtx{done: make(chan struct{})}
	done := make(chan struct{})
	go func() {
		r.Run(ctx)
		close(done)
	}()

	time.Sleep(35 * time.Millisecond)
	close(ctx.done)

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("reporter did not stop after context cancel")
	}
}

type cancelCtx struct{ done chan struct{} }

func (c *cancelCtx) Done() <-chan struct{} { return c.done }
