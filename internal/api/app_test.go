package api_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"world-wide-bulb/internal/api"
	"world-wide-bulb/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestNewApp(t *testing.T) {
	t.Run("initializes App services successfully", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		require.NoError(t, err)
		defer func() { _ = db.Close() }()

		cfg := &config.Config{
			BackendPort:  "5000",
			IPSalt:       "test_salt",
			DBPath:       ":memory:",
			AllowedHosts: []string{"localhost"},
			IsProd:       false,
		}

		ctx := context.Background()
		app, err := api.NewApp(ctx, cfg, db)
		require.NoError(t, err)
		assert.NotNil(t, app)
		assert.NotNil(t, app.Router)
		assert.NotNil(t, app.DB)
		assert.NotNil(t, app.Queries)
		assert.NotNil(t, app.Engine)
		assert.NotNil(t, app.Hub)
	})

	t.Run("guards toggles with Turnstile only when it is configured", func(t *testing.T) {
		toggle := func(cfg *config.Config) *httptest.ResponseRecorder {
			db, err := sql.Open("sqlite", ":memory:")
			require.NoError(t, err)
			defer func() { _ = db.Close() }()

			app, err := api.NewApp(context.Background(), cfg, db)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			app.Router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil))
			return rec
		}

		open := toggle(&config.Config{IPSalt: "test_salt"})
		assert.Equal(t, http.StatusOK, open.Code)

		guarded := toggle(&config.Config{IPSalt: "test_salt", TurnstileSiteKey: "site-key", TurnstileSecret: "secret"})
		assert.Equal(t, http.StatusForbidden, guarded.Code)
		assert.Contains(t, guarded.Body.String(), `"turnstile_site_key":"site-key"`)
	})

	t.Run("returns error when database connection is closed", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		require.NoError(t, err)
		_ = db.Close()

		cfg := &config.Config{
			BackendPort:  "5000",
			IPSalt:       "test_salt",
			DBPath:       ":memory:",
			AllowedHosts: []string{"localhost"},
			IsProd:       false,
		}

		ctx := context.Background()
		app, err := api.NewApp(ctx, cfg, db)
		assert.Error(t, err)
		assert.Nil(t, app)
	})
}
