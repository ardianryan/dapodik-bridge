package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ardianryan/dapodik-bridge/internal/config"
	"github.com/ardianryan/dapodik-bridge/internal/models"
)

func (s *Server) handleRombel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", 0)
		return
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	semesterID := r.URL.Query().Get("semester_id")
	limit, offset, page, perPage := parsePagination(r)

	results, total, err := s.db.GetRombel(ctx, semesterID, limit, offset)
	execMs := time.Since(start).Milliseconds()
	if err != nil {
		log.Printf("[ERROR] handleRombel error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Gagal membaca data rombel", execMs)
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
