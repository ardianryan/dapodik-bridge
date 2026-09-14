package server

import (
	"crypto/subtle"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ardianryan/dapodik-bridge/internal/config"
)

// ResponseWriterWrapper intercepts status code for access logging
type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *ResponseWriterWrapper) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingAndMetricsMiddleware logs incoming HTTP requests with reverse proxy IP detection
func LoggingAndMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		clientIP := r.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = r.Header.Get("X-Real-IP")
		}
		if clientIP == "" {
			clientIP = r.RemoteAddr
		} else {
			// First IP in list is the original client
			if idx := strings.Index(clientIP, ","); idx != -1 {
				clientIP = strings.TrimSpace(clientIP[:idx])
			}
		}

		wrapper := &ResponseWriterWrapper{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)
		log.Printf("[HTTP] %s %s %d %s (client: %s)",
			r.Method, r.URL.Path, wrapper.StatusCode, duration, clientIP)
	})
}

// CORSMiddleware enables cross-origin requests for school dashboards
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, X-Requested-With")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware validates Bearer API key if configured
func AuthMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Health check is always public so reverse proxy / load balancer health probes work
		if r.URL.Path == "/api/v1/health" || r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		// If no API key is set, authentication is disabled
		if cfg.APIKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Check Authorization header or query param
		authHeader := r.Header.Get("Authorization")
		token := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			token = r.URL.Query().Get("api_key")
		}

		if token == "" || subtle.ConstantTimeCompare([]byte(token), []byte(cfg.APIKey)) != 1 {
			writeJSONError(w, http.StatusUnauthorized, "Unauthorized: missing or invalid API key", 0)
			return
		}

		next.ServeHTTP(w, r)
	})
}
