package rest_test

import (
	"context"
	"database/sql"
	"testing"
	"world-wide-bulb/internal/api/rest"
	"world-wide-bulb/internal/api/sse"
	"world-wide-bulb/internal/api/ws"
	"world-wide-bulb/internal/bulb"
	"world-wide-bulb/internal/store"
	"world-wide-bulb/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

type testEnv struct {
	router  *gin.Engine
	queries *store.Queries
	engine  *bulb.Engine
	broker  *sse.Broker
	hub     *ws.Hub
	db      *sql.DB
}

func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	ctx := context.Background()
	err = store.Migrate(ctx, db)
	require.NoError(t, err)

	queries := store.New(db)
	engine := bulb.NewEngine(ctx, queries)
	hasher := utils.NewHasher("test_salt")
	hub := ws.NewHub()
	broker := sse.NewBroker()

	handler := rest.NewHandler(queries, engine, hub, broker, hasher, false)

	r := gin.New()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/state", handler.GetState)
		v1.GET("/history", handler.GetHistory)
		v1.GET("/events", handler.StreamEvents)
		v1.POST("/toggle", handler.PostToggle)
		v1.POST("/reason", handler.PostReason)
	}

	return &testEnv{
		router:  r,
		queries: queries,
		engine:  engine,
		hub:     hub,
		broker:  broker,
		db:      db,
	}
}
