package watch

import (
	"testing"
	"time"
)

func makeFuse() *Fuse {
	return NewFuse(FusePolicy{MaxErrors: 3, ResetAfter: 10 * time.Second})
}

func TestDefaultFusePolicyValues(t *testing.T) {
	p := DefaultFusePolicy()
	if p.MaxErrors != 5 {
		t.Fatalf("expected MaxErrors 5, got %d", p.MaxErrors)
	}
	if p.ResetAfter != 30*time.Second {
		t.Fatalf("expected ResetAfter 30s, got %v", p.ResetAfter)
	}
}

func TestFuseNotBlownOnInit(t *testing.T) {
	f := makeFuse()
	if f.Blown() {
		t.Fatal("expected fuse not blown on init")
	}
}

func TestFuseBlowsAfterMaxErrors(t *testing.T) {
	f := makeFuse()
	f.Record()
	f.Record()
	if f.Blown() {
		t.Fatal("should not be blown before max")
	}
	f.Record()
	if !f.Blown() {
		t.Fatal("expected fuse blown after max errors")
	}
}

func TestFuseResetClearsState(t *testing.T) {
	f := makeFuse()
	f.Record()
	f.Record()
	f.Record()
	if !f.Blown() {
		t.Fatal("expected blown")
	}
	f.Reset()
	if f.Blown() {
		t.Fatal("expected not blown after reset")
	}
	if f.Errors() != 0 {
		t.Fatalf("expected 0 errors after reset, got %d", f.Errors())
	}
}

func TestFuseErrorsCount(t *testing.T) {
	f := makeFuse()
	f.Record()
	f.Record()
	if f.Errors() != 2 {
		t.Fatalf("expected 2 errors, got %d", f.Errors())
	}
}

func TestFuseAutoResetAfterWindow(t *testing.T) {
	f := NewFuse(FusePolicy{MaxErrors: 3, ResetAfter: 50 * time.Millisecond})
	f.Record()
	f.Record()
	time.Sleep(60 * time.Millisecond)
	f.Record()
	if f.Blown() {
		t.Fatal("fuse should have auto-reset after window")
	}
	if f.Errors() != 1 {
		t.Fatalf("expected 1 error after auto-reset, got %d", f.Errors())
	}
}
