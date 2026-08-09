package utils

import (
	"errors"
	"os"
)

// GetEnv returns the value of the environment variable key.
// It returns defaultVal if the environment variable is not set.
func GetEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

// GetEnvOrErr returns the value of the environment variable key.
// It returns an error if the environment variable is not set.
func GetEnvOrErr(key string) (string, error) {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val, nil
	}
	return "", errors.New("missing required environment variable: " + key)
}

// GetEnvOrErrInProd returns the value of the environment variable key.
// It returns an error if the environment variable is not set in production.
// In development, it returns devFallback.
func GetEnvOrErrInProd(key string, devFallback string) (string, error) {
	if !IsProd() {
		return GetEnv(key, devFallback), nil
	}
	return GetEnvOrErr(key)
}
