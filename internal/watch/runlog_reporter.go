package watch

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// RunLogReporter periodically writes a summary of recent runs to a writer.
type RunLogReporter struct {
	log      *RunLog
	out      io.Writer
	interval time.Duration
}

// NewRunLogReporter creates a reporter. Zero interval defaults to 60s.
func NewRunLogReporter(rl *RunLog, out io.Writer, interval time.Duration) *RunLogReporter {
	if out == nil {
		out = os.Stderr
	}
	if interval <= 0 {
		interval = 60 * time.Second
	}
	return &RunLogReporter{log: rl, out: out, interval: interval}
}

// Report writes a single summary to the writer.
func (r *RunLogReporter) Report() {
	entries := r.log.Snapshot()
	total := len(entries)
	var failed, totalPorts int
	for _, e := range entries {
		if e.Err != nil {
			failed++
		}
		totalPorts += e.PortsFound
	}
	avgPorts := 0
	if total > 0 {
		avgPorts = totalPorts / total
	}
	fmt.Fprintf(r.out, "[runlog] total=%d failed=%d avg_ports=%d\n", total, failed, avgPorts)
}

// Run reports on the given interval until ctx is cancelled.
func (r *RunLogReporter) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.Report()
		}
	}
}
