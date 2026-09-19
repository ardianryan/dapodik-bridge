package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ardianryan/dapodik-bridge/internal/config"
	"github.com/ardianryan/dapodik-bridge/internal/models"
)

func (s *Server) handleSchemaTables(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", 0)
		return
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	tables, err := s.db.GetSchemaTables(ctx)
	execMs := time.Since(start).Milliseconds()
	if err != nil {
		log.Printf("[ERROR] handleSchemaTables error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Gagal inspeksi skema tabel", execMs)
		return
	}

	meta := &models.MetaInfo{
		TotalCount:  len(tables),
		Count:       len(tables),
		Database:    s.db.ActiveDB(),
		BridgeVer:   config.AppVersion,
		ExecutionMs: execMs,
	}

	writeJSON(w, http.StatusOK, tables, meta)
}

func (s *Server) handleSchemaColumns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", 0)
		return
	}

	table := r.URL.Query().Get("table")
	if table == "" {
		writeJSONError(w, http.StatusBadRequest, "Parameter query 'table' wajib diisi", 0)
		return
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cols, err := s.db.GetSchemaColumns(ctx, table)
	execMs := time.Since(start).Milliseconds()
	if err != nil {
		log.Printf("[ERROR] handleSchemaColumns error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Gagal inspeksi kolom tabel", execMs)
		return
	}

	meta := &models.MetaInfo{
		TotalCount:  len(cols),
		Count:       len(cols),
		Database:    s.db.ActiveDB(),
		BridgeVer:   config.AppVersion,
		ExecutionMs: execMs,
	}

	writeJSON(w, http.StatusOK, cols, meta)
}
