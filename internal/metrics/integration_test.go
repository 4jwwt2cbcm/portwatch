package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/example/portwatch/internal/metrics"
)

// TestReporterIntegration verifies that Collector and Reporter work together
// end-to-end: record several scans, then confirm the report reflects them.
func TestReporterIntegration(t *testing.T) {
	c := metrics.NewCollector()

	c.RecordScan(10, 10, 0, 5*time.Millisecond)
	c.RecordScan(9, 0, 1, 6*time.Millisecond)
	c.RecordScan(11, 2, 0, 4*time.Millisecond)

	var buf bytes.Buffer
	r := metrics.NewReporter(c, time.Hour, &buf)
	r.Report()

	out := buf.String()

	mustContain := func(sub string) {
		t.Helper()
		if !strings.Contains(out, sub) {
			t.Errorf("report missing %q; full output:\n%s", sub, out)
		}
	}

	mustContain("total_scans")
	mustContain("3") // 3 scan cycles
	mustContain("ports_added")
	mustContain("ports_removed")
	mustContain("last_scan_at")
}
