// Package main is the entry point for the World Wide Bulb server.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"world-wide-bulb/internal/api"
	"world-wide-bulb/internal/api/config"
	"world-wide-bulb/internal/api/rest"
	"world-wide-bulb/internal/api/ws"
	"world-wide-bulb/internal/bulb"
	"world-wide-bulb/internal/store"
	"world-wide-bulb/internal/utils"

	_ "modernc.org/sqlite"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000;"); err != nil {
		return fmt.Errorf("failed to set sqlite pragmas: %w", err)
	}

	if err := store.Migrate(ctx, db); err != nil {
		return fmt.Errorf("failed to run database migration: %w", err)
	}

	queries := store.New(db)
	engine := bulb.NewEngine(ctx, queries)
	hasher := utils.NewHasher(cfg.IPSalt)
	hub := ws.NewHub()

	restHandler := rest.NewHandler(queries, engine, hub, hasher)
	wsHandler := ws.NewHandler(hub, cfg.IsProd, cfg.AllowedHosts)
	router := api.NewRouter(restHandler, wsHandler)

	if err := router.Run(":" + cfg.Port); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
