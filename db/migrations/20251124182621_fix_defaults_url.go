package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upFixDefaultsUrl, downFixDefaultsUrl)
}

func upFixDefaultsUrl(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `UPDATE urls SET active = TRUE WHERE active IS NULL;`); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `
            ALTER TABLE urls
                ALTER COLUMN active SET DEFAULT TRUE,
                ALTER COLUMN active SET NOT NULL,
                ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP,
                ALTER COLUMN created_at SET NOT NULL;
        `)
	return err
}

func downFixDefaultsUrl(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
            ALTER TABLE urls
                ALTER COLUMN active DROP DEFAULT,
                ALTER COLUMN active DROP NOT NULL,
                ALTER COLUMN created_at DROP DEFAULT,
                ALTER COLUMN created_at DROP NOT NULL;
        `)
	return err
}
