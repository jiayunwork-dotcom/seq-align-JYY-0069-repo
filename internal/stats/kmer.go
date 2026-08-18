// Package stats - k-mer frequency analysis for sequence comparison without alignment.
package stats

import "math"

// KmerProfile holds k-mer frequencies for a sequence.
type KmerProfile struct {
	K      int
	Counts map[string]int
	Total  int
}

// ComputeKmerProfile builds a k-mer frequency profile for a sequence.
// k must be >= 1; returns an empty profile for k > len(seq).
func ComputeKmerProfile(seq string, k int) KmerProfile {
	p := KmerProfile{K: k, Counts: make(map[string]int)}
	if k < 1 || k > len(seq) {
		return p
	}
	for i := 0; i <= len(seq)-k; i++ {
		kmer := seq[i : i+k]
		p.Counts[kmer]++
		p.Total++
	}
	return p
}

// KmerDistance computes the Euclidean distance between two k-mer profiles.
// Useful for fast pre-filtering before expensive alignment.
func KmerDistance(a, b KmerProfile) float64 {
	allKmers := make(map[string]bool)
	for k := range a.Counts {
		allKmers[k] = true
	}
	for k := range b.Counts {
		allKmers[k] = true
	}
	sumSq := 0.0
	for kmer := range allKmers {
		freqA := float64(a.Counts[kmer]) / math.Max(float64(a.Total), 1)
		freqB := float64(b.Counts[kmer]) / math.Max(float64(b.Total), 1)
		diff := freqA - freqB
		sumSq += diff * diff
	}
	return math.Sqrt(sumSq)
}

// TopKmers returns the n most frequent k-mers in the profile.
// If n > number of distinct k-mers, returns all.
func TopKmers(p KmerProfile, n int) []string {
	type kv struct {
		kmer  string
		count int
	}
	var pairs []kv
	for k, c := range p.Counts {
		pairs = append(pairs, kv{k, c})
	}
	// 简单选择排序，适合小规模数据
	for i := 0; i < len(pairs) && i < n; i++ {
		maxIdx := i
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].count > pairs[maxIdx].count {
				maxIdx = j
			}
		}
		pairs[i], pairs[maxIdx] = pairs[maxIdx], pairs[i]
	}
	limit := n
	if limit > len(pairs) {
		limit = len(pairs)
	}
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		result[i] = pairs[i].kmer
	}
	return result
}
