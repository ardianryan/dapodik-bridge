//go:build !windows

package autostart

type noopManager struct{}

// NewManager creates a no-op autostart manager for non-Windows platforms
func NewManager() Manager {
	return &noopManager{}
}

func (m *noopManager) IsEnabled() bool {
	return false
}

func (m *noopManager) Enable() error {
	return nil
}

func (m *noopManager) Disable() error {
	return nil
}
