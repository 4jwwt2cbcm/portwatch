package watch

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestNewLoggerDefaultsToStderr(t *testing.T) {
	l := NewLogger(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
	dl, ok := l.(*defaultLogger)
	if !ok {
		t.Fatal("expected *defaultLogger")
	}
	if dl.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestInfoWritesMessage(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	l.Info("port scan complete")

	out := buf.String()
	if !strings.Contains(out, "[INFO]") {
		t.Errorf("expected [INFO] in output, got: %s", out)
	}
	if !strings.Contains(out, "port scan complete") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestErrorWritesMessageAndErr(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	l.Error("scan failed", errors.New("connection refused"))

	out := buf.String()
	if !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected [ERROR] in output, got: %s", out)
	}
	if !strings.Contains(out, "scan failed") {
		t.Errorf("expected message in output, got: %s", out)
	}
	if !strings.Contains(out, "connection refused") {
		t.Errorf("expected error text in output, got: %s", out)
	}
}

func TestErrorWithNilErrOmitsColon(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	l.Error("unexpected state", nil)

	out := buf.String()
	if !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected [ERROR] in output, got: %s", out)
	}
	if strings.Contains(out, "<nil>") {
		t.Errorf("did not expect <nil> in output, got: %s", out)
	}
}

func TestInfoOutputContainsTimestamp(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	l.Info("tick")

	out := buf.String()
	// timestamp format starts with year
	if !strings.Contains(out, "20") {
		t.Errorf("expected timestamp in output, got: %s", out)
	}
}
