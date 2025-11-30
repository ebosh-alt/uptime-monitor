package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upRenameColumnUrlsUrlToLink, downRenameColumnUrlsUrlToLink)
}

func upRenameColumnUrlsUrlToLink(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.ExecContext(ctx, `
ALTER TABLE urls RENAME COLUMN link TO url;
            `)
	return err
}

func downRenameColumnUrlsUrlToLink(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
ALTER TABLE urls RENAME COLUMN url TO link;
            `)
	return err
}
