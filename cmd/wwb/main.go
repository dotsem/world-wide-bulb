// Package main is the entry point for the World Wide Bulb server.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"world-wide-bulb/internal/api"
	"world-wide-bulb/internal/config"

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

	if dir := filepath.Dir(cfg.DBPath); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	app, err := api.NewApp(ctx, cfg, db)
	if err != nil {
		return err
	}

	if err := app.Router.Run(":" + cfg.BackendPort); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
