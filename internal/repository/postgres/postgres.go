package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"uptime-monitor/config"
	"uptime-monitor/internal/usecase/port"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	_ "uptime-monitor/db/migrations"
)

type Postgres struct {
	ctx context.Context
	log *zap.SugaredLogger
	db  *pgxpool.Pool
	cfg *config.Config
}

var _ port.UrlRepository = (*Postgres)(nil)

func New(log *zap.SugaredLogger, cfg *config.Config, ctx context.Context) *Postgres {
	return &Postgres{
		ctx: ctx,
		log: log,
		cfg: cfg,
	}
}

func (p *Postgres) OnStart(_ context.Context) error {
	connectionURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.cfg.Postgres.Host,
		p.cfg.Postgres.Port,
		p.cfg.Postgres.User,
		p.cfg.Postgres.Password,
		p.cfg.Postgres.DBName,
		p.cfg.Postgres.SSLMode,
	)

	pool, err := pgxpool.New(p.ctx, connectionURL)
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", connectionURL)
	if err != nil {
		p.log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			p.log.Fatal(err)
		}
	}(db)

	if err := goose.UpContext(context.Background(), db, "db/migrations"); err != nil {
		p.log.Fatal(err)
	}
	p.log.Info("Database migration completed successfully")

	p.db = pool
	p.log.Info("Postgres started")
	return nil
}

func (p *Postgres) OnStop(_ context.Context) error {
	p.db.Close()
	return nil
}
