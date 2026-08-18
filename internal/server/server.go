// Package server provides an HTTP API for the sequence alignment tool,
// enabling the frontend to submit sequences and receive alignment results.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"seq-align/internal/align"
	"seq-align/internal/scoring"
)

// Server is the HTTP API server.
type Server struct {
	mux  *http.ServeMux
	addr string
}

// New creates a server.
func New(addr string) *Server {
	s := &Server{mux: http.NewServeMux(), addr: addr}
	s.mux.HandleFunc("/api/align", s.handleAlign)
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.Handle("/", http.FileServer(http.Dir("frontend")))
	return s
}

// Handler returns the HTTP handler.
func (s *Server) Handler() http.Handler { return s.mux }

// ListenAndServe starts the server.
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s.mux)
}

// AlignRequest is the input.
type AlignRequest struct {
	SeqA      string `json:"seq_a"`
	SeqB      string `json:"seq_b"`
	Mode      string `json:"mode"` // "global" or "local"
	Match     int    `json:"match"`
	Mismatch  int    `json:"mismatch"`
	GapOpen   int    `json:"gap_open"`
	GapExtend int    `json:"gap_extend"`
}

// AlignResponse is the output.
type AlignResponse struct {
	AlignedA string `json:"aligned_a"`
	AlignedB string `json:"aligned_b"`
	Score    int    `json:"score"`
	Identity float64 `json:"identity"`
	Gaps     int    `json:"gaps"`
	Length   int    `json:"length"`
}

func (s *Server) handleAlign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req AlignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.SeqA == "" || req.SeqB == "" {
		http.Error(w, "sequences required", http.StatusBadRequest)
		return
	}

	match := req.Match
	if match == 0 {
		match = 1
	}
	mismatch := req.Mismatch
	if mismatch == 0 {
		mismatch = -1
	}
	gapOpen := req.GapOpen
	if gapOpen == 0 {
		gapOpen = -2
	}
	gapExtend := req.GapExtend
	if gapExtend == 0 {
		gapExtend = -1
	}

	scheme := scoring.Scheme{
		Match:     match,
		Mismatch:  mismatch,
		GapOpen:   gapOpen,
		GapExtend: gapExtend,
	}

	var result align.Result
	var err error
	if req.Mode == "local" {
		result, err = align.Local(req.SeqA, req.SeqB, scheme)
	} else {
		result, err = align.Global(req.SeqA, req.SeqB, scheme)
	}
	if err != nil {
		http.Error(w, "alignment error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := AlignResponse{
		AlignedA: result.A,
		AlignedB: result.B,
		Score:    result.Score,
		Identity: result.Identity,
		Gaps:     countGapsStr(result.A, result.B),
		Length:   result.AlignedLen,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok"}`)
}

func countGapsStr(a, b string) int {
	gaps := 0
	for _, c := range a {
		if c == '-' {
			gaps++
		}
	}
	for _, c := range b {
		if c == '-' {
			gaps++
		}
	}
	return gaps
}
