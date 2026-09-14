package server

import (
	"context"
	"net/http"
	"time"

	"github.com/ardianryan/dapodik-bridge/internal/config"
	"github.com/ardianryan/dapodik-bridge/internal/models"
)

func (s *Server) handleRapor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", 0)
		return
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	semesterID := r.URL.Query().Get("semester_id")
	rombelID := r.URL.Query().Get("rombel_id")
	nisn := r.URL.Query().Get("nisn")
	limit, offset, page, perPage := parsePagination(r)

	results, total, err := s.db.GetRaporGrades(ctx, semesterID, rombelID, nisn, limit, offset)
	execMs := time.Since(start).Milliseconds()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Gagal membaca nilai rapor: "+err.Error(), execMs)
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
