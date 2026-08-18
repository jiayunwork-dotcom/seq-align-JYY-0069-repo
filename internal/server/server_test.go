package server

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	s := New(":0")
	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestAlignGlobal(t *testing.T) {
	s := New(":0")
	body := AlignRequest{SeqA: "ACGT", SeqB: "ACCT", Mode: "global", Match: 1, Mismatch: -1, GapOpen: -2, GapExtend: -1}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/align", bytes.NewReader(data))
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp AlignResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Score <= 0 {
		t.Fatalf("expected positive score, got %d", resp.Score)
	}
}

func TestAlignLocal(t *testing.T) {
	s := New(":0")
	body := AlignRequest{SeqA: "XXXXACGTXXXX", SeqB: "ACGT", Mode: "local", Match: 2, Mismatch: -1, GapOpen: -3, GapExtend: -1}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/align", bytes.NewReader(data))
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestAlignEmpty(t *testing.T) {
	s := New(":0")
	body := AlignRequest{SeqA: "", SeqB: "ACGT"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/align", bytes.NewReader(data))
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
