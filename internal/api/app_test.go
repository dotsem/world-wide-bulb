package api_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
	"world-wide-bulb/internal/api"
	"world-wide-bulb/internal/config"
	"world-wide-bulb/internal/utils"

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

	t.Run("verify and state follow the Turnstile setting", func(t *testing.T) {
		newApp := func(cfg *config.Config) *api.App {
			db, err := sql.Open("sqlite", ":memory:")
			require.NoError(t, err)
			t.Cleanup(func() { _ = db.Close() })
			app, err := api.NewApp(context.Background(), cfg, db)
			require.NoError(t, err)
			return app
		}
		serve := func(app *api.App, method, path, cookie string) *httptest.ResponseRecorder {
			req := httptest.NewRequest(method, path, nil)
			if cookie != "" {
				req.Header.Set("Cookie", cookie)
			}
			rec := httptest.NewRecorder()
			app.Router.ServeHTTP(rec, req)
			return rec
		}

		open := newApp(&config.Config{IPSalt: "test_salt"})
		assert.Equal(t, http.StatusNoContent, serve(open, http.MethodPost, "/api/v1/verify", "").Code)
		assert.NotContains(t, serve(open, http.MethodGet, "/api/v1/state", "").Body.String(), "turnstile_site_key")

		guarded := newApp(&config.Config{IPSalt: "test_salt", TurnstileSiteKey: "site-key", TurnstileSecret: "secret"})
		assert.Equal(t, http.StatusForbidden, serve(guarded, http.MethodPost, "/api/v1/verify", "").Code)
		assert.Contains(t, serve(guarded, http.MethodGet, "/api/v1/state", "").Body.String(), `"turnstile_site_key":"site-key"`)

		exp := strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10)
		human := "human=id." + exp + "." + utils.NewHasher("test_salt").Sign("human|192.0.2.1|id|"+exp)
		assert.Equal(t, http.StatusNoContent, serve(guarded, http.MethodPost, "/api/v1/verify", human).Code)
		assert.NotContains(t, serve(guarded, http.MethodGet, "/api/v1/state", human).Body.String(), "turnstile_site_key")
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
