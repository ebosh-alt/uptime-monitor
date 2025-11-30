package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUpdateUrlsTableActive, downUpdateUrlsTableActive)
}

func upUpdateUrlsTableActive(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
                    ALTER TABLE urls ADD COLUMN active BOOLEAN DEFAULT TRUE;
            `)
	return err
}

func downUpdateUrlsTableActive(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `ALTER TABLE urls DROP COLUMN active;`)
	return err
}
