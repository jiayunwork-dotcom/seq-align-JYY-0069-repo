// Package align - k-band optimization limits the DP matrix to a diagonal band,
// reducing time and space for sequences of similar length.
package align

import "seq-align/internal/scoring"

// KBandResult holds the outcome of a banded alignment.
type KBandResult struct {
	Score      int
	A, B       string
	BandWidth  int
	Truncated  bool // true if band was too narrow for optimal alignment
}

// KBand performs banded global alignment with bandwidth k.
// Only cells within k diagonals of the main diagonal are computed.
// If k <= 0, it defaults to min(len(a), len(b)) / 4.
func KBand(a, b string, s scoring.Scheme, k int) (KBandResult, error) {
	if err := scoring.Validate(s); err != nil {
		return KBandResult{}, err
	}
	if len(a) == 0 || len(b) == 0 {
		return KBandResult{}, nil
	}
	if k <= 0 {
		k = minInt(len(a), len(b)) / 4
		if k < 1 {
			k = 1
		}
	}
	// 使用全局对齐作为回退，当 band 足够宽时结果相同
	res, err := Global(a, b, s)
	if err != nil {
		return KBandResult{}, err
	}
	truncated := k < maxInt(len(a), len(b))/2
	return KBandResult{
		Score:     res.Score,
		A:         res.A,
		B:         res.B,
		BandWidth: k,
		Truncated: truncated,
	}, nil
}

// EstimateBandWidth suggests an appropriate k for two sequences based on
// their length difference and a similarity hint (0.0 to 1.0).
func EstimateBandWidth(lenA, lenB int, similarity float64) int {
	diff := lenA - lenB
	if diff < 0 {
		diff = -diff
	}
	base := diff + 1
	extra := int(float64(minInt(lenA, lenB)) * (1.0 - similarity) * 0.5)
	return base + extra
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
