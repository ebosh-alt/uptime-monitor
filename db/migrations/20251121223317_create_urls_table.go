package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUrlsTable, downCreateUrlsTable)
}

func upCreateUrlsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
                    CREATE TABLE IF NOT EXISTS urls (
                            id SERIAL PRIMARY KEY,
                            link VARCHAR(255) UNIQUE NOT NULL,
                            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
                    );
            `)
	return err
}

func downCreateUrlsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx,
		`DROP TABLE IF EXISTS urls;`,
	)
	return err
}
