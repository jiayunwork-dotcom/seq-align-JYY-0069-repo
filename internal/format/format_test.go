package format

import (
	"strings"
	"testing"
)

func TestFASTA(t *testing.T) {
	out := FASTA("seq1", "ACGTACGTACGT", 4)
	if !strings.HasPrefix(out, ">seq1\n") {
		t.Fatal("should start with >seq1")
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 4 { // header + 3 lines of 4
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}
}

func TestClustal(t *testing.T) {
	out := Clustal("seqA", "ACGT-A", "seqB", "AC-TGA", 60)
	if !strings.Contains(out, "CLUSTAL") {
		t.Fatal("should contain CLUSTAL header")
	}
	if !strings.Contains(out, "seqA") {
		t.Fatal("should contain seqA name")
	}
}

func TestMatchBar(t *testing.T) {
	bar := MatchBar("ACGT", "ACCT")
	if bar != "||.|" {
		t.Fatalf("expected '||.|', got %q", bar)
	}
}

func TestPlain(t *testing.T) {
	out := Plain("a", "ACGT", "b", "ACCT")
	if !strings.Contains(out, "a: ACGT") {
		t.Fatal("should contain labeled output")
	}
}
