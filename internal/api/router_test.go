package api_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"world-wide-bulb/internal/api"
	"world-wide-bulb/internal/api/rest"
	"world-wide-bulb/internal/api/sse"
	"world-wide-bulb/internal/api/ws"
	"world-wide-bulb/internal/bulb"
	"world-wide-bulb/internal/store"
	"world-wide-bulb/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestNewRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	require.NoError(t, store.Migrate(ctx, db))

	queries := store.New(db)
	engine := bulb.NewEngine(ctx, queries)
	hub := ws.NewHub()
	broker := sse.NewBroker()
	hasher := utils.NewHasher("test_salt")

	restH := rest.NewHandler(queries, engine, hub, broker, hasher, false)
	wsH := ws.NewHandler(hub, false, []string{"*"})

	testFS := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<html>SPA</html>")},
		"test.js":    &fstest.MapFile{Data: []byte("console.log('hi');")},
	}

	isProd := false
	allowedHosts := []string{"localhost"}
	router := api.NewRouter(restH, wsH, testFS, isProd, allowedHosts)
	assert.NotNil(t, router)

	t.Run("Public API endpoint applies PublicCORS", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/state", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("Web API endpoint rejects unauthorized origin in production", func(t *testing.T) {
		prodRouter := api.NewRouter(restH, wsH, testFS, true, []string{"example.com"})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req.Header.Set("Origin", "https://evil.com")
		rec := httptest.NewRecorder()
		prodRouter.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("Static asset", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test.js", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "console.log")
	})

	t.Run("SPA Fallback route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/history", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "<html>SPA</html>")
	})

	t.Run("Production mode sets release mode", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		_ = api.NewRouter(restH, wsH, testFS, true, allowedHosts)
		assert.Equal(t, gin.ReleaseMode, gin.Mode())
	})

	t.Run("Non-production mode retains mode", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		_ = api.NewRouter(restH, wsH, testFS, false, allowedHosts)
		assert.Equal(t, gin.TestMode, gin.Mode())
	})
}
