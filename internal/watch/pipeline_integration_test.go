package watch_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/user/portwatch/internal/watch"
)

func TestPipelineIntegration(t *testing.T) {
	var log []string

	record := func(label string) watch.Stage[string] {
		return func(_ context.Context, in string) (string, error) {
			log = append(log, label)
			return fmt.Sprintf("%s|%s", in, label), nil
		}
	}

	p := watch.NewPipeline(
		record("a"),
		record("b"),
		record("c"),
	)

	out, err := p.Run(context.Background(), "start")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "start|a|b|c"
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}

	if len(log) != 3 {
		t.Errorf("expected 3 log entries, got %d", len(log))
	}
}
