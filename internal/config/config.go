// Package config manages application configuration and environment variable loading.
package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds runtime configuration settings for the service.
type Config struct {
	BackendPort    string
	FrontendPort   string
	IsProd         bool
	IPSalt         string
	DBPath         string
	AllowedHosts   []string
	RetentionLimit int64
}

// Load parses environment variables and returns a validated Config.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			slog.Warn("failed to load .env file", slog.Any("err", err))
		} else {
			slog.Debug("no .env file found, relying on process environment")
		}
	}
	var hosts []string
	if rawHosts := os.Getenv("ALLOWED_HOSTS"); rawHosts != "" {
		hosts = strings.Split(rawHosts, ",")
	}
	var retentionLimit int64
	if rawRetention := os.Getenv("RETENTION_LIMIT"); rawRetention != "" {
		if parsed, err := strconv.ParseInt(rawRetention, 10, 64); err == nil && parsed > 0 {
			retentionLimit = parsed
		}
	}
	cfg := &Config{
		BackendPort:    getEnv("8080", "BACKEND_PORT", "PORT"),
		FrontendPort:   getEnv("5001", "FRONTEND_PORT"),
		IsProd:         os.Getenv("APP_ENV") == "production",
		IPSalt:         os.Getenv("IP_SALT"),
		DBPath:         getEnv("bulb.db", "DB_PATH"),
		AllowedHosts:   hosts,
		RetentionLimit: retentionLimit,
	}
	if cfg.IsProd && cfg.IPSalt == "" {
		return nil, errors.New("IP_SALT must be set in production")
	}
	if cfg.IPSalt == "" {
		cfg.IPSalt = "dev_fallback"
	}

	return cfg, nil
}

func getEnv(fallback string, keys ...string) string {
	for _, key := range keys {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	slog.Debug("using default value for keys", slog.Any("keys", keys), slog.String("default", fallback))
	return fallback
}
