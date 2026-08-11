// Package api initializes HTTP routing and application dependency wiring.
package api

import (
	"context"
	"database/sql"
	"fmt"
	"world-wide-bulb/internal/api/config"
	"world-wide-bulb/internal/api/rest"
	"world-wide-bulb/internal/api/ws"
	"world-wide-bulb/internal/bulb"
	"world-wide-bulb/internal/store"
	"world-wide-bulb/internal/utils"

	"github.com/gin-gonic/gin"
)

// App encapsulates all application services and the HTTP router.
type App struct {
	Router  *gin.Engine
	DB      *sql.DB
	Queries *store.Queries
	Engine  *bulb.Engine
	Hub     *ws.Hub
}

// NewApp initializes migrations, repositories, state engine, websocket hub, and router.
func NewApp(ctx context.Context, cfg *config.Config, db *sql.DB) (*App, error) {
	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000;"); err != nil {
		return nil, fmt.Errorf("failed to set sqlite pragmas: %w", err)
	}

	if err := store.Migrate(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to run database migration: %w", err)
	}

	queries := store.New(db)
	engine := bulb.NewEngine(ctx, queries)
	hasher := utils.NewHasher(cfg.IPSalt)
	hub := ws.NewHub()

	restHandler := rest.NewHandler(queries, engine, hub, hasher)
	wsHandler := ws.NewHandler(hub, cfg.IsProd, cfg.AllowedHosts)
	router := NewRouter(restHandler, wsHandler)

	return &App{
		Router:  router,
		DB:      db,
		Queries: queries,
		Engine:  engine,
		Hub:     hub,
	}, nil
}
