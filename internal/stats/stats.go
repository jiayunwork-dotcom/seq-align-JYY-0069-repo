// Package stats computes sequence and alignment statistics.
package stats

import "strings"

// SeqStats holds basic sequence statistics.
type SeqStats struct {
	Length    int
	GC       float64
	Composition map[byte]int
}

// ComputeSeqStats computes stats for a nucleotide sequence.
func ComputeSeqStats(seq string) SeqStats {
	s := SeqStats{Length: len(seq), Composition: make(map[byte]int)}
	var gc int
	for i := 0; i < len(seq); i++ {
		c := seq[i]
		s.Composition[c]++
		if c == 'G' || c == 'C' || c == 'g' || c == 'c' {
			gc++
		}
	}
	if s.Length > 0 {
		s.GC = float64(gc) / float64(s.Length)
	}
	return s
}

// AlignStats holds alignment quality statistics.
type AlignStats struct {
	Length       int
	Matches     int
	Mismatches  int
	Gaps        int
	Identity    float64
	Similarity  float64
	GapFraction float64
}

// ComputeAlignStats computes stats from aligned sequences.
func ComputeAlignStats(alignedA, alignedB string) AlignStats {
	n := len(alignedA)
	if len(alignedB) < n {
		n = len(alignedB)
	}
	s := AlignStats{Length: n}
	for i := 0; i < n; i++ {
		a, b := alignedA[i], alignedB[i]
		if a == '-' || b == '-' {
			s.Gaps++
		} else if a == b {
			s.Matches++
		} else {
			s.Mismatches++
		}
	}
	if n > 0 {
		s.Identity = float64(s.Matches) / float64(n)
		s.Similarity = float64(s.Matches+s.Mismatches) / float64(n) // non-gap positions
		s.GapFraction = float64(s.Gaps) / float64(n)
	}
	return s
}

// ConsensusSequence builds a consensus from two aligned sequences.
func ConsensusSequence(alignedA, alignedB string) string {
	n := len(alignedA)
	if len(alignedB) < n {
		n = len(alignedB)
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		a, b := alignedA[i], alignedB[i]
		if a == b {
			sb.WriteByte(a)
		} else if a == '-' {
			sb.WriteByte(b)
		} else if b == '-' {
			sb.WriteByte(a)
		} else {
			sb.WriteByte('X') // ambiguity
		}
	}
	return sb.String()
}

// MatchLine creates a visual match indicator (| for match, . for mismatch, space for gap).
func MatchLine(alignedA, alignedB string) string {
	n := len(alignedA)
	if len(alignedB) < n {
		n = len(alignedB)
	}
	line := make([]byte, n)
	for i := 0; i < n; i++ {
		a, b := alignedA[i], alignedB[i]
		if a == '-' || b == '-' {
			line[i] = ' '
		} else if a == b {
			line[i] = '|'
		} else {
			line[i] = '.'
		}
	}
	return string(line)
}
