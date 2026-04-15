package alert

import (
	"github.com/user/portwatch/internal/scanner"
)

// Dispatcher watches for diffs and dispatches alerts via a Notifier.
type Dispatcher struct {
	notifier *Notifier
}

// NewDispatcher creates a Dispatcher backed by the given Notifier.
func NewDispatcher(n *Notifier) *Dispatcher {
	return &Dispatcher{notifier: n}
}

// Dispatch takes a Diff and emits alerts for each change.
func (d *Dispatcher) Dispatch(diff scanner.Diff) error {
	for _, ps := range diff.Added {
		a := NewPortAlert(ps)
		if err := d.notifier.Notify(a); err != nil {
			return err
		}
	}
	for _, ps := range diff.Removed {
		a := ClosedPortAlert(ps)
		if err := d.notifier.Notify(a); err != nil {
			return err
		}
	}
	return nil
}
