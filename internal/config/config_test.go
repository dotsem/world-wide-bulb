package config_test

import (
	"os"
	"testing"
	"world-wide-bulb/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Run("loads dev defaults when env vars are unset", func(t *testing.T) {
		os.Clearenv()

		cfg, err := config.Load()
		require.NoError(t, err)

		assert.Equal(t, "8080", cfg.BackendPort)
		assert.False(t, cfg.IsProd)
		assert.Equal(t, "dev_fallback", cfg.IPSalt)
		assert.Equal(t, "bulb.db", cfg.DBPath)
		assert.Empty(t, cfg.AllowedHosts)
		assert.Equal(t, int64(0), cfg.RetentionLimit)
	})

	t.Run("returns error in prod when IP_SALT is missing", func(t *testing.T) {
		os.Clearenv()
		t.Setenv("APP_ENV", "production")

		_, err := config.Load()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "IP_SALT must be set in production")
	})

	t.Run("loads production config successfully when valid", func(t *testing.T) {
		os.Clearenv()
		t.Setenv("APP_ENV", "production")
		t.Setenv("IP_SALT", "super_secret_salt")
		t.Setenv("BACKEND_PORT", "3000")
		t.Setenv("DB_PATH", "/data/prod.db")
		t.Setenv("ALLOWED_HOSTS", "example.com,bulb.example.com")
		t.Setenv("RETENTION_LIMIT", "10000")

		cfg, err := config.Load()
		require.NoError(t, err)

		assert.True(t, cfg.IsProd)
		assert.Equal(t, "3000", cfg.BackendPort)
		assert.Equal(t, "super_secret_salt", cfg.IPSalt)
		assert.Equal(t, "/data/prod.db", cfg.DBPath)
		assert.Equal(t, []string{"example.com", "bulb.example.com"}, cfg.AllowedHosts)
		assert.Equal(t, int64(10000), cfg.RetentionLimit)
	})

	t.Run("handles invalid and negative RETENTION_LIMIT gracefully", func(t *testing.T) {
		os.Clearenv()
		t.Setenv("RETENTION_LIMIT", "invalid_number")

		cfg, err := config.Load()
		require.NoError(t, err)
		assert.Equal(t, int64(0), cfg.RetentionLimit)

		t.Setenv("RETENTION_LIMIT", "-50")
		cfg2, err := config.Load()
		require.NoError(t, err)
		assert.Equal(t, int64(0), cfg2.RetentionLimit)
	})
}
