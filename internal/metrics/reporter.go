package metrics

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// Reporter prints metric snapshots to a writer.
type Reporter struct {
	out       io.Writer
	collector *Collector
}

// NewReporter returns a Reporter that writes to w.
// If w is nil, os.Stdout is used.
func NewReporter(w io.Writer, c *Collector) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{out: w, collector: c}
}

// Print writes a formatted snapshot to the reporter's writer.
func (r *Reporter) Print() {
	s := r.collector.Snapshot()
	tw := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "--- portwatch metrics ---")
	fmt.Fprintf(tw, "Total scans:\t%d\n", s.TotalScans)
	fmt.Fprintf(tw, "Ports added:\t%d\n", s.PortsAdded)
	fmt.Fprintf(tw, "Ports removed:\t%d\n", s.PortsRemoved)
	if s.LastScanAt.IsZero() {
		fmt.Fprintf(tw, "Last scan:\t-\n")
		fmt.Fprintf(tw, "Last duration:\t-\n")
	} else {
		fmt.Fprintf(tw, "Last scan:\t%s\n", s.LastScanAt.Format(time.RFC3339))
		fmt.Fprintf(tw, "Last duration:\t%s\n", s.LastScanDur)
	}
	_ = tw.Flush()
}
