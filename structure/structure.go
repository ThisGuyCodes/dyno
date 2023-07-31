package structure

import (
	"context"
	"database/sql"
	"fmt"
)

var tableStatements = []string{
	`CREATE TABLE IF NOT EXISTS .sequences (
		UUID TEXT UNIQUE,
		description TEXT NOT NULL,
		sequence TEXT NOT NULL
	);`,
}

func InitTables(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, statement := range tableStatements {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("could not create / replace table: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit table creation transaction: %w", err)
	}

	return nil
}
