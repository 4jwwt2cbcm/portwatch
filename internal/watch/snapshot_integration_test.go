package watch_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watch"
)

func TestSnapshotStoreConcurrentAccess(t *testing.T) {
	s := watch.NewSnapshotStore()
	ports := []scanner.PortState{
		{Port: 80, Protocol: "tcp", Open: true},
		{Port: 443, Protocol: "tcp", Open: true},
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			s.Set(ports)
		}()
		go func() {
			defer wg.Done()
			snap, ok := s.Get()
			if ok && len(snap.Ports) != 2 {
				t.Errorf("unexpected port count: %d", len(snap.Ports))
			}
		}()
	}
	wg.Wait()

	snap, ok := s.Get()
	if !ok {
		t.Fatal("expected snapshot after concurrent writes")
	}
	if len(snap.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(snap.Ports))
	}
}
