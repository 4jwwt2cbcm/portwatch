package scanner

// Diff holds the result of comparing two port snapshots.
type Diff struct {
	Added   []PortState
	Removed []PortState
}

// HasChanges reports whether the diff contains any changes.
func (d Diff) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0
}

// Compare returns the diff between a previous and current list of PortStates.
func Compare(prev, curr []PortState) Diff {
	prevSet := toSet(prev)
	currSet := toSet(curr)

	var diff Diff

	for key, ps := range currSet {
		if _, exists := prevSet[key]; !exists {
			diff.Added = append(diff.Added, ps)
		}
	}

	for key, ps := range prevSet {
		if _, exists := currSet[key]; !exists {
			diff.Removed = append(diff.Removed, ps)
		}
	}

	return diff
}

// toSet converts a slice of PortState into a map keyed by "port/protocol".
func toSet(ports []PortState) map[string]PortState {
	m := make(map[string]PortState, len(ports))
	for _, p := range ports {
		key := p.String()
		m[key] = p
	}
	return m
}
