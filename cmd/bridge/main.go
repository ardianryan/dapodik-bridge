package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardianryan/dapodik-bridge/internal/autostart"
	"github.com/ardianryan/dapodik-bridge/internal/config"
	"github.com/ardianryan/dapodik-bridge/internal/database"
	"github.com/ardianryan/dapodik-bridge/internal/server"
	"github.com/ardianryan/dapodik-bridge/internal/service"
	"github.com/ardianryan/dapodik-bridge/internal/tray"
)

const banner = `
===================================================================
    ____                    _ _    ____       _     _            
   |  _ \  __ _ _ __   ___   | (_) | __ ) _ __(_) __| | __ _  ___ 
   | | | |/ _` + "`" + ` | '_ \ / _ \  | | | |  _ \| '__| |/ _` + "`" + ` |/ _` + "`" + ` |/ _ \
   | |_| | (_| | |_) | (_) | | | | | |_) | |  | | (_| | (_| |  __/
   |____/ \__,_| .__/ \___/  |_|_| |____/|_|  |_|\__,_|\__, |\___|
               |_|                                     |___/      
   Dapodik Read-Only Bridge Daemon (Version %s)
   Strict Read-Only Mode: ENFORCED
===================================================================
`

func runBridgeEngine(ctx context.Context, cfg *config.Config) error {
	log.Printf("[INFO] Initializing Dapodik Bridge on port %d (host: %s)...", cfg.Port, cfg.Host)
	log.Printf("[INFO] Dapodik target DB: %s:%d/%s (user: %s)", cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser)

	dbMgr, err := database.NewDBManager(ctx, cfg)
	if err != nil {
		log.Printf("[WARN] Initial database connection check: %v", err)
	}
	defer dbMgr.Close()

	srv := server.NewServer(cfg, dbMgr)

	errChan := make(chan error, 1)
	go func() {
		log.Printf("[INFO] Dapodik Bridge Daemon is listening on http://%s:%d", cfg.Host, cfg.Port)
		log.Printf("[INFO] Endpoints ready: /api/v1/health, /api/v1/kesejahteraan, /api/v1/pip, /api/v1/rapor, /api/v1/schema/tables")
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Printf("[INFO] Shutting down bridge server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errChan:
		return err
	}
}

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]

		// 1. Version check
		if arg == "version" || arg == "-version" || arg == "--version" {
			fmt.Printf("dapodik-bridge version %s\n", config.AppVersion)
			os.Exit(0)
		}

		// 2. Windows Service / Linux Systemd CLI Command
		if arg == "service" {
			cmd := "status"
			if len(os.Args) > 2 {
				cmd = os.Args[2]
			}
			cfg, err := config.LoadConfig(os.Args[3:])
			if err != nil {
				log.Fatalf("[FATAL] Configuration error: %v", err)
			}
			svc, err := service.New(func(ctx context.Context) error {
				return runBridgeEngine(ctx, cfg)
			})
			if err != nil {
				log.Fatalf("[FATAL] Service initialization error: %v", err)
			}
			if err := service.HandleCommand(cmd, svc); err != nil {
				log.Fatalf("[ERROR] %v", err)
			}
			return
		}

		// 3. Explicit CLI Console Run Mode
		if arg == "run" {
			fmt.Printf(banner, config.AppVersion)
			cfg, err := config.LoadConfig(os.Args[2:])
			if err != nil {
				log.Fatalf("[FATAL] Configuration error: %v", err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			stopChan := make(chan os.Signal, 1)
			signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
			go func() {
				sig := <-stopChan
				log.Printf("[INFO] Received signal %v, shutting down...", sig)
				cancel()
			}()

			if err := runBridgeEngine(ctx, cfg); err != nil {
				log.Fatalf("[FATAL] %v", err)
			}
			return
		}
	}

	// 4. Default Mode: System Tray GUI (Windows Desktop) or Auto-Headless Console
	isHeadless := runtime.GOOS != "windows" && os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == ""

	var appArgs []string
	for _, a := range os.Args[1:] {
		if a != "--minimized" && a != "--tray" && a != "-tray" {
			appArgs = append(appArgs, a)
		}
	}

	cfg, err := config.LoadConfig(appArgs)
	if err != nil {
		log.Fatalf("[FATAL] Configuration error: %v", err)
	}

	if isHeadless {
		// Headless Linux / Docker container fallback
		fmt.Printf(banner, config.AppVersion)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
		go func() {
			sig := <-stopChan
			log.Printf("[INFO] Received signal %v, shutting down...", sig)
			cancel()
		}()

		if err := runBridgeEngine(ctx, cfg); err != nil {
			log.Fatalf("[FATAL] %v", err)
		}
		return
	}

	// Desktop Tray Mode
	// On Windows first-run, enable autostart in HKCU automatically
	autoMgr := autostart.NewManager()
	if !autoMgr.IsEnabled() {
		_ = autoMgr.Enable()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Bridge HTTP server in background
	go func() {
		if err := runBridgeEngine(ctx, cfg); err != nil {
			log.Printf("[ERROR] Bridge engine error: %v", err)
			cancel()
		}
	}()

	// Run System Tray (blocks until user quits or ctx is canceled)
	appState := &tray.AppState{
		Config: cfg,
		TriggerSync: func() error {
			log.Printf("[INFO] Manual sync requested from System Tray")
			return nil
		},
		OnQuit: func() {
			log.Printf("[INFO] System Tray requested shutdown")
			cancel()
		},
	}

	tray.Run(ctx, appState)
}
