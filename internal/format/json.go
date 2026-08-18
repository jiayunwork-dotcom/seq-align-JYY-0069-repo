// Package format - JSON output format for alignment results.
package format

import (
	"encoding/json"
)

// JSONAlignment represents an alignment result in JSON-serializable form.
type JSONAlignment struct {
	SeqA     string  `json:"seq_a"`
	SeqB     string  `json:"seq_b"`
	Score    int     `json:"score"`
	Identity float64 `json:"identity"`
	Length   int     `json:"length"`
	Gaps     int     `json:"gaps"`
	Mode     string  `json:"mode"`
}

// JSONReport holds a complete alignment report with metadata.
type JSONReport struct {
	Version    string          `json:"version"`
	Alignments []JSONAlignment `json:"alignments"`
	Summary    JSONSummary     `json:"summary"`
}

// JSONSummary provides aggregate stats for a multi-alignment report.
type JSONSummary struct {
	TotalPairs    int     `json:"total_pairs"`
	MeanIdentity  float64 `json:"mean_identity"`
	MeanScore     float64 `json:"mean_score"`
}

// RenderJSON formats a list of alignments as a JSON report string.
func RenderJSON(alignments []JSONAlignment) (string, error) {
	if len(alignments) == 0 {
		return "[]", nil
	}
	var totalIdentity, totalScore float64
	for _, a := range alignments {
		totalIdentity += a.Identity
		totalScore += float64(a.Score)
	}
	n := float64(len(alignments))
	report := JSONReport{
		Version:    "1.0",
		Alignments: alignments,
		Summary: JSONSummary{
			TotalPairs:   len(alignments),
			MeanIdentity: totalIdentity / n,
			MeanScore:    totalScore / n,
		},
	}
	b, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ParseJSONReport parses a JSON report string back into a JSONReport struct.
func ParseJSONReport(data []byte) (*JSONReport, error) {
	var report JSONReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, err
	}
	return &report, nil
}
