package watch

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestRunLogIntegration(t *testing.T) {
	rl := NewRunLog(5)

	// Simulate several cycle outcomes.
	for i := 0; i < 3; i++ {
		rl.Append(RunEntry{
			StartedAt:  time.Now(),
			FinishedAt: time.Now().Add(5 * time.Millisecond),
			PortsFound: (i + 1) * 2,
		})
	}

	if rl.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", rl.Len())
	}

	var buf bytes.Buffer
	reporter := NewRunLogReporter(rl, &buf, time.Minute)
	reporter.Report()

	out := buf.String()
	if !strings.Contains(out, "total=3") {
		t.Errorf("report missing total=3: %q", out)
	}
	if !strings.Contains(out, "failed=0") {
		t.Errorf("report missing failed=0: %q", out)
	}

	rl.Clear()
	if rl.Len() != 0 {
		t.Fatal("expected empty log after clear")
	}
}

// TestRunLogIntegration_WithFailures verifies that failed entries are counted
// correctly in the reporter output when some run entries have errors.
func TestRunLogIntegration_WithFailures(t *testing.T) {
	rl := NewRunLog(10)

	// Add two successful entries and one failed entry.
	rl.Append(RunEntry{StartedAt: time.Now(), FinishedAt: time.Now().Add(2 * time.Millisecond), PortsFound: 4})
	rl.Append(RunEntry{StartedAt: time.Now(), FinishedAt: time.Now().Add(2 * time.Millisecond), PortsFound: 0, Err: fmt.Errorf("scan failed")})
	rl.Append(RunEntry{StartedAt: time.Now(), FinishedAt: time.Now().Add(2 * time.Millisecond), PortsFound: 6})

	if rl.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", rl.Len())
	}

	var buf bytes.Buffer
	reporter := NewRunLogReporter(rl, &buf, time.Minute)
	reporter.Report()

	out := buf.String()
	if !strings.Contains(out, "total=3") {
		t.Errorf("report missing total=3: %q", out)
	}
	if !strings.Contains(out, "failed=1") {
		t.Errorf("report missing failed=1: %q", out)
	}
}
