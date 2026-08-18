// Package matrix - sparse matrix utilities for large distance matrices
// where most values are above a threshold and can be omitted.
package matrix

import "fmt"

// SparseEntry represents a single non-zero entry in a sparse matrix.
type SparseEntry struct {
	Row, Col int
	Value    float64
}

// SparseMatrix stores only entries below a distance threshold.
type SparseMatrix struct {
	Dim     int
	IDs     []string
	Entries []SparseEntry
}

// Sparsify converts a dense Matrix into a SparseMatrix, keeping only entries
// where value <= threshold. Diagonal entries are always excluded.
func Sparsify(m *Matrix, threshold float64) *SparseMatrix {
	n := m.Dim()
	sm := &SparseMatrix{Dim: n, IDs: make([]string, n)}
	copy(sm.IDs, m.IDs)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if m.Data[i][j] <= threshold {
				sm.Entries = append(sm.Entries, SparseEntry{Row: i, Col: j, Value: m.Data[i][j]})
			}
		}
	}
	return sm
}

// Density returns the fraction of upper-triangle entries that are stored.
func (sm *SparseMatrix) Density() float64 {
	total := sm.Dim * (sm.Dim - 1) / 2
	if total == 0 {
		return 0
	}
	return float64(len(sm.Entries)) / float64(total)
}

// Get retrieves the distance between i and j. Returns 0 if not stored,
// and a boolean indicating whether the entry exists.
func (sm *SparseMatrix) Get(i, j int) (float64, bool) {
	if i > j {
		i, j = j, i
	}
	for _, e := range sm.Entries {
		if e.Row == i && e.Col == j {
			return e.Value, true
		}
	}
	return 0, false
}

// Neighbors returns all IDs within threshold distance of the given index.
func (sm *SparseMatrix) Neighbors(idx int) []string {
	var result []string
	for _, e := range sm.Entries {
		if e.Row == idx {
			result = append(result, sm.IDs[e.Col])
		} else if e.Col == idx {
			result = append(result, sm.IDs[e.Row])
		}
	}
	return result
}

// String returns a summary of the sparse matrix.
func (sm *SparseMatrix) String() string {
	return fmt.Sprintf("SparseMatrix{dim=%d, entries=%d, density=%.2f%%}",
		sm.Dim, len(sm.Entries), sm.Density()*100)
}
