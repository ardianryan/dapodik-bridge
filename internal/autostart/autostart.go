package autostart

// Manager defines the interface for managing user-level startup applications
type Manager interface {
	IsEnabled() bool
	Enable() error
	Disable() error
}
