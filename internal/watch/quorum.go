package watch

import (
	"sync"
	"time"
)

// DefaultQuorumPolicy returns a sensible default policy for Quorum.
func DefaultQuorumPolicy() QuorumPolicy {
	return QuorumPolicy{
		MinVoters: 3,
		Threshold: 0.51,
		Window:    30 * time.Second,
	}
}

// QuorumPolicy configures the Quorum behaviour.
type QuorumPolicy struct {
	MinVoters int           // minimum number of voters required
	Threshold float64       // fraction of yes-votes required (0.0–1.0)
	Window    time.Duration // how long a vote remains valid
}

// vote is an internal record of a single cast vote.
type vote struct {
	value bool
	at    time.Time
}

// Quorum tracks votes from named participants and reports whether a
// configurable threshold of agreement has been reached within a time window.
type Quorum struct {
	mu     sync.Mutex
	policy QuorumPolicy
	votes  map[string]vote
}

// NewQuorum creates a Quorum with the given policy.
// Zero-value fields are replaced with defaults.
func NewQuorum(p QuorumPolicy) *Quorum {
	def := DefaultQuorumPolicy()
	if p.MinVoters <= 0 {
		p.MinVoters = def.MinVoters
	}
	if p.Threshold <= 0 {
		p.Threshold = def.Threshold
	}
	if p.Window <= 0 {
		p.Window = def.Window
	}
	return &Quorum{policy: p, votes: make(map[string]vote)}
}

// Cast records a yes/no vote from the named voter.
func (q *Quorum) Cast(voter string, yes bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.votes[voter] = vote{value: yes, at: time.Now()}
}

// Reached reports whether quorum has been achieved.
// Expired votes (outside the configured window) are ignored.
func (q *Quorum) Reached() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	var total, yes int
	for _, v := range q.votes {
		if now.Sub(v.at) > q.policy.Window {
			continue
		}
		total++
		if v.value {
			yes++
		}
	}
	if total < q.policy.MinVoters {
		return false
	}
	return float64(yes)/float64(total) >= q.policy.Threshold
}

// Reset clears all recorded votes.
func (q *Quorum) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.votes = make(map[string]vote)
}

// VoterCount returns the number of active (non-expired) votes.
func (q *Quorum) VoterCount() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	now := time.Now()
	count := 0
	for _, v := range q.votes {
		if now.Sub(v.at) <= q.policy.Window {
			count++
		}
	}
	return count
}
