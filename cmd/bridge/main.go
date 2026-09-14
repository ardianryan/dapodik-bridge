package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardianryan/dapodik-bridge/internal/config"
	"github.com/ardianryan/dapodik-bridge/internal/database"
	"github.com/ardianryan/dapodik-bridge/internal/server"
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

func main() {
	// Handle quick version argument
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "-version" || os.Args[1] == "--version") {
		fmt.Printf("dapodik-bridge version %s\n", config.AppVersion)
		os.Exit(0)
	}

	fmt.Printf(banner, config.AppVersion)

	cfg, err := config.LoadConfig(os.Args[1:])
	if err != nil {
		log.Fatalf("[FATAL] Configuration error: %v", err)
	}

	log.Printf("[INFO] Initializing Dapodik Bridge on port %d (host: %s)...", cfg.Port, cfg.Host)
	log.Printf("[INFO] Dapodik target DB: %s:%d/%s (user: %s)", cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser)
	if cfg.APIKey != "" {
		log.Printf("[INFO] API Key authentication is ENABLED.")
	} else {
		log.Printf("[WARN] API Key is not set. All incoming requests will be accepted.")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize DB manager
	dbMgr, err := database.NewDBManager(ctx, cfg)
	if err != nil {
		log.Printf("[WARN] Initial database connection check: %v", err)
	}
	defer dbMgr.Close()

	// Initialize HTTP Server
	srv := server.NewServer(cfg, dbMgr)

	// Graceful shutdown channel
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[INFO] Dapodik Bridge Daemon is listening on http://%s:%d", cfg.Host, cfg.Port)
		log.Printf("[INFO] Endpoints ready: /api/v1/health, /api/v1/kesejahteraan, /api/v1/pip, /api/v1/rapor, /api/v1/schema/tables")
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] HTTP server failed to start: %v", err)
		}
	}()

	// Block until OS signal received
	sig := <-stopChan
	log.Printf("[INFO] Received signal %v, shutting down gracefully...", sig)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("[ERROR] Server shutdown error: %v", err)
	} else {
		log.Printf("[INFO] Server stopped gracefully.")
	}
}
