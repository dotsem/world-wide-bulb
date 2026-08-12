// Package config manages application configuration and environment variable loading.
package config

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds runtime configuration settings for the service.
type Config struct {
	Port         string
	IsProd       bool
	IPSalt       string
	DBPath       string
	AllowedHosts []string
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
	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		IsProd:       os.Getenv("APP_ENV") == "production",
		IPSalt:       os.Getenv("IP_SALT"),
		DBPath:       getEnv("DB_PATH", "bulb.db"),
		AllowedHosts: hosts,
	}
	if cfg.IsProd && cfg.IPSalt == "" {
		return nil, errors.New("IP_SALT must be set in production")
	}
	if cfg.IPSalt == "" {
		cfg.IPSalt = "dev_fallback"
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
