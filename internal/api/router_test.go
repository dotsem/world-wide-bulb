package api_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
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

	router := api.NewRouter(restH, wsH)
	assert.NotNil(t, router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/state", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}
