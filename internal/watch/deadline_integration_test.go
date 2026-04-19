package watch_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func TestDeadlineIntegration(t *testing.T) {
	// Deadline fires and DeadlineRunner stops iterating via RunUntil.
	expiry := time.Now().Add(60 * time.Millisecond)
	d := watch.NewDeadline(watch.DeadlinePolicy{At: expiry})

	var count int
	r := watch.NewDeadlineRunner(d, func(ctx context.Context) error {
		count++
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := r.RunUntil(ctx, 15*time.Millisecond)
	if !errors.Is(err, watch.ErrDeadlineExceeded) {
		t.Fatalf("expected ErrDeadlineExceeded, got %v", err)
	}
	if count == 0 {
		t.Fatal("expected fn to have been called at least once")
	}
	if !d.Expired() {
		t.Fatal("expected deadline to be expired")
	}
}
