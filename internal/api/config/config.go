package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	Port         string
	IsProd       bool
	IPSalt       string
	DBPath       string
	AllowedHosts []string
}

func Load() (*Config, error) {
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
