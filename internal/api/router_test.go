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
	hasher := utils.NewHasher("test_salt")

	restH := rest.NewHandler(queries, engine, hub, hasher, false)
	wsH := ws.NewHandler(hub, false, []string{"*"})

	testFS := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<html>SPA</html>")},
		"test.js":    &fstest.MapFile{Data: []byte("console.log('hi');")},
	}

	isProd := false
	router := api.NewRouter(restH, wsH, testFS, isProd)
	assert.NotNil(t, router)

	t.Run("API endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/state", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
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
}
