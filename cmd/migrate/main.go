package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"uptime-monitor/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	_ "uptime-monitor/db/migrations"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <goose-command> [args]", os.Args[0])
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("failed to close db: %v", err)
		}
	}(db)

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		return
	}

	command := os.Args[1]
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	if err := goose.RunContext(context.Background(), command, db, "db/migrations", args...); err != nil {
		log.Fatalf("goose %s: %v", command, err)
	}
}
