package watch

import (
	"testing"
	"time"
)

func makeQuorum(min int, threshold float64, window time.Duration) *Quorum {
	return NewQuorum(QuorumPolicy{
		MinVoters: min,
		Threshold: threshold,
		Window:    window,
	})
}

func TestDefaultQuorumPolicyValues(t *testing.T) {
	p := DefaultQuorumPolicy()
	if p.MinVoters != 3 {
		t.Errorf("expected MinVoters 3, got %d", p.MinVoters)
	}
	if p.Threshold != 0.51 {
		t.Errorf("expected Threshold 0.51, got %f", p.Threshold)
	}
	if p.Window != 30*time.Second {
		t.Errorf("expected Window 30s, got %v", p.Window)
	}
}

func TestQuorumNotReachedWithNoVotes(t *testing.T) {
	q := makeQuorum(1, 0.5, time.Minute)
	if q.Reached() {
		t.Error("expected quorum not reached with no votes")
	}
}

func TestQuorumReachedWithMajority(t *testing.T) {
	q := makeQuorum(3, 0.5, time.Minute)
	q.Cast("a", true)
	q.Cast("b", true)
	q.Cast("c", false)
	if !q.Reached() {
		t.Error("expected quorum reached with 2/3 yes votes")
	}
}

func TestQuorumNotReachedBelowThreshold(t *testing.T) {
	q := makeQuorum(3, 0.67, time.Minute)
	q.Cast("a", true)
	q.Cast("b", false)
	q.Cast("c", false)
	if q.Reached() {
		t.Error("expected quorum not reached with 1/3 yes votes")
	}
}

func TestQuorumNotReachedBelowMinVoters(t *testing.T) {
	q := makeQuorum(3, 0.5, time.Minute)
	q.Cast("a", true)
	q.Cast("b", true)
	if q.Reached() {
		t.Error("expected quorum not reached: only 2 of 3 required voters")
	}
}

func TestQuorumExpiredVotesIgnored(t *testing.T) {
	q := makeQuorum(1, 0.5, time.Millisecond)
	q.Cast("a", true)
	time.Sleep(5 * time.Millisecond)
	if q.Reached() {
		t.Error("expected quorum not reached after vote expiry")
	}
}

func TestQuorumResetClearsVotes(t *testing.T) {
	q := makeQuorum(1, 0.5, time.Minute)
	q.Cast("a", true)
	q.Reset()
	if q.VoterCount() != 0 {
		t.Errorf("expected 0 voters after reset, got %d", q.VoterCount())
	}
}

func TestQuorumVoterCountExcludesExpired(t *testing.T) {
	q := makeQuorum(1, 0.5, 10*time.Millisecond)
	q.Cast("a", true)
	q.Cast("b", false)
	time.Sleep(20 * time.Millisecond)
	q.Cast("c", true)
	if n := q.VoterCount(); n != 1 {
		t.Errorf("expected 1 active voter, got %d", n)
	}
}

func TestQuorumDefaultsAppliedOnZero(t *testing.T) {
	q := NewQuorum(QuorumPolicy{})
	def := DefaultQuorumPolicy()
	if q.policy.MinVoters != def.MinVoters {
		t.Errorf("expected default MinVoters %d, got %d", def.MinVoters, q.policy.MinVoters)
	}
}
