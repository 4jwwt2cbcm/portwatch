package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makePort() scanner.PortState {
	return scanner.PortState{Port: 8080, Protocol: "tcp", State: "open"}
}

func TestNotifyWritesOutput(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)

	a := alert.Alert{
		Timestamp: time.Now(),
		Level:     alert.LevelAlert,
		Message:   "test message",
		Port:      makePort(),
	}

	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", out)
	}
	if !strings.Contains(out, "test message") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestNewPortAlert(t *testing.T) {
	ps := makePort()
	a := alert.NewPortAlert(ps)

	if a.Level != alert.LevelAlert {
		t.Errorf("expected LevelAlert, got %s", a.Level)
	}
	if !strings.Contains(a.Message, "new port open") {
		t.Errorf("unexpected message: %s", a.Message)
	}
	if a.Port != ps {
		t.Errorf("port mismatch")
	}
}

func TestClosedPortAlert(t *testing.T) {
	ps := makePort()
	a := alert.ClosedPortAlert(ps)

	if a.Level != alert.LevelWarn {
		t.Errorf("expected LevelWarn, got %s", a.Level)
	}
	if !strings.Contains(a.Message, "port closed") {
		t.Errorf("unexpected message: %s", a.Message)
	}
}

func TestNewNotifierDefaultsToStdout(t *testing.T) {
	n := alert.NewNotifier(nil)
	if n.Writer == nil {
		t.Error("expected non-nil writer")
	}
}
