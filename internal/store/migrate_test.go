package store_test

import (
	"context"
	"database/sql"
	"testing"
	"world-wide-bulb/internal/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestMigrate(t *testing.T) {
	t.Run("creates tables and indexes on fresh database", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		require.NoError(t, err)
		defer db.Close()

		ctx := context.Background()
		err = store.Migrate(ctx, db)
		require.NoError(t, err)

		var tableName string
		err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name='toggles'").Scan(&tableName)
		require.NoError(t, err)
		assert.Equal(t, "toggles", tableName)
	})

	t.Run("is idempotent when executed multiple times", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		require.NoError(t, err)
		defer db.Close()

		ctx := context.Background()
		err = store.Migrate(ctx, db)
		require.NoError(t, err)

		err = store.Migrate(ctx, db)
		assert.NoError(t, err)
	})
}
