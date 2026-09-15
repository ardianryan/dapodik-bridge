package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ardianryan/dapodik-bridge/internal/config"
	"github.com/ardianryan/dapodik-bridge/internal/database"
	"github.com/ardianryan/dapodik-bridge/internal/models"
)

// Server represents the HTTP Bridge Daemon
type Server struct {
	cfg       *config.Config
	db        *database.DBManager
	startTime time.Time
	httpSrv   *http.Server
}

// NewServer creates a new Server instance
func NewServer(cfg *config.Config, db *database.DBManager) *Server {
	return &Server{
		cfg:       cfg,
		db:        db,
		startTime: time.Now(),
	}
}

// Routes initializes all HTTP routes and middleware
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	// Root, Dashboard & Health
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/dashboard", s.handleDashboard)
	mux.HandleFunc("/api/v1/health", s.handleHealth)

	// Welfare / Bansos
	mux.HandleFunc("/api/v1/kesejahteraan", s.handleKesejahteraan)
	mux.HandleFunc("/api/v1/pip", s.handlePIP)

	// Grades / Rapor
	mux.HandleFunc("/api/v1/rapor", s.handleRapor)

	// Students & GTK
	mux.HandleFunc("/api/v1/siswa/komprehensif", s.handleSiswaKomprehensif)
	mux.HandleFunc("/api/v1/gtk/lengkap", s.handleGTKLengkap)
	mux.HandleFunc("/api/v1/rombel", s.handleRombel)

	// Schema Inspector
	mux.HandleFunc("/api/v1/schema/tables", s.handleSchemaTables)
	mux.HandleFunc("/api/v1/schema/columns", s.handleSchemaColumns)

	// Cloud Sync Push
	mux.HandleFunc("/api/v1/sync/push", s.handleSyncPush)

	// Wrap middleware stack: Logger -> CORS -> Auth -> Mux
	handler := AuthMiddleware(s.cfg, mux)
	handler = CORSMiddleware(handler)
	handler = LoggingAndMetricsMiddleware(handler)

	return handler
}

// Start runs the HTTP server on configured host and port
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	s.httpSrv = &http.Server{
		Addr:         addr,
		Handler:      s.Routes(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return s.httpSrv.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpSrv != nil {
		return s.httpSrv.Shutdown(ctx)
	}
	return nil
}

// Helper methods for JSON responses

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}, meta *models.MetaInfo) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	resp := models.Response{
		Success:   statusCode >= 200 && statusCode < 300,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now(),
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string, execMs int64) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	resp := models.Response{
		Success: false,
		Error:   message,
		Message: message,
		Meta: &models.MetaInfo{
			ExecutionMs: execMs,
			BridgeVer:   config.AppVersion,
		},
		Timestamp: time.Now(),
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func parsePagination(r *http.Request) (limit, offset, page, perPage int) {
	page = 1
	perPage = 50

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if pp, err := strconv.Atoi(r.URL.Query().Get("per_page")); err == nil && pp > 0 {
		if pp > 500 {
			pp = 500 // Safety cap
		}
		perPage = pp
	}

	offset = (page - 1) * perPage
	limit = perPage
	return
}
