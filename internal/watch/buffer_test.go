package watch

import (
	"testing"
)

func TestBufferDefaultsCapacityOnZero(t *testing.T) {
	b := NewBuffer[int](0)
	if b.Cap() != 1 {
		t.Fatalf("expected cap 1, got %d", b.Cap())
	}
}

func TestBufferAddReturnsFalseWhenNotFull(t *testing.T) {
	b := NewBuffer[int](3)
	full := b.Add(1)
	if full {
		t.Fatal("expected not full after first add")
	}
}

func TestBufferAddReturnsTrueWhenFull(t *testing.T) {
	b := NewBuffer[int](2)
	b.Add(1)
	full := b.Add(2)
	if !full {
		t.Fatal("expected full after reaching capacity")
	}
}

func TestBufferFlushReturnsItems(t *testing.T) {
	b := NewBuffer[string](4)
	b.Add("a")
	b.Add("b")
	items := b.Flush()
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0] != "a" || items[1] != "b" {
		t.Fatalf("unexpected items: %v", items)
	}
}

func TestBufferFlushResetsLen(t *testing.T) {
	b := NewBuffer[int](4)
	b.Add(1)
	b.Add(2)
	b.Flush()
	if b.Len() != 0 {
		t.Fatalf("expected len 0 after flush, got %d", b.Len())
	}
}

func TestBufferFlushEmptyReturnsEmpty(t *testing.T) {
	b := NewBuffer[int](4)
	items := b.Flush()
	if len(items) != 0 {
		t.Fatalf("expected empty flush, got %d items", len(items))
	}
}

func TestBufferLenTracksAdds(t *testing.T) {
	b := NewBuffer[int](10)
	for i := 0; i < 5; i++ {
		b.Add(i)
	}
	if b.Len() != 5 {
		t.Fatalf("expected len 5, got %d", b.Len())
	}
}
