package align

import "seq-align/internal/scoring"

// SemiGlobal performs semi-global alignment (no penalty for leading/trailing gaps).
// Useful for aligning a short read to a longer reference.
func SemiGlobal(query, reference string, s scoring.Scheme) (Result, error) {
	m := len(query)
	n := len(reference)
	if m == 0 || n == 0 {
		return Result{}, nil
	}

	// Initialize DP matrix.
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	// No penalty for leading gaps in reference (free start in query direction).
	for i := 1; i <= m; i++ {
		dp[i][0] = i * s.Gap
	}
	// Free start in reference: dp[0][j] = 0 (no penalty for leading gaps in reference).

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			match := dp[i-1][j-1]
			if query[i-1] == reference[j-1] {
				match += s.Match
			} else {
				match += s.Mismatch
			}
			del := dp[i-1][j] + s.Gap
			ins := dp[i][j-1] + s.Gap
			dp[i][j] = max3(match, del, ins)
		}
	}

	// Find best score in last row (free trailing gaps in reference).
	bestScore := dp[m][0]
	bestJ := 0
	for j := 1; j <= n; j++ {
		if dp[m][j] > bestScore {
			bestScore = dp[m][j]
			bestJ = j
		}
	}

	// Traceback.
	var alignA, alignB []byte
	i, j := m, bestJ
	for i > 0 && j > 0 {
		score := dp[i][j]
		diag := dp[i-1][j-1]
		matchScore := s.Mismatch
		if query[i-1] == reference[j-1] {
			matchScore = s.Match
		}
		if score == diag+matchScore {
			alignA = append(alignA, query[i-1])
			alignB = append(alignB, reference[j-1])
			i--
			j--
		} else if score == dp[i-1][j]+s.Gap {
			alignA = append(alignA, query[i-1])
			alignB = append(alignB, '-')
			i--
		} else {
			alignA = append(alignA, '-')
			alignB = append(alignB, reference[j-1])
			j--
		}
	}
	for i > 0 {
		alignA = append(alignA, query[i-1])
		alignB = append(alignB, '-')
		i--
	}
	for j > 0 {
		alignA = append(alignA, '-')
		alignB = append(alignB, reference[j-1])
		j--
	}

	reverse(alignA)
	reverse(alignB)

	return buildResult(alignA, alignB, bestScore, 0, 0), nil
}

func max3(a, b, c int) int {
	if b > a {
		a = b
	}
	if c > a {
		a = c
	}
	return a
}
