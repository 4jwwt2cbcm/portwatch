package watch

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Logger defines the interface for watch-cycle logging.
type Logger interface {
	Info(msg string)
	Error(msg string, err error)
}

// defaultLogger writes timestamped log lines to a writer.
type defaultLogger struct {
	out io.Writer
}

// NewLogger returns a Logger that writes to w.
// If w is nil, os.Stderr is used.
func NewLogger(w io.Writer) Logger {
	if w == nil {
		w = os.Stderr
	}
	return &defaultLogger{out: w}
}

func (l *defaultLogger) Info(msg string) {
	fmt.Fprintf(l.out, "%s [INFO]  %s\n", timestamp(), msg)
}

func (l *defaultLogger) Error(msg string, err error) {
	if err != nil {
		fmt.Fprintf(l.out, "%s [ERROR] %s: %v\n", timestamp(), msg, err)
		return
	}
	fmt.Fprintf(l.out, "%s [ERROR] %s\n", timestamp(), msg)
}

func timestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}
