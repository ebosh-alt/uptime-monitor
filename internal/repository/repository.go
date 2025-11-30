package repository

import (
	"context"
	"fmt"

	"uptime-monitor/config"
	"uptime-monitor/internal/repository/postgres"
	"uptime-monitor/internal/usecase/port"

	"go.uber.org/zap"
)

type InterfaceLifecycle interface {
	OnStart(_ context.Context) error
	OnStop(_ context.Context) error
}

type Repository interface {
	InterfaceLifecycle
	port.UrlRepository
}

func New(name string, log *zap.SugaredLogger, cfg *config.Config, ctx context.Context) (Repository, error) {
	switch name {
	case "postgres":
		return postgres.New(log, cfg, ctx), nil
	default:
		return nil, fmt.Errorf("unknown repo backend: %s", name)
	}
}
