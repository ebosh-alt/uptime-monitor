package usecase

import (
	"context"

	"uptime-monitor/config"
	"uptime-monitor/internal/repository"
	"uptime-monitor/internal/usecase/domain"

	"go.uber.org/zap"
)

type InterfaceUsecase interface {
	InterfaceUrl
}

func New(cfg *config.Config, log *zap.SugaredLogger, cxt context.Context, repo repository.Repository) InterfaceUsecase {
	return domain.New(cfg, log, cxt, repo)
}
