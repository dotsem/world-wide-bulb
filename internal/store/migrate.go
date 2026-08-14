// Package store provides database access, generated queries, and migration execution.
package store

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"strings"
)

//go:embed schema/*.sql
var schemaFS embed.FS

// Migrate executes all schema SQL files against the database in alphabetical order.
func Migrate(ctx context.Context, db *sql.DB) error {
	entries, err := fs.ReadDir(schemaFS, "schema")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		content, err := schemaFS.ReadFile("schema/" + entry.Name())
		if err != nil {
			return err
		}
		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			if strings.Contains(err.Error(), "duplicate column name") {
				continue
			}
			return err
		}
	}
	return nil
}
