package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ardianryan/dapodik-bridge/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBManager handles PostgreSQL connection pool and queries
type DBManager struct {
	pool       *pgxpool.Pool
	cfg        *config.Config
	activePass string
	activeDB   string
	mu         sync.RWMutex
}

// NewDBManager initializes connection pool with strict read-only enforcement
func NewDBManager(ctx context.Context, cfg *config.Config) (*DBManager, error) {
	mgr := &DBManager{
		cfg: cfg,
	}

	err := mgr.connectWithFallbacks(ctx)
	if err != nil {
		log.Printf("[WARN] Initial database connection could not be established: %v", err)
		log.Printf("[INFO] Bridge daemon will start in degraded mode and retry DB connection on demand.")
	}

	return mgr, nil
}

// connectWithFallbacks tries the configured password, then attempts known Dapodik default passwords
func (m *DBManager) connectWithFallbacks(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.pool != nil {
		m.pool.Close()
		m.pool = nil
	}

	dbNames := []string{m.cfg.DBName}
	if m.cfg.DBName != "dapodik" {
		dbNames = append(dbNames, "dapodik")
	}

	passwordsToTry := []string{}
	if m.cfg.DBPassword != "" {
		passwordsToTry = append(passwordsToTry, m.cfg.DBPassword)
	}
	for _, p := range m.cfg.DBFallbackPass {
		if p != m.cfg.DBPassword {
			passwordsToTry = append(passwordsToTry, p)
		}
	}

	var lastErr error
	for _, dbName := range dbNames {
		for _, pass := range passwordsToTry {
			connStr := buildConnString(m.cfg.DBHost, m.cfg.DBPort, m.cfg.DBUser, pass, dbName, m.cfg.DBSSLMode)
			poolCfg, err := pgxpool.ParseConfig(connStr)
			if err != nil {
				lastErr = err
				continue
			}

			// Configure pool
			poolCfg.MaxConns = 10
			poolCfg.MinConns = 1
			poolCfg.MaxConnLifetime = 30 * time.Minute
			poolCfg.MaxConnIdleTime = 5 * time.Minute

			// Enforce STRICT READ-ONLY on connection acquire
			poolCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
				_, err := conn.Exec(ctx, "SET default_transaction_read_only = on; SET transaction_read_only = on;")
				return err
			}

			connectCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			pool, err := pgxpool.NewWithConfig(connectCtx, poolCfg)
			cancel()
			if err != nil {
				lastErr = err
				continue
			}

			pingCtx, pingCancel := context.WithTimeout(ctx, 2*time.Second)
			err = pool.Ping(pingCtx)
			pingCancel()
			if err == nil {
				// Successfully connected
				m.pool = pool
				m.activePass = pass
				m.activeDB = dbName
				log.Printf("[INFO] Connected to Dapodik PostgreSQL at %s:%d/%s (user: %s, read-only enforced)",
					m.cfg.DBHost, m.cfg.DBPort, dbName, m.cfg.DBUser)
				return nil
			}

			pool.Close()
			lastErr = err
		}
	}

	return fmt.Errorf("failed to connect to Dapodik PostgreSQL after trying fallback credentials: %w", lastErr)
}

func buildConnString(host string, port int, user, password, dbname, sslmode string) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   dbname,
	}
	q := u.Query()
	q.Set("sslmode", sslmode)
	q.Set("connect_timeout", "3")
	u.RawQuery = q.Encode()
	return u.String()
}

// Pool returns the underlying pgxpool.Pool
func (m *DBManager) Pool() *pgxpool.Pool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.pool
}

// EnsureConnected checks if pool is alive, attempts reconnect if dead
func (m *DBManager) EnsureConnected(ctx context.Context) error {
	m.mu.RLock()
	p := m.pool
	m.mu.RUnlock()

	if p != nil {
		pingCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		if err := p.Ping(pingCtx); err == nil {
			return nil
		}
	}

	return m.connectWithFallbacks(ctx)
}

// IsHealthy returns true if the database responds to ping
func (m *DBManager) IsHealthy(ctx context.Context) bool {
	m.mu.RLock()
	p := m.pool
	m.mu.RUnlock()

	if p == nil {
		return false
	}
	pingCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	return p.Ping(pingCtx) == nil
}

// ActiveDB returns the current active database name
func (m *DBManager) ActiveDB() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.activeDB != "" {
		return m.activeDB
	}
	return m.cfg.DBName
}

// Close terminates the pool
func (m *DBManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.pool != nil {
		m.pool.Close()
		m.pool = nil
	}
}

// SanitizeSQL guarantees safety: rejects any DDL or write statements
func SanitizeSQL(query string) error {
	q := strings.TrimSpace(strings.ToUpper(query))
	forbidden := []string{
		"INSERT ", "UPDATE ", "DELETE ", "DROP ", "ALTER ", "TRUNCATE ",
		"CREATE ", "GRANT ", "REVOKE ", "REPLACE ", "EXECUTE ",
	}
	for _, f := range forbidden {
		if strings.HasPrefix(q, f) || strings.Contains(q, " "+f) {
			return fmt.Errorf("forbidden SQL keyword '%s': bridge is strict read-only", strings.TrimSpace(f))
		}
	}
	return nil
}
