package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardianryan/dapodik-bridge/internal/config"
	"github.com/ardianryan/dapodik-bridge/internal/database"
	"github.com/ardianryan/dapodik-bridge/internal/models"
)

func TestRootEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port: 4712,
		Host: "0.0.0.0",
	}
	db, _ := database.NewDBManager(t.Context(), cfg)
	srv := NewServer(cfg, db)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp models.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success true, got false")
	}
}

func TestHealthEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port:   4712,
		DBHost: "127.0.0.1",
		DBPort: 5432,
		DBName: "dapodik_dasmen",
	}
	db, _ := database.NewDBManager(t.Context(), cfg)
	srv := NewServer(cfg, db)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rr := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rr, req)

	// Since local postgres might not be active during unit tests, expect either 200 (if active) or 503 (degraded)
	if rr.Code != http.StatusOK && rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 200 or 503, got %d", rr.Code)
	}

	var resp models.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

func TestAuthMiddleware(t *testing.T) {
	cfg := &config.Config{
		Port:   4712,
		APIKey: "my-secret-token",
	}
	db, _ := database.NewDBManager(t.Context(), cfg)
	srv := NewServer(cfg, db)

	// 1. Without Token -> 401
	reqUnauth := httptest.NewRequest(http.MethodGet, "/api/v1/kesejahteraan", nil)
	rrUnauth := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rrUnauth, reqUnauth)

	if rrUnauth.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 Unauthorized, got %d", rrUnauth.Code)
	}

	// 2. With Wrong Token -> 401
	reqWrong := httptest.NewRequest(http.MethodGet, "/api/v1/kesejahteraan", nil)
	reqWrong.Header.Set("Authorization", "Bearer wrong-key")
	rrWrong := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rrWrong, reqWrong)

	if rrWrong.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 Unauthorized, got %d", rrWrong.Code)
	}

	// 3. Health check should bypass auth -> not 401
	reqHealth := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rrHealth := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rrHealth, reqHealth)

	if rrHealth.Code == http.StatusUnauthorized {
		t.Errorf("health endpoint should not require auth, got %d", rrHealth.Code)
	}
}

func TestCORSOptions(t *testing.T) {
	cfg := &config.Config{
		Port: 4712,
	}
	db, _ := database.NewDBManager(t.Context(), cfg)
	srv := NewServer(cfg, db)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/kesejahteraan", nil)
	rr := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status 204 No Content for OPTIONS, got %d", rr.Code)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected CORS allow origin *, got %s", rr.Header().Get("Access-Control-Allow-Origin"))
	}
}
