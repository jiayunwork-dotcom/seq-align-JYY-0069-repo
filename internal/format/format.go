// Package format provides output formatting for sequence alignments in
// various standard formats (FASTA, CLUSTAL, plain text).
package format

import (
	"fmt"
	"strings"
)

// FASTA formats a sequence in FASTA format with line wrapping.
func FASTA(name, seq string, lineWidth int) string {
	if lineWidth <= 0 {
		lineWidth = 60
	}
	var sb strings.Builder
	sb.WriteString(">")
	sb.WriteString(name)
	sb.WriteByte('\n')
	for i := 0; i < len(seq); i += lineWidth {
		end := i + lineWidth
		if end > len(seq) {
			end = len(seq)
		}
		sb.WriteString(seq[i:end])
		sb.WriteByte('\n')
	}
	return sb.String()
}

// Clustal formats an alignment in CLUSTAL-like format.
func Clustal(nameA, seqA, nameB, seqB string, lineWidth int) string {
	if lineWidth <= 0 {
		lineWidth = 60
	}
	var sb strings.Builder
	sb.WriteString("CLUSTAL-like alignment\n\n")

	maxName := len(nameA)
	if len(nameB) > maxName {
		maxName = len(nameB)
	}
	pad := maxName + 4

	n := len(seqA)
	if len(seqB) < n {
		n = len(seqB)
	}

	for start := 0; start < n; start += lineWidth {
		end := start + lineWidth
		if end > n {
			end = n
		}
		sb.WriteString(padRight(nameA, pad))
		sb.WriteString(seqA[start:end])
		sb.WriteByte('\n')

		// Match line.
		sb.WriteString(padRight("", pad))
		for i := start; i < end; i++ {
			if seqA[i] == seqB[i] && seqA[i] != '-' {
				sb.WriteByte('*')
			} else if seqA[i] != '-' && seqB[i] != '-' {
				sb.WriteByte('.')
			} else {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')

		sb.WriteString(padRight(nameB, pad))
		sb.WriteString(seqB[start:end])
		sb.WriteByte('\n')
		sb.WriteByte('\n')
	}
	return sb.String()
}

// Plain formats an alignment in simple labeled text.
func Plain(nameA, seqA, nameB, seqB string) string {
	return fmt.Sprintf("%s: %s\n%s: %s\n", nameA, seqA, nameB, seqB)
}

// MatchBar creates a visual match indicator string.
func MatchBar(seqA, seqB string) string {
	n := len(seqA)
	if len(seqB) < n {
		n = len(seqB)
	}
	bar := make([]byte, n)
	for i := 0; i < n; i++ {
		if seqA[i] == seqB[i] && seqA[i] != '-' {
			bar[i] = '|'
		} else if seqA[i] == '-' || seqB[i] == '-' {
			bar[i] = ' '
		} else {
			bar[i] = '.'
		}
	}
	return string(bar)
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
