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
