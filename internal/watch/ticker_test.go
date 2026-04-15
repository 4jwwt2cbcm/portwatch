package watch

import (
	"testing"
	"time"
)

func TestIntervalFromSeconds(t *testing.T) {
	cases := []struct {
		input    int
		expected time.Duration
	}{
		{10, 10 * time.Second},
		{1, time.Second},
		{0, time.Second},
		{-5, time.Second},
		{60, 60 * time.Second},
	}

	for _, tc := range cases {
		got := IntervalFromSeconds(tc.input)
		if got != tc.expected {
			t.Errorf("IntervalFromSeconds(%d) = %s; want %s", tc.input, got, tc.expected)
		}
	}
}

func TestDefaultTickerFactory(t *testing.T) {
	ticker := DefaultTickerFactory(50 * time.Millisecond)
	if ticker == nil {
		t.Fatal("expected non-nil ticker")
	}
	ticker.Stop()
}
