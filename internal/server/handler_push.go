package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ardianryan/dapodik-bridge/internal/config"
	"github.com/ardianryan/dapodik-bridge/internal/models"
)

type PushRequest struct {
	TargetURL   string `json:"target_url,omitempty"`
	SecretToken string `json:"secret_token,omitempty"`
	Type        string `json:"type,omitempty"` // "welfare", "rapor", "siswa", "gtk", "all"
	SemesterID  string `json:"semester_id,omitempty"`
}

type PushResult struct {
	TargetURL    string `json:"target_url"`
	Type         string `json:"type"`
	RecordsSent  int    `json:"records_sent"`
	TargetStatus int    `json:"target_http_status"`
	ResponseBody string `json:"target_response_body,omitempty"`
	Success      bool   `json:"success"`
}

func (s *Server) handleSyncPush(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed. Use POST", 0)
		return
	}

	start := time.Now()
	var req PushRequest

	// Support both JSON body and query parameters
	if r.Header.Get("Content-Type") == "application/json" {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	if req.Type == "" {
		req.Type = r.URL.Query().Get("type")
	}
	if req.Type == "" {
		req.Type = "welfare" // Default push type
	}

	targetURL := req.TargetURL
	if targetURL == "" {
		targetURL = s.cfg.PushTargetURL
	}
	if targetURL == "" {
		writeJSONError(w, http.StatusBadRequest, "Target URL belum dikonfigurasi. Isi 'target_url' pada body JSON atau atur PUSH_TARGET_URL pada environment.", 0)
		return
	}

	secretToken := req.SecretToken
	if secretToken == "" {
		secretToken = s.cfg.PushSecretToken
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	var dataToPush interface{}
	recordCount := 0

	switch req.Type {
	case "welfare", "kesejahteraan", "pip":
		items, count, err := s.db.GetKesejahteraan(ctx, "", 10000, 0)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Gagal mengambil data kesejahteraan: "+err.Error(), time.Since(start).Milliseconds())
			return
		}
		dataToPush = items
		recordCount = count

	case "rapor":
		items, count, err := s.db.GetRaporGrades(ctx, req.SemesterID, "", "", 10000, 0)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Gagal mengambil data rapor: "+err.Error(), time.Since(start).Milliseconds())
			return
		}
		dataToPush = items
		recordCount = count

	case "siswa":
		items, count, err := s.db.GetSiswaKomprehensif(ctx, "", "", 10000, 0)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Gagal mengambil data siswa: "+err.Error(), time.Since(start).Milliseconds())
			return
		}
		dataToPush = items
		recordCount = count

	case "gtk":
		items, count, err := s.db.GetGTKLengkap(ctx, "", 10000, 0)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Gagal mengambil data GTK: "+err.Error(), time.Since(start).Milliseconds())
			return
		}
		dataToPush = items
		recordCount = count

	default:
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Push type '%s' tidak valid. Pilihan: welfare, rapor, siswa, gtk", req.Type), 0)
		return
	}

	hostname, _ := os.Hostname()
	payload := models.PushPayload{
		SourceHost:  hostname,
		SchoolNPSN:  s.cfg.SchoolNPSN,
		PushType:    req.Type,
		RecordCount: recordCount,
		PushedAt:    time.Now(),
		Data:        dataToPush,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Gagal membuat payload push: "+err.Error(), time.Since(start).Milliseconds())
		return
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Gagal membuat HTTP request push: "+err.Error(), time.Since(start).Milliseconds())
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "Dapodik-Bridge/"+config.AppVersion)
	if secretToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+secretToken)
		httpReq.Header.Set("X-Bridge-Token", secretToken)
	}
	if s.cfg.SchoolNPSN != "" {
		httpReq.Header.Set("X-School-NPSN", s.cfg.SchoolNPSN)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, fmt.Sprintf("Gagal menghubungi target push (%s): %v", targetURL, err), time.Since(start).Milliseconds())
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))

	result := PushResult{
		TargetURL:    targetURL,
		Type:         req.Type,
		RecordsSent:  recordCount,
		TargetStatus: resp.StatusCode,
		ResponseBody: string(respBody),
		Success:      resp.StatusCode >= 200 && resp.StatusCode < 300,
	}

	meta := &models.MetaInfo{
		Count:       recordCount,
		BridgeVer:   config.AppVersion,
		ExecutionMs: time.Since(start).Milliseconds(),
	}

	statusCode := http.StatusOK
	if !result.Success {
		statusCode = http.StatusBadGateway
	}

	writeJSON(w, statusCode, result, meta)
}
