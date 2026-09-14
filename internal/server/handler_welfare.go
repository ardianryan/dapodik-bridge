package server

import (
	"context"
	"net/http"
	"time"

	"github.com/ardianryan/dapodik-bridge/internal/config"
	"github.com/ardianryan/dapodik-bridge/internal/models"
)

func (s *Server) handleKesejahteraan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", 0)
		return
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	jenis := r.URL.Query().Get("jenis")
	limit, offset, page, perPage := parsePagination(r)

	results, total, err := s.db.GetKesejahteraan(ctx, jenis, limit, offset)
	execMs := time.Since(start).Milliseconds()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Gagal membaca data kesejahteraan: "+err.Error(), execMs)
		return
	}

	meta := &models.MetaInfo{
		TotalCount:  total,
		Count:       len(results),
		Page:        page,
		PerPage:     perPage,
		Database:    s.db.ActiveDB(),
		BridgeVer:   config.AppVersion,
		ExecutionMs: execMs,
	}

	writeJSON(w, http.StatusOK, results, meta)
}

func (s *Server) handlePIP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", 0)
		return
	}

	// Alias specifically filtering for PIP / KIP
	start := time.Now()
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	limit, offset, page, perPage := parsePagination(r)
	results, total, err := s.db.GetKesejahteraan(ctx, "pip", limit, offset)
	execMs := time.Since(start).Milliseconds()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Gagal membaca data PIP: "+err.Error(), execMs)
		return
	}

	meta := &models.MetaInfo{
		TotalCount:  total,
		Count:       len(results),
		Page:        page,
		PerPage:     perPage,
		Database:    s.db.ActiveDB(),
		BridgeVer:   config.AppVersion,
		ExecutionMs: execMs,
	}

	writeJSON(w, http.StatusOK, results, meta)
}
