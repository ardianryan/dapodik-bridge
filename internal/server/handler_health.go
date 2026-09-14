package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ardianryan/dapodik-bridge/internal/config"
	"github.com/ardianryan/dapodik-bridge/internal/models"
)

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"service":     "dapodik-bridge",
		"description": "High-Performance Read-Only Local Bridge Daemon for Dapodik Desktop PostgreSQL",
		"version":     config.AppVersion,
		"status":      "running",
		"endpoints": []string{
			"GET /api/v1/health",
			"GET /api/v1/kesejahteraan?jenis={pip|kip|pkh|kks}&page=1&per_page=50",
			"GET /api/v1/pip?page=1&per_page=50",
			"GET /api/v1/rapor?semester_id={id}&rombel_id={id}&nisn={nisn}",
			"GET /api/v1/siswa/komprehensif?q={nama|nisn|nik}&page=1&per_page=50",
			"GET /api/v1/gtk/lengkap?q={nama|nuptk}&page=1&per_page=50",
			"GET /api/v1/rombel?semester_id={id}",
			"GET /api/v1/schema/tables",
			"GET /api/v1/schema/columns?table={table_name}",
			"POST /api/v1/sync/push?type={welfare|rapor|siswa|all}",
		},
		"docs": "https://github.com/ardianryan/dapodik-bridge",
	}

	writeJSON(w, http.StatusOK, data, nil)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	dbUp := s.db.IsHealthy(ctx)

	status := "ok"
	statusCode := http.StatusOK
	if !dbUp {
		status = "degraded"
		// If DB connection is down, return 503 or 200 with degraded status
		statusCode = http.StatusServiceUnavailable
	}

	uptime := time.Since(s.startTime).Truncate(time.Second).String()

	health := models.HealthStatus{
		Status:       status,
		BridgeUp:     true,
		DatabaseUp:   dbUp,
		DatabaseHost: fmt.Sprintf("%s:%d", s.cfg.DBHost, s.cfg.DBPort),
		DatabaseName: s.db.ActiveDB(),
		DatabaseUser: s.cfg.DBUser,
		ReadOnly:     true,
		Uptime:       uptime,
		Version:      config.AppVersion,
		CurrentTime:  time.Now(),
	}

	meta := &models.MetaInfo{
		BridgeVer:   config.AppVersion,
		Database:    s.db.ActiveDB(),
		ExecutionMs: time.Since(start).Milliseconds(),
	}

	writeJSON(w, statusCode, health, meta)
}
