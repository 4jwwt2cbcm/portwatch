package scanner

// Diff computes the difference between two port snapshots.
// It returns newly opened ports and ports that have been closed.
type Diff struct {
	Opened []PortState
	Closed []PortState
}

// HasChanges returns true when there are any opened or closed ports.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// Compare takes a previous and current snapshot and returns a Diff.
func Compare(previous, current []PortState) Diff {
	prevSet := toSet(previous)
	currSet := toSet(current)

	var opened, closed []PortState

	for key, ps := range currSet {
		if _, exists := prevSet[key]; !exists {
			opened = append(opened, ps)
		}
	}

	for key, ps := range prevSet {
		if _, exists := currSet[key]; !exists {
			closed = append(closed, ps)
		}
	}

	return Diff{Opened: opened, Closed: closed}
}

func toSet(states []PortState) map[string]PortState {
	m := make(map[string]PortState, len(states))
	for _, s := range states {
		key := s.Protocol + "|" + s.Address + "|" + string(rune(s.Port))
		m[key] = s
	}
	return m
}
