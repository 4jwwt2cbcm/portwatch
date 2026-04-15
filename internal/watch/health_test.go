package watch

import (
	"errors"
	"testing"
)

func TestNewHealthTrackerDefaultsMaxErrors(t *testing.T) {
	h := NewHealthTracker(0)
	if h.maxErrors != 3 {
		t.Fatalf("expected default maxErrors=3, got %d", h.maxErrors)
	}
}

func TestHealthyOnInit(t *testing.T) {
	h := NewHealthTracker(3)
	if !h.Status().Healthy {
		t.Fatal("expected healthy on init")
	}
}

func TestRecordSuccessResetsErrors(t *testing.T) {
	h := NewHealthTracker(3)
	h.RecordError(errors.New("boom"))
	h.RecordError(errors.New("boom"))
	h.RecordSuccess()

	s := h.Status()
	if s.ConsecErrors != 0 {
		t.Fatalf("expected 0 consec errors after success, got %d", s.ConsecErrors)
	}
	if s.LastError != nil {
		t.Fatalf("expected nil LastError after success, got %v", s.LastError)
	}
	if s.LastSuccess.IsZero() {
		t.Fatal("expected LastSuccess to be set")
	}
}

func TestRecordErrorIncrementsCounter(t *testing.T) {
	h := NewHealthTracker(3)
	err := errors.New("scan failed")
	h.RecordError(err)

	s := h.Status()
	if s.ConsecErrors != 1 {
		t.Fatalf("expected 1 consec error, got %d", s.ConsecErrors)
	}
	if !errors.Is(s.LastError, err) {
		t.Fatalf("expected LastError=%v, got %v", err, s.LastError)
	}
}

func TestUnhealthyAfterMaxErrors(t *testing.T) {
	h := NewHealthTracker(2)
	h.RecordError(errors.New("e1"))
	if !h.Status().Healthy {
		t.Fatal("should still be healthy after 1 error with max=2")
	}
	h.RecordError(errors.New("e2"))
	if h.Status().Healthy {
		t.Fatal("should be unhealthy after reaching maxErrors")
	}
}

func TestRecordSuccessRestoresHealth(t *testing.T) {
	h := NewHealthTracker(2)
	h.RecordError(errors.New("e1"))
	h.RecordError(errors.New("e2"))
	if h.Status().Healthy {
		t.Fatal("should be unhealthy before success")
	}
	h.RecordSuccess()
	if !h.Status().Healthy {
		t.Fatal("should be healthy after success")
	}
}
