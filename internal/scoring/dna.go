// Package scoring - DNA-specific substitution matrices for nucleotide alignment.
package scoring

// DNAIdentity returns a simple identity scoring scheme for DNA:
// match=1, mismatch=-1, gap=-2.
func DNAIdentity() Scheme {
	return Scheme{Match: 1, Mismatch: -1, Gap: -2}
}

// DNATransitionTransversion returns a scheme that penalizes transversions
// more heavily than transitions. Uses affine gaps.
// Transitions (A<->G, C<->T) get mismatch=-1, transversions get -2.
// This is approximated with mismatch=-2, suitable for divergent sequences.
func DNATransitionTransversion() Scheme {
	return Scheme{
		Match:     2,
		Mismatch:  -2,
		GapOpen:   -5,
		GapExtend: -1,
	}
}

// DNAHighSimilarity returns a scheme optimized for closely related sequences
// (>90% identity), with a tight gap penalty to avoid spurious indels.
func DNAHighSimilarity() Scheme {
	return Scheme{
		Match:     1,
		Mismatch:  -3,
		GapOpen:   -8,
		GapExtend: -2,
	}
}

// IsTransition reports whether two nucleotides form a transition pair.
// Transitions: A<->G (purines) and C<->T (pyrimidines).
func IsTransition(a, b byte) bool {
	a = toUpper(a)
	b = toUpper(b)
	if a == b {
		return false
	}
	return (a == 'A' && b == 'G') || (a == 'G' && b == 'A') ||
		(a == 'C' && b == 'T') || (a == 'T' && b == 'C')
}

// IsTransversion reports whether two nucleotides form a transversion pair.
func IsTransversion(a, b byte) bool {
	a = toUpper(a)
	b = toUpper(b)
	if a == b {
		return false
	}
	return !IsTransition(a, b)
}

func toUpper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 32
	}
	return c
}
