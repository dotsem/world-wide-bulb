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
		defer func() { _ = db.Close() }()

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
		defer func() { _ = db.Close() }()

		ctx := context.Background()
		err = store.Migrate(ctx, db)
		require.NoError(t, err)

		err = store.Migrate(ctx, db)
		assert.NoError(t, err)
	})

	t.Run("returns error when context is canceled", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		require.NoError(t, err)
		defer func() { _ = db.Close() }()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = store.Migrate(ctx, db)
		assert.Error(t, err)
	})

	t.Run("successfully populates uuids when migrating database with existing rows", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		require.NoError(t, err)
		defer func() { _ = db.Close() }()

		ctx := context.Background()
		_, err = db.ExecContext(ctx, `
			CREATE TABLE toggles (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				state BOOLEAN NOT NULL,
				reason TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				ip_hash TEXT NOT NULL
			);
			INSERT INTO toggles (state, ip_hash) VALUES (1, 'hash1'), (0, 'hash2');
		`)
		require.NoError(t, err)

		err = store.Migrate(ctx, db)
		assert.NoError(t, err)

		rows, err := db.QueryContext(ctx, "SELECT uuid FROM toggles")
		require.NoError(t, err)
		defer func() { _ = rows.Close() }()

		var uuids []string
		for rows.Next() {
			var u string
			require.NoError(t, rows.Scan(&u))
			uuids = append(uuids, u)
		}
		require.Len(t, uuids, 2)
		assert.NotEmpty(t, uuids[0])
		assert.NotEmpty(t, uuids[1])
		assert.NotEqual(t, uuids[0], uuids[1])
	})
}
