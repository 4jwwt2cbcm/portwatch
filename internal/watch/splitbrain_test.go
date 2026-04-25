package watch

import (
	"testing"
	"time"
)

func makeSplitBrain(quorum int, window time.Duration) *SplitBrain {
	return NewSplitBrain(SplitBrainPolicy{
		QuorumSize: quorum,
		Window:     window,
	})
}

func TestDefaultSplitBrainPolicyValues(t *testing.T) {
	p := DefaultSplitBrainPolicy()
	if p.QuorumSize != 2 {
		t.Errorf("expected QuorumSize=2, got %d", p.QuorumSize)
	}
	if p.Window != 10*time.Second {
		t.Errorf("expected Window=10s, got %v", p.Window)
	}
}

func TestSplitBrainNoVotesNoQuorum(t *testing.T) {
	sb := makeSplitBrain(2, time.Second)
	if sb.HasQuorum() {
		t.Error("expected no quorum with no votes")
	}
}

func TestSplitBrainQuorumReached(t *testing.T) {
	sb := makeSplitBrain(2, time.Second)
	sb.Vote("a")
	sb.Vote("a")
	if !sb.HasQuorum() {
		t.Error("expected quorum after 2 votes for same value")
	}
}

func TestSplitBrainQuorumNotReachedWithSplit(t *testing.T) {
	sb := makeSplitBrain(2, time.Second)
	sb.Vote("a")
	sb.Vote("b")
	if sb.HasQuorum() {
		t.Error("expected no quorum when votes are split")
	}
}

func TestSplitBrainConflictedDetected(t *testing.T) {
	sb := makeSplitBrain(2, time.Second)
	sb.Vote("a")
	sb.Vote("b")
	if !sb.Conflicted() {
		t.Error("expected conflict when two values have votes")
	}
}

func TestSplitBrainNotConflictedWithSingleValue(t *testing.T) {
	sb := makeSplitBrain(2, time.Second)
	sb.Vote("a")
	sb.Vote("a")
	if sb.Conflicted() {
		t.Error("expected no conflict when all votes agree")
	}
}

func TestSplitBrainResetClearsState(t *testing.T) {
	sb := makeSplitBrain(2, time.Second)
	sb.Vote("a")
	sb.Vote("a")
	sb.Reset()
	if sb.HasQuorum() {
		t.Error("expected no quorum after reset")
	}
	if sb.Conflicted() {
		t.Error("expected no conflict after reset")
	}
}

func TestSplitBrainVotesExpireAfterWindow(t *testing.T) {
	sb := makeSplitBrain(2, 10*time.Millisecond)
	sb.Vote("a")
	sb.Vote("a")
	time.Sleep(20 * time.Millisecond)
	if sb.HasQuorum() {
		t.Error("expected quorum to expire after window")
	}
}

func TestSplitBrainDefaultsOnZeroPolicy(t *testing.T) {
	sb := NewSplitBrain(SplitBrainPolicy{})
	if sb.policy.QuorumSize != DefaultSplitBrainPolicy().QuorumSize {
		t.Errorf("expected default QuorumSize, got %d", sb.policy.QuorumSize)
	}
	if sb.policy.Window != DefaultSplitBrainPolicy().Window {
		t.Errorf("expected default Window, got %v", sb.policy.Window)
	}
}
