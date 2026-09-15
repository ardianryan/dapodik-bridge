package service

import (
	"context"
	"testing"
)

func TestServiceConfig(t *testing.T) {
	cfg := Config()
	if cfg.Name != "DapodikBridge" {
		t.Fatalf("expected service name DapodikBridge, got %s", cfg.Name)
	}
	if len(cfg.Arguments) == 0 || cfg.Arguments[0] != "run" {
		t.Fatalf("expected service argument 'run', got %v", cfg.Arguments)
	}
}

func TestNewService(t *testing.T) {
	dummyRun := func(ctx context.Context) error {
		return nil
	}
	svc, err := New(dummyRun)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	if svc == nil {
		t.Fatal("expected non-nil service instance")
	}
}
