package watch

import (
	"errors"
	"testing"
	"time"
)

func TestNewEnvelopeSetsFields(t *testing.T) {
	before := time.Now()
	env := NewEnvelope("id-1", "scanner", 42)
	after := time.Now()

	if env.Meta.ID != "id-1" {
		t.Errorf("expected id-1, got %s", env.Meta.ID)
	}
	if env.Meta.Source != "scanner" {
		t.Errorf("expected scanner, got %s", env.Meta.Source)
	}
	if env.Payload != 42 {
		t.Errorf("expected payload 42, got %d", env.Payload)
	}
	if env.Meta.Attempt != 1 {
		t.Errorf("expected attempt 1, got %d", env.Meta.Attempt)
	}
	if env.Meta.Timestamp.Before(before) || env.Meta.Timestamp.After(after) {
		t.Error("timestamp out of expected range")
	}
}

func TestEnvelopeOKNoError(t *testing.T) {
	env := NewEnvelope("id-2", "src", "hello")
	if !env.OK() {
		t.Error("expected OK to be true with no error")
	}
}

func TestEnvelopeWithError(t *testing.T) {
	env := NewEnvelope("id-3", "src", "hello")
	errored := env.WithError(errors.New("boom"))

	if errored.OK() {
		t.Error("expected OK to be false after WithError")
	}
	if errored.Err.Error() != "boom" {
		t.Errorf("unexpected error: %v", errored.Err)
	}
	// original unchanged
	if env.Err != nil {
		t.Error("original envelope should not be mutated")
	}
}

func TestEnvelopeWithAttempt(t *testing.T) {
	env := NewEnvelope("id-4", "src", 0)
	updated := env.WithAttempt(3)

	if updated.Meta.Attempt != 3 {
		t.Errorf("expected attempt 3, got %d", updated.Meta.Attempt)
	}
	if env.Meta.Attempt != 1 {
		t.Error("original attempt should remain 1")
	}
}

func TestEnvelopeWithErrorPreservesPayload(t *testing.T) {
	env := NewEnvelope("id-5", "src", 99)
	errored := env.WithError(errors.New("fail"))
	if errored.Payload != 99 {
		t.Errorf("expected payload 99, got %d", errored.Payload)
	}
}
