package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUrlsHistoryTableFixed, downCreateUrlsHistoryTableFixed)
}

func upCreateUrlsHistoryTableFixed(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS urls_history (
            id SERIAL PRIMARY KEY,
            url_id INTEGER NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
            latency_ms INTEGER NOT NULL,
            status_code INTEGER NOT NULL,
            created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
        CREATE INDEX IF NOT EXISTS idx_urls_history_url_id_created_at ON urls_history (url_id, created_at DESC);
    `)
	return err
}

func downCreateUrlsHistoryTableFixed(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS urls_history;`)
	return err
}
