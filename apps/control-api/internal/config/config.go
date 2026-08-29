package config

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	ListenAddr          string
	DatabaseURL         string
	RedisURL            string
	LogLevel            slog.Level
	AdminInitialEmail   string
	AdminInitialPassword string
}

// Load reads configuration from environment variables.
// Required values fail fast; no defaults are substituted for secrets.
func Load() (*Config, error) {
	logLevel := slog.LevelInfo
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}

	cfg := &Config{
		ListenAddr:  getenv("LISTEN_ADDR", ":8443"),
		LogLevel:    logLevel,
		AdminInitialEmail:   os.Getenv("ADMIN_INITIAL_EMAIL"),
		AdminInitialPassword: os.Getenv("ADMIN_INITIAL_PASSWORD"),
	}

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	cfg.RedisURL = getenv("REDIS_URL", "")

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.AdminInitialEmail == "" || cfg.AdminInitialPassword == "" {
		return nil, fmt.Errorf("ADMIN_INITIAL_EMAIL and ADMIN_INITIAL_PASSWORD are required (no default credentials)")
	}
	return cfg, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
