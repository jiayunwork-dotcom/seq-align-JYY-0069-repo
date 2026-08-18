package stats

import "testing"

func TestSeqStats(t *testing.T) {
	s := ComputeSeqStats("ACGTACGT")
	if s.Length != 8 {
		t.Fatalf("expected 8, got %d", s.Length)
	}
	if s.GC != 0.5 {
		t.Fatalf("expected GC=0.5, got %f", s.GC)
	}
}

func TestAlignStats(t *testing.T) {
	s := ComputeAlignStats("ACGT-A", "ACCT-A")
	// A=A(match), C=C(match), G vs C(mismatch), T=T(match), -=-(gap), A=A(match)
	if s.Matches != 4 {
		t.Fatalf("expected 4 matches, got %d", s.Matches)
	}
	if s.Gaps != 1 { // position 4: both are '-', counted as gap
		t.Fatalf("expected 1 gap, got %d", s.Gaps)
	}
}

func TestMatchLine(t *testing.T) {
	line := MatchLine("ACGT", "ACCT")
	if line != "||.|" {
		t.Fatalf("expected '||.|', got %q", line)
	}
}

func TestConsensus(t *testing.T) {
	c := ConsensusSequence("AC-T", "ACGT")
	if c != "ACGT" {
		t.Fatalf("expected ACGT, got %q", c)
	}
}
