package store

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
)

//go:embed schema/*.sql
var schemaFS embed.FS

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
			return err
		}
	}
	return nil
}
