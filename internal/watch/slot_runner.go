package watch

import "errors"

// ErrSlotFull is returned when the event count in the current slot has
// reached the configured limit.
var ErrSlotFull = errors.New("slot: rate limit reached for current time slot")

// SlotRunner wraps a function and enforces a per-slot event limit using a
// Slot and a configurable maximum count.
type SlotRunner struct {
	slot  *Slot
	max   int
	clock func() interface{ } // unused; kept for future injectable clock
}

// NewSlotRunner creates a SlotRunner. If slot is nil a default Slot is
// created. If max is zero it defaults to 1.
func NewSlotRunner(slot *Slot, max int) *SlotRunner {
	if slot == nil {
		slot = NewSlot(DefaultSlotPolicy())
	}
	if max <= 0 {
		max = 1
	}
	return &SlotRunner{slot: slot, max: max}
}

// Run executes fn if the current slot has not yet reached the maximum event
// count. On success the slot counter is incremented. Returns ErrSlotFull
// when the limit is exceeded without calling fn.
func (r *SlotRunner) Run(fn func() error) error {
	now := now()
	if r.slot.Count(now) >= r.max {
		return ErrSlotFull
	}
	r.slot.Record(now)
	if fn == nil {
		return nil
	}
	return fn()
}
