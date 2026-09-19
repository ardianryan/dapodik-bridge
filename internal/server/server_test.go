package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestDashboardEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port: 4712,
		Host: "0.0.0.0",
	}
	db, _ := database.NewDBManager(t.Context(), cfg)
	srv := NewServer(cfg, db)

	// Test /dashboard directly
	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	rr := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if rr.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fatalf("expected text/html content type, got %s", rr.Header().Get("Content-Type"))
	}

	// Test / with Accept: text/html (browser request)
	reqBrowser := httptest.NewRequest(http.MethodGet, "/", nil)
	reqBrowser.Header.Set("Accept", "text/html,application/xhtml+xml")
	rrBrowser := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rrBrowser, reqBrowser)

	if rrBrowser.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rrBrowser.Code)
	}
	if rrBrowser.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fatalf("expected text/html content type, got %s", rrBrowser.Header().Get("Content-Type"))
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

	// 1. Allowed origin (localhost)
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/kesejahteraan", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rr := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status 204 No Content for OPTIONS with localhost, got %d", rr.Code)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Errorf("expected CORS allow origin http://localhost:3000, got %s", rr.Header().Get("Access-Control-Allow-Origin"))
	}

	// 2. Disallowed origin (unauthorized domain)
	reqEvil := httptest.NewRequest(http.MethodOptions, "/api/v1/kesejahteraan", nil)
	reqEvil.Header.Set("Origin", "https://evil-attacker.com")
	rrEvil := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rrEvil, reqEvil)

	if rrEvil.Code != http.StatusForbidden {
		t.Errorf("expected status 403 Forbidden for unauthorized origin OPTIONS, got %d", rrEvil.Code)
	}
	if rrEvil.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("expected no CORS allow origin header for evil origin, got %s", rrEvil.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestSyncPushSecurity(t *testing.T) {
	// 1. Without server PUSH_TARGET_URL configured -> should return 400 Bad Request
	cfg := &config.Config{
		Port:          4712,
		PushTargetURL: "", // Not configured
	}
	db, _ := database.NewDBManager(t.Context(), cfg)
	srv := NewServer(cfg, db)

	body := strings.NewReader(`{"target_url": "https://evil.attacker.com/steal-data", "type": "welfare"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sync/push", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 when server PushTargetURL is empty, got %d", rr.Code)
	}
}

