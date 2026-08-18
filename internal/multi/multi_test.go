package multi

import (
	"seq-align/internal/scoring"
	"testing"
)

func TestProgressiveBasic(t *testing.T) {
	seqs := []string{"ACGT", "ACCT", "ACGG"}
	scheme := scoring.Scheme{Match: 1, Mismatch: -1, GapOpen: -2, GapExtend: -1}
	result := Progressive(seqs, scheme)
	if len(result.Aligned) != 3 {
		t.Fatalf("expected 3, got %d", len(result.Aligned))
	}
}

func TestProgressiveSingle(t *testing.T) {
	result := Progressive([]string{"ACGT"}, scoring.Scheme{})
	if len(result.Aligned) != 1 {
		t.Fatal("expected 1")
	}
}

func TestPairwiseDistanceMatrix(t *testing.T) {
	seqs := []string{"ACGT", "ACGT", "TTTT"}
	scheme := scoring.Scheme{Match: 1, Mismatch: -1, GapOpen: -2, GapExtend: -1}
	dist := PairwiseDistanceMatrix(seqs, scheme)
	// Identical sequences should have small distance (identity close to 1).
	if dist[0][1] > 0.01 {
		t.Fatalf("identical seqs should have near-zero dist, got %f", dist[0][1])
	}
	// Different sequences should have larger distance.
	if dist[0][2] <= dist[0][1] {
		t.Fatal("different seqs should have larger distance")
	}
}

func TestProfile(t *testing.T) {
	p := NewProfile("ACGT")
	if p.Len() != 4 {
		t.Fatal("expected 4")
	}
	if p.Consensus() != "ACGT" {
		t.Fatalf("expected ACGT, got %s", p.Consensus())
	}
}
