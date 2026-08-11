package config

import (
	"errors"
	"os"
)

type Config struct {
	Port   string
	IsProd bool
	IPSalt string
	DBPath string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:   getEnv("PORT", "8080"),
		IsProd: os.Getenv("APP_ENV") == "production",
		IPSalt: os.Getenv("IP_SALT"),
		DBPath: getEnv("DB_PATH", "bulb.db"),
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
