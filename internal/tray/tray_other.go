//go:build !windows && !(darwin && cgo)

package tray

import (
	"context"
)

// Run on headless platforms acts as a graceful wait-on-context stub
func Run(ctx context.Context, state *AppState) {
	<-ctx.Done()
	if state.OnQuit != nil {
		state.OnQuit()
	}
}
