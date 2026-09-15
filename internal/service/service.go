package service

import (
	"context"
	"fmt"
	"log"

	"github.com/kardianos/service"
)

// Program implements service.Interface for background daemon execution
type Program struct {
	runFunc func(ctx context.Context) error
	cancel  context.CancelFunc
}

func (p *Program) Start(s service.Service) error {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	go func() {
		if err := p.runFunc(ctx); err != nil {
			log.Printf("[ERROR] Service background runner failed: %v", err)
		}
	}()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	if p.cancel != nil {
		p.cancel()
	}
	return nil
}

// Config returns standard service definition
func Config() *service.Config {
	return &service.Config{
		Name:        "DapodikBridge",
		DisplayName: "Dapodik Read-Only Bridge Service",
		Description: "Background API daemon providing secure read-only access to local Dapodik PostgreSQL database on port 4712.",
		Arguments:   []string{"run"},
	}
}

// New creates a new service instance
func New(runFunc func(ctx context.Context) error) (service.Service, error) {
	prg := &Program{runFunc: runFunc}
	return service.New(prg, Config())
}

// HandleCommand processes CLI service commands: install, uninstall, start, stop, restart, status
func HandleCommand(cmd string, s service.Service) error {
	switch cmd {
	case "install":
		err := s.Install()
		if err != nil {
			return fmt.Errorf("failed to install service: %w", err)
		}
		fmt.Println("[SUCCESS] Service 'DapodikBridge' successfully installed.")
		fmt.Println("To start the service now, run: dapodik-bridge service start")
		return nil

	case "uninstall":
		err := s.Uninstall()
		if err != nil {
			return fmt.Errorf("failed to uninstall service: %w", err)
		}
		fmt.Println("[SUCCESS] Service 'DapodikBridge' successfully uninstalled.")
		return nil

	case "start":
		err := s.Start()
		if err != nil {
			return fmt.Errorf("failed to start service: %w", err)
		}
		fmt.Println("[SUCCESS] Service 'DapodikBridge' started.")
		return nil

	case "stop":
		err := s.Stop()
		if err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}
		fmt.Println("[SUCCESS] Service 'DapodikBridge' stopped.")
		return nil

	case "restart":
		err := s.Restart()
		if err != nil {
			return fmt.Errorf("failed to restart service: %w", err)
		}
		fmt.Println("[SUCCESS] Service 'DapodikBridge' restarted.")
		return nil

	case "status":
		status, err := s.Status()
		if err != nil {
			return fmt.Errorf("failed to check service status: %w", err)
		}
		switch status {
		case service.StatusRunning:
			fmt.Println("Service Status: RUNNING")
		case service.StatusStopped:
			fmt.Println("Service Status: STOPPED")
		default:
			fmt.Println("Service Status: UNKNOWN")
		}
		return nil

	default:
		return fmt.Errorf("unknown service action: %s (available: install, uninstall, start, stop, restart, status)", cmd)
	}
}
