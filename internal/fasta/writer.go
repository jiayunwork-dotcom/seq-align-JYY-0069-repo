package fasta

import (
	"fmt"
	"io"
	"strings"
)

// Sequence is a named sequence record.
type Sequence struct {
	Name string
	Seq  string
}

// Write writes sequences in FASTA format to w.
func Write(w io.Writer, seqs []Sequence, lineWidth int) error {
	if lineWidth <= 0 {
		lineWidth = 60
	}
	for _, s := range seqs {
		if _, err := fmt.Fprintf(w, ">%s\n", s.Name); err != nil {
			return err
		}
		for i := 0; i < len(s.Seq); i += lineWidth {
			end := i + lineWidth
			if end > len(s.Seq) {
				end = len(s.Seq)
			}
			if _, err := fmt.Fprintf(w, "%s\n", s.Seq[i:end]); err != nil {
				return err
			}
		}
	}
	return nil
}

// FormatOne formats a single sequence.
func FormatOne(name, seq string) string {
	var sb strings.Builder
	sb.WriteString(">")
	sb.WriteString(name)
	sb.WriteByte('\n')
	for i := 0; i < len(seq); i += 60 {
		end := i + 60
		if end > len(seq) {
			end = len(seq)
		}
		sb.WriteString(seq[i:end])
		sb.WriteByte('\n')
	}
	return sb.String()
}

// ReverseComplement returns the reverse complement of a DNA sequence.
func ReverseComplement(seq string) string {
	n := len(seq)
	out := make([]byte, n)
	for i, c := range seq {
		out[n-1-i] = complement(byte(c))
	}
	return string(out)
}

func complement(c byte) byte {
	switch c {
	case 'A', 'a':
		return 'T'
	case 'T', 't':
		return 'A'
	case 'G', 'g':
		return 'C'
	case 'C', 'c':
		return 'G'
	default:
		return 'N'
	}
}

// Translate converts a DNA sequence to protein (standard genetic code).
func Translate(dna string) string {
	var sb strings.Builder
	for i := 0; i+2 < len(dna); i += 3 {
		codon := strings.ToUpper(dna[i : i+3])
		sb.WriteByte(codonToAA(codon))
	}
	return sb.String()
}

func codonToAA(codon string) byte {
	table := map[string]byte{
		"TTT": 'F', "TTC": 'F', "TTA": 'L', "TTG": 'L',
		"CTT": 'L', "CTC": 'L', "CTA": 'L', "CTG": 'L',
		"ATT": 'I', "ATC": 'I', "ATA": 'I', "ATG": 'M',
		"GTT": 'V', "GTC": 'V', "GTA": 'V', "GTG": 'V',
		"TCT": 'S', "TCC": 'S', "TCA": 'S', "TCG": 'S',
		"CCT": 'P', "CCC": 'P', "CCA": 'P', "CCG": 'P',
		"ACT": 'T', "ACC": 'T', "ACA": 'T', "ACG": 'T',
		"GCT": 'A', "GCC": 'A', "GCA": 'A', "GCG": 'A',
		"TAT": 'Y', "TAC": 'Y', "TAA": '*', "TAG": '*',
		"CAT": 'H', "CAC": 'H', "CAA": 'Q', "CAG": 'Q',
		"AAT": 'N', "AAC": 'N', "AAA": 'K', "AAG": 'K',
		"GAT": 'D', "GAC": 'D', "GAA": 'E', "GAG": 'E',
		"TGT": 'C', "TGC": 'C', "TGA": '*', "TGG": 'W',
		"CGT": 'R', "CGC": 'R', "CGA": 'R', "CGG": 'R',
		"AGT": 'S', "AGC": 'S', "AGA": 'R', "AGG": 'R',
		"GGT": 'G', "GGC": 'G', "GGA": 'G', "GGG": 'G',
	}
	if aa, ok := table[codon]; ok {
		return aa
	}
	return 'X'
}
