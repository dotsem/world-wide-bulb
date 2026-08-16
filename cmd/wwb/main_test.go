package main

import (
	"context"
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
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		err := run(ctx)
		assert.NoError(t, err)
	})
}
