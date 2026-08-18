package fasta

import "testing"

func TestReverseComplement(t *testing.T) {
	rc := ReverseComplement("ACGT")
	if rc != "ACGT" { // palindrome
		t.Fatalf("expected ACGT, got %s", rc)
	}
	rc2 := ReverseComplement("AACG")
	if rc2 != "CGTT" {
		t.Fatalf("expected CGTT, got %s", rc2)
	}
}

func TestTranslate(t *testing.T) {
	// ATG = M, GCT = A, TAA = *
	protein := Translate("ATGGCTTAA")
	if protein != "MA*" {
		t.Fatalf("expected MA*, got %s", protein)
	}
}

func TestFormatOne(t *testing.T) {
	out := FormatOne("test", "ACGT")
	if out != ">test\nACGT\n" {
		t.Fatalf("unexpected: %q", out)
	}
}
