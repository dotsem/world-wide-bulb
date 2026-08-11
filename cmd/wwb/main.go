package main

import (
	"context"
	"database/sql"
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
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		return
	}

	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
		return
	}
	defer db.Close()

	queries := store.New(db)
	engine := bulb.NewEngine(ctx, queries)
	hasher := utils.NewHasher(cfg.IPSalt)
	hub := ws.NewHub()

	restHandler := rest.NewHandler(queries, engine, hub, hasher)
	wsHandler := ws.NewHandler(hub, cfg.IsProd, cfg.AllowedHosts)
	router := api.NewRouter(restHandler, wsHandler)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
		return
	}

}
