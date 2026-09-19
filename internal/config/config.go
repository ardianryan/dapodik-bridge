package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration options for Dapodik Bridge
type Config struct {
	Port            int      `json:"port"`
	Host            string   `json:"host"`
	DBHost          string   `json:"db_host"`
	DBPort          int      `json:"db_port"`
	DBUser          string   `json:"db_user"`
	DBPassword      string   `json:"db_password"`
	DBName          string   `json:"db_name"`
	DBSSLMode       string   `json:"db_sslmode"`
	DBFallbackPass  []string `json:"-"`
	APIKey          string   `json:"-"`
	PushTargetURL   string   `json:"push_target_url"`
	PushSecretToken string   `json:"-"`
	SchoolNPSN      string   `json:"school_npsn"`
	CORSAllowedOrigins string   `json:"cors_allowed_origins,omitempty"`
	Version            string   `json:"version"`
}

const (
	DefaultPort    = 4712
	DefaultHost    = "127.0.0.1"
	DefaultDBHost  = "127.0.0.1"
	DefaultDBPort  = 5432
	DefaultDBUser  = "postgres"
	DefaultDBName  = "dapodik_dasmen"
	DefaultSSLMode = "disable"
	AppVersion     = "1.1.0"
)

// LoadConfig loads configuration from flags, environment variables, and .env file
func LoadConfig(args []string) (*Config, error) {
	// 1. Try loading .env file if it exists
	_ = godotenv.Load()

	fs := flag.NewFlagSet("dapodik-bridge", flag.ContinueOnError)

	var (
		portFlag        = fs.Int("port", getEnvInt("BRIDGE_PORT", DefaultPort), "Bridge HTTP port (default 4712)")
		hostFlag        = fs.String("host", getEnvString("BRIDGE_HOST", DefaultHost), "Bridge HTTP listen host (default 127.0.0.1)")
		corsOriginsFlag = fs.String("cors-origins", getEnvString("BRIDGE_CORS_ORIGINS", ""), "Allowed CORS origins, comma-separated")
		dbHostFlag      = fs.String("db-host", getEnvString("DB_HOST", DefaultDBHost), "Dapodik PostgreSQL host")
		dbPortFlag      = fs.Int("db-port", getEnvInt("DB_PORT", DefaultDBPort), "Dapodik PostgreSQL port")
		dbUserFlag      = fs.String("db-user", getEnvString("DB_USER", DefaultDBUser), "Dapodik PostgreSQL user")
		dbPassFlag      = fs.String("db-pass", getEnvString("DB_PASSWORD", ""), "Dapodik PostgreSQL password")
		dbNameFlag      = fs.String("db-name", getEnvString("DB_NAME", DefaultDBName), "Dapodik PostgreSQL db name")
		apiKeyFlag      = fs.String("api-key", getEnvString("BRIDGE_API_KEY", ""), "Bearer API key required for incoming requests (optional on localhost, required on 0.0.0.0)")
		pushURLFlag     = fs.String("push-url", getEnvString("PUSH_TARGET_URL", ""), "Target school cloud webhook URL")
		pushSecretFlag  = fs.String("push-secret", getEnvString("PUSH_SECRET_TOKEN", ""), "Secret token sent to push target")
		npsnFlag        = fs.String("npsn", getEnvString("SCHOOL_NPSN", ""), "School NPSN identifier")
		envFileFlag     = fs.String("env-file", "", "Custom path to .env file")
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if *envFileFlag != "" {
		_ = godotenv.Load(*envFileFlag)
	}

	cfg := &Config{
		Port:               *portFlag,
		Host:               *hostFlag,
		CORSAllowedOrigins: *corsOriginsFlag,
		DBHost:             *dbHostFlag,
		DBPort:             *dbPortFlag,
		DBUser:             *dbUserFlag,
		DBPassword:         *dbPassFlag,
		DBName:             *dbNameFlag,
		DBSSLMode:          getEnvString("DB_SSLMODE", DefaultSSLMode),
		DBFallbackPass:     []string{"", "password", "dapodik", "dapodik123", "123456", "admin"},
		APIKey:             *apiKeyFlag,
		PushTargetURL:      *pushURLFlag,
		PushSecretToken:    *pushSecretFlag,
		SchoolNPSN:         *npsnFlag,
		Version:            AppVersion,
	}

	if (cfg.Host == "0.0.0.0" || cfg.Host == "") && cfg.APIKey == "" {
		return nil, fmt.Errorf("BRIDGE_API_KEY wajib diisi jika host diatur ke 0.0.0.0 (semua antarmuka jaringan) untuk keamanan")
	}

	return cfg, nil
}

func getEnvString(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
