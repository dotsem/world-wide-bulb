// Package main is the entry point for the World Wide Bulb server.
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"world-wide-bulb/internal/api"
	"world-wide-bulb/internal/config"

	_ "modernc.org/sqlite"
)

func main() {
	if err := runRoot(); err != nil {
		os.Exit(1)
	}
}

func runRoot() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server error", slog.Any("error", err))
		return err
	}
	return nil
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

	if cfg.RetentionLimit > 0 {
		// ponytail: single background ticker for pruning; upgrade to cron scheduler if retention rules expand
		go func() {
			ticker := time.NewTicker(1 * time.Hour)
			defer ticker.Stop()
			_ = app.Queries.PruneOldToggles(ctx, cfg.RetentionLimit)
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if err := app.Queries.PruneOldToggles(ctx, cfg.RetentionLimit); err != nil {
						slog.Warn("failed to prune old toggles", slog.Any("err", err))
					}
				}
			}
		}()
	}

	srv := &http.Server{
		Addr:              ":" + cfg.BackendPort,
		Handler:           app.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("failed to start server: %w", err)
	case <-ctx.Done():
		slog.Info("shutting down server gracefully")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}
		return nil
	}
}
