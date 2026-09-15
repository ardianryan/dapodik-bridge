//go:build !windows

package tray

import (
	"context"
)

// Run on non-Windows platforms acts as a graceful wait-on-context stub
func Run(ctx context.Context, state *AppState) {
	<-ctx.Done()
	if state.OnQuit != nil {
		state.OnQuit()
	}
}
