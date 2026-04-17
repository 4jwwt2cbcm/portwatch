package watch

import (
	"sync"
	"testing"
	"time"
)

func TestPausableNotPausedByDefault(t *testing.T) {
	p := NewPausable()
	if p.IsPaused() {
		t.Fatal("expected not paused on init")
	}
}

func TestPausableWaitReturnsImmediatelyWhenNotPaused(t *testing.T) {
	p := NewPausable()
	done := make(chan struct{})
	go func() {
		p.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Wait should return immediately when not paused")
	}
}

func TestPausableWaitBlocksWhilePaused(t *testing.T) {
	p := NewPausable()
	p.Pause()

	if !p.IsPaused() {
		t.Fatal("expected paused")
	}

	done := make(chan struct{})
	go func() {
		p.Wait()
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("Wait should block while paused")
	case <-time.After(50 * time.Millisecond):
	}

	p.Resume()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Wait should unblock after Resume")
	}
}

func TestPausableResumeUnblocksMultipleWaiters(t *testing.T) {
	p := NewPausable()
	p.Pause()

	const n = 5
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			p.Wait()
			wg.Done()
		}()
	}

	time.Sleep(20 * time.Millisecond)
	p.Resume()

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("all waiters should unblock after Resume")
	}
}

func TestPausableIsPausedAfterPause(t *testing.T) {
	p := NewPausable()
	p.Pause()
	if !p.IsPaused() {
		t.Fatal("expected paused")
	}
	p.Resume()
	if p.IsPaused() {
		t.Fatal("expected not paused after resume")
	}
}
