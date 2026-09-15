//go:build windows || (darwin && cgo)

package tray

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/ardianryan/dapodik-bridge/internal/autostart"
	"github.com/ardianryan/dapodik-bridge/internal/config"
	"github.com/energye/systray"
)

// Run starts the system tray event loop. Blocks the calling thread.
func Run(ctx context.Context, state *AppState) {
	onReady := func() {
		systray.SetIcon(DefaultIconBytes())
		systray.SetTitle("Dapodik Bridge")
		systray.SetTooltip(fmt.Sprintf("Dapodik Bridge v%s - Port %d", config.AppVersion, state.Config.Port))

		// Menu Item 1: Status Header (Disabled/Display only)
		mStatus := systray.AddMenuItem("● Dapodik Bridge: Active", "Status koneksi server")
		mStatus.Disable()

		// Menu Item 2: Target DB & Port Info
		mInfo := systray.AddMenuItem(fmt.Sprintf("Target: %s:%d (Port %d)", state.Config.DBHost, state.Config.DBPort, state.Config.Port), "Informasi port dan database")
		mInfo.Disable()

		systray.AddSeparator()

		// Menu Item 3: Open Web Dashboard
		dashboardURL := fmt.Sprintf("http://localhost:%d/dashboard", state.Config.Port)
		mDashboard := systray.AddMenuItem("Buka Dashboard (Browser)", "Buka status server di web browser")
		mDashboard.Click(func() {
			if err := openBrowser(dashboardURL); err != nil {
				log.Printf("[WARN] Failed to open browser: %v", err)
			}
		})

		// Menu Item 4: Trigger Sync
		mSync := systray.AddMenuItem("Sinkronisasi Sekarang", "Picu pengiriman auto-push webhook sekarang")
		mSync.Click(func() {
			if state.TriggerSync != nil {
				mSync.SetTitle("Sedang Sinkronisasi...")
				mSync.Disable()
				go func() {
					err := state.TriggerSync()
					if err != nil {
						log.Printf("[WARN] Manual sync trigger returned: %v", err)
					}
					mSync.Enable()
					mSync.SetTitle("Sinkronisasi Sekarang")
				}()
			}
		})

		// Menu Item 5: Autostart toggle (Windows only)
		if runtime.GOOS == "windows" {
			systray.AddSeparator()
			autoMgr := autostart.NewManager()
			mAutostart := systray.AddMenuItemCheckbox("Jalankan saat Startup Windows", "Aktifkan startup otomatis di background", autoMgr.IsEnabled())
			mAutostart.Click(func() {
				if mAutostart.Checked() {
					if err := autoMgr.Disable(); err == nil {
						mAutostart.Uncheck()
						log.Printf("[INFO] Autostart disabled")
					} else {
						log.Printf("[WARN] Failed to disable autostart: %v", err)
					}
				} else {
					if err := autoMgr.Enable(); err == nil {
						mAutostart.Check()
						log.Printf("[INFO] Autostart enabled")
					} else {
						log.Printf("[WARN] Failed to enable autostart: %v", err)
					}
				}
			})
		}

		systray.AddSeparator()

		// Menu Item 6: Quit
		mQuit := systray.AddMenuItem("Keluar dari Dapodik Bridge", "Hentikan bridge daemon")
		mQuit.Click(func() {
			log.Printf("[INFO] Tray exit requested by user")
			if state.OnQuit != nil {
				state.OnQuit()
			}
			systray.Quit()
		})

		// Goroutine to handle external context cancellation
		go func() {
			<-ctx.Done()
			systray.Quit()
		}()
	}

	onExit := func() {
		if state.OnQuit != nil {
			state.OnQuit()
		}
	}

	systray.Run(onReady, onExit)
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
