package scoring

import "testing"

func TestBlosumScore(t *testing.T) {
	// A-A should be 4.
	if BlosumScore('A', 'A') != 4 {
		t.Fatalf("expected 4, got %d", BlosumScore('A', 'A'))
	}
	// A-R should be -1.
	if BlosumScore('A', 'R') != -1 {
		t.Fatalf("expected -1, got %d", BlosumScore('A', 'R'))
	}
	// C-C should be 9.
	if BlosumScore('C', 'C') != 9 {
		t.Fatalf("expected 9, got %d", BlosumScore('C', 'C'))
	}
}

func TestProteinScheme(t *testing.T) {
	s := ProteinScheme(-10, -1)
	if s.GapOpen != -10 {
		t.Fatal("bad gap open")
	}
}
