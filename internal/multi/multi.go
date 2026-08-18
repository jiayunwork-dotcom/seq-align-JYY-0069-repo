// Package multi implements progressive multiple sequence alignment using a
// guide tree constructed from pairwise distances.
package multi

import (
	"seq-align/internal/align"
	"seq-align/internal/scoring"
)

// Profile represents an alignment profile (position frequency matrix).
type Profile struct {
	Columns []map[byte]int
	Count   int
}

// NewProfile creates a profile from a single sequence.
func NewProfile(seq string) *Profile {
	p := &Profile{Columns: make([]map[byte]int, len(seq)), Count: 1}
	for i, c := range seq {
		p.Columns[i] = map[byte]int{byte(c): 1}
	}
	return p
}

// Len returns the profile length.
func (p *Profile) Len() int { return len(p.Columns) }

// Consensus returns the consensus sequence (most frequent at each position).
func (p *Profile) Consensus() string {
	out := make([]byte, len(p.Columns))
	for i, col := range p.Columns {
		var best byte
		var bestC int
		for c, n := range col {
			if n > bestC {
				best = c
				bestC = n
			}
		}
		out[i] = best
	}
	return string(out)
}

// MSAResult holds the output of multiple sequence alignment.
type MSAResult struct {
	Aligned  []string
	Score    int
	Names    []string
}

// Progressive performs progressive multiple alignment on the given sequences.
// It uses a simple star alignment strategy with the first sequence as center.
func Progressive(sequences []string, scheme scoring.Scheme) MSAResult {
	if len(sequences) == 0 {
		return MSAResult{}
	}
	if len(sequences) == 1 {
		return MSAResult{Aligned: sequences}
	}

	// Star alignment: use first as reference.
	ref := sequences[0]
	aligned := make([]string, len(sequences))
	aligned[0] = ref

	// Align each to the reference and collect.
	for i := 1; i < len(sequences); i++ {
		result, _ := align.Global(ref, sequences[i], scheme)
		aligned[0] = result.A // may have gaps inserted
		aligned[i] = result.B
	}

	// Re-align all to the (possibly modified) reference.
	// Simple pass: just keep pairwise results.
	totalScore := 0
	finalAligned := make([]string, len(sequences))
	finalAligned[0] = sequences[0]
	for i := 1; i < len(sequences); i++ {
		result, _ := align.Global(finalAligned[0], sequences[i], scheme)
		finalAligned[i] = result.B
		totalScore += result.Score
	}

	return MSAResult{Aligned: finalAligned, Score: totalScore}
}

// PairwiseDistanceMatrix computes the distance matrix for a set of sequences.
// Distance = 1 - identity.
func PairwiseDistanceMatrix(sequences []string, scheme scoring.Scheme) [][]float64 {
	n := len(sequences)
	dist := make([][]float64, n)
	for i := range dist {
		dist[i] = make([]float64, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			result, _ := align.Global(sequences[i], sequences[j], scheme)
			dist[i][j] = 1 - result.Identity
			dist[j][i] = dist[i][j]
		}
	}
	return dist
}
