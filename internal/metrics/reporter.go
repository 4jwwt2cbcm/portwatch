package metrics

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// Reporter periodically writes metrics snapshots to an output writer.
type Reporter struct {
	collector *Collector
	interval  time.Duration
	out       io.Writer
}

// NewReporter creates a Reporter that prints metrics from the given collector.
// If out is nil, os.Stdout is used.
func NewReporter(c *Collector, interval time.Duration, out io.Writer) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{collector: c, interval: interval, out: out}
}

// Run blocks, printing a metrics report every interval until ctx is done.
func (r *Reporter) Run(ctx interface{ Done() <-chan struct{} }) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.Report()
		case <-ctx.Done():
			return
		}
	}
}

// Report writes a single snapshot of current metrics to the output writer.
func (r *Reporter) Report() {
	s := r.collector.Snapshot()
	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "--- portwatch metrics ---")
	fmt.Fprintf(w, "total_scans\t%d\n", s.TotalScans)
	fmt.Fprintf(w, "total_ports_seen\t%d\n", s.TotalPortsSeen)
	fmt.Fprintf(w, "ports_added\t%d\n", s.PortsAdded)
	fmt.Fprintf(w, "ports_removed\t%d\n", s.PortsRemoved)
	if !s.LastScanAt.IsZero() {
		fmt.Fprintf(w, "last_scan_at\t%s\n", s.LastScanAt.Format(time.RFC3339))
		fmt.Fprintf(w, "last_scan_duration\t%s\n", s.LastScanDuration)
		fmt.Fprintf(w, "last_port_count\t%d\n", s.LastPortCount)
	}
	_ = w.Flush()
}
