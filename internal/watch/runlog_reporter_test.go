package watch

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestRunLogReporterDefaultsToStderr(t *testing.T) {
	rl := NewRunLog(10)
	r := NewRunLogReporter(rl, nil, 0)
	if r.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestRunLogReporterDefaultInterval(t *testing.T) {
	rl := NewRunLog(10)
	r := NewRunLogReporter(rl, nil, 0)
	if r.interval != 60*time.Second {
		t.Fatalf("expected 60s, got %v", r.interval)
	}
}

func TestRunLogReporterReportOutput(t *testing.T) {
	rl := NewRunLog(10)
	rl.Append(makeEntry(4, nil))
	rl.Append(makeEntry(6, errors.New("err")))

	var buf bytes.Buffer
	r := NewRunLogReporter(rl, &buf, time.Minute)
	r.Report()

	out := buf.String()
	if !strings.Contains(out, "total=2") {
		t.Errorf("expected total=2 in output: %q", out)
	}
	if !strings.Contains(out, "failed=1") {
		t.Errorf("expected failed=1 in output: %q", out)
	}
	if !strings.Contains(out, "avg_ports=5") {
		t.Errorf("expected avg_ports=5 in output: %q", out)
	}
}

func TestRunLogReporterEmptyLog(t *testing.T) {
	rl := NewRunLog(10)
	var buf bytes.Buffer
	r := NewRunLogReporter(rl, &buf, time.Minute)
	r.Report()

	if !strings.Contains(buf.String(), "total=0") {
		t.Errorf("expected total=0: %q", buf.String())
	}
}

func TestRunLogReporterStopsOnCancel(t *testing.T) {
	rl := NewRunLog(10)
	var buf bytes.Buffer
	r := NewRunLogReporter(rl, &buf, 10*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Millisecond)
	defer cancel()
	r.Run(ctx)
}
