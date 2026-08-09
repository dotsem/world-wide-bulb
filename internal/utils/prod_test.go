package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsProd(t *testing.T) {
	originalEnv := os.Getenv("APP_ENV")
	defer func() { _ = os.Setenv("APP_ENV", originalEnv) }()

	_ = os.Setenv("APP_ENV", "development")
	assert.False(t, IsProd())

	_ = os.Setenv("APP_ENV", "production")
	assert.True(t, IsProd())
}
