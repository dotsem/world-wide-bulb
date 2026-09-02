package main

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("returns error on invalid port", func(t *testing.T) {
		t.Setenv("DB_PATH", ":memory:")
		t.Setenv("BACKEND_PORT", "-1")

		ctx := context.Background()
		err := run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to start server")
	})

	t.Run("shuts down cleanly on context cancellation", func(t *testing.T) {
		t.Setenv("DB_PATH", ":memory:")
		t.Setenv("BACKEND_PORT", "0")

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(200 * time.Millisecond)
			cancel()
		}()

		err := run(ctx)
		assert.NoError(t, err)
	})

	t.Run("creates nested db directory and runs with retention limit enabled", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbFile := filepath.Join(tmpDir, "nested", "sub", "test.db")
		t.Setenv("DB_PATH", dbFile)
		t.Setenv("BACKEND_PORT", "0")
		t.Setenv("RETENTION_LIMIT", "500")

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(250 * time.Millisecond)
			cancel()
		}()

		err := run(ctx)
		assert.NoError(t, err)
		assert.FileExists(t, dbFile)
	})

	t.Run("returns error when config loading fails", func(t *testing.T) {
		t.Setenv("APP_ENV", "production")
		t.Setenv("IP_SALT", "")

		ctx := context.Background()
		err := run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to load config")
	})
}

func TestRunRoot(t *testing.T) {
	t.Run("returns error when run fails", func(t *testing.T) {
		t.Setenv("APP_ENV", "production")
		t.Setenv("IP_SALT", "")

		err := runRoot()
		assert.Error(t, err)
	})
}
