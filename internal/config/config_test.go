package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	// Clean env
	os.Unsetenv("BRIDGE_PORT")
	os.Unsetenv("BRIDGE_HOST")
	os.Unsetenv("DB_HOST")

	cfg, err := LoadConfig([]string{})
	if err != nil {
		t.Fatalf("unexpected error loading default config: %v", err)
	}

	if cfg.Port != DefaultPort {
		t.Errorf("expected default port %d, got %d", DefaultPort, cfg.Port)
	}
	if cfg.Port != 4712 {
		t.Errorf("expected port 4712 as requested, got %d", cfg.Port)
	}
	if cfg.DBHost != "127.0.0.1" {
		t.Errorf("expected default DBHost 127.0.0.1, got %s", cfg.DBHost)
	}
	if cfg.DBPort != 5432 {
		t.Errorf("expected default DBPort 5432, got %d", cfg.DBPort)
	}
}

func TestConfigFlags(t *testing.T) {
	args := []string{
		"-port=8080",
		"-host=127.0.0.1",
		"-db-host=localhost",
		"-db-port=5433",
		"-db-user=dapodik_user",
		"-db-name=dapodik_custom",
		"-api-key=secret123",
	}

	cfg, err := LoadConfig(args)
	if err != nil {
		t.Fatalf("unexpected error parsing flags: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
	if cfg.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", cfg.Host)
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("expected db-host localhost, got %s", cfg.DBHost)
	}
	if cfg.DBPort != 5433 {
		t.Errorf("expected db-port 5433, got %d", cfg.DBPort)
	}
	if cfg.DBUser != "dapodik_user" {
		t.Errorf("expected db-user dapodik_user, got %s", cfg.DBUser)
	}
	if cfg.DBName != "dapodik_custom" {
		t.Errorf("expected db-name dapodik_custom, got %s", cfg.DBName)
	}
	if cfg.APIKey != "secret123" {
		t.Errorf("expected api-key secret123, got %s", cfg.APIKey)
	}
}
