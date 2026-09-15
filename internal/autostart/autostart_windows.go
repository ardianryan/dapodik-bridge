//go:build windows

package autostart

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const (
	runKey  = `Software\Microsoft\Windows\CurrentVersion\Run`
	appName = "DapodikBridge"
)

type windowsManager struct{}

// NewManager creates a new autostart manager for Windows
func NewManager() Manager {
	return &windowsManager{}
}

func (m *windowsManager) IsEnabled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	val, _, err := k.GetStringValue(appName)
	if err != nil {
		return false
	}
	return strings.TrimSpace(val) != ""
}

func (m *windowsManager) Enable() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer k.Close()

	// Add --minimized flag so upon boot it runs silently in system tray
	cmd := fmt.Sprintf(`"%s" --minimized`, exePath)
	return k.SetStringValue(appName, cmd)
}

func (m *windowsManager) Disable() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer k.Close()

	err = k.DeleteValue(appName)
	if err == registry.ErrNotExist {
		return nil
	}
	return err
}
