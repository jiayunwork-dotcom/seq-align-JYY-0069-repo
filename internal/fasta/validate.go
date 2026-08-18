// Package fasta - sequence validation functions for FASTA records.
package fasta

import (
	"fmt"
	"strings"
)

// SeqType represents the type of biological sequence.
type SeqType int

const (
	SeqDNA     SeqType = iota // DNA nucleotides: A, C, G, T
	SeqRNA                    // RNA nucleotides: A, C, G, U
	SeqProtein                // Amino acid single-letter codes
)

// validChars defines allowed characters for each sequence type (uppercase).
var validChars = map[SeqType]string{
	SeqDNA:     "ACGTNRYSWKMBDHV",
	SeqRNA:     "ACGUNRYSWKMBDHV",
	SeqProtein: "ACDEFGHIKLMNPQRSTVWY*X",
}

// ValidateSequence checks that a sequence contains only valid characters
// for the given type. Returns nil if valid, or an error describing the
// first invalid character found.
func ValidateSequence(seq string, st SeqType) error {
	allowed, ok := validChars[st]
	if !ok {
		return fmt.Errorf("fasta: unknown sequence type %d", st)
	}
	upper := strings.ToUpper(seq)
	for i, c := range upper {
		if !strings.ContainsRune(allowed, c) {
			return fmt.Errorf("fasta: invalid character %q at position %d for type %d", c, i, st)
		}
	}
	return nil
}

// DetectType guesses the sequence type based on character composition.
// Returns SeqDNA if >80% ACGT, SeqRNA if contains U but no T, else SeqProtein.
func DetectType(seq string) SeqType {
	upper := strings.ToUpper(seq)
	var dnaCount, uCount, tCount int
	for _, c := range upper {
		switch c {
		case 'A', 'C', 'G':
			dnaCount++
		case 'T':
			dnaCount++
			tCount++
		case 'U':
			uCount++
		}
	}
	if len(upper) == 0 {
		return SeqDNA
	}
	if uCount > 0 && tCount == 0 {
		return SeqRNA
	}
	if float64(dnaCount+tCount)/float64(len(upper)) > 0.8 {
		return SeqDNA
	}
	return SeqProtein
}

// ValidateRecords checks a slice of Records, ensuring each has a non-empty
// ID and a valid sequence of the given type.
func ValidateRecords(records []Record, st SeqType) error {
	for i, r := range records {
		if r.ID == "" {
			return fmt.Errorf("fasta: record %d has empty ID", i)
		}
		if r.Seq == "" {
			return fmt.Errorf("fasta: record %q has empty sequence", r.ID)
		}
		if err := ValidateSequence(r.Seq, st); err != nil {
			return fmt.Errorf("fasta: record %q: %w", r.ID, err)
		}
	}
	return nil
}
