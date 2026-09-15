package tray

import "github.com/ardianryan/dapodik-bridge/internal/config"

// AppState holds references for updating tray status dynamically
type AppState struct {
	Config      *config.Config
	TriggerSync func() error
	OnQuit      func()
}
