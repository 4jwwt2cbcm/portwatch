package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func TestDispatchAdded(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)
	d := alert.NewDispatcher(n)

	diff := scanner.Diff{
		Added: []scanner.PortState{
			{Port: 443, Protocol: "tcp", State: "open"},
		},
		Removed: nil,
	}

	if err := d.Dispatch(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", out)
	}
	if !strings.Contains(out, "443") {
		t.Errorf("expected port 443 in output, got: %s", out)
	}
}

func TestDispatchRemoved(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)
	d := alert.NewDispatcher(n)

	diff := scanner.Diff{
		Added: nil,
		Removed: []scanner.PortState{
			{Port: 22, Protocol: "tcp", State: "open"},
		},
	}

	if err := d.Dispatch(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output, got: %s", out)
	}
	if !strings.Contains(out, "22") {
		t.Errorf("expected port 22 in output, got: %s", out)
	}
}

func TestDispatchNoDiff(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)
	d := alert.NewDispatcher(n)

	if err := d.Dispatch(scanner.Diff{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}
