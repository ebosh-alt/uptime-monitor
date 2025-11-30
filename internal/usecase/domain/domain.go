package domain

import (
	"context"
	"errors"

	"uptime-monitor/config"
	"uptime-monitor/internal/entities"
	"uptime-monitor/internal/usecase/port"
	"uptime-monitor/internal/usecase/response"

	repoResponse "uptime-monitor/internal/repository/response"

	"go.uber.org/zap"
)

type Usecase struct {
	cfg  *config.Config
	log  *zap.SugaredLogger
	ctx  context.Context
	repo port.UrlRepository
}

func New(cfg *config.Config, log *zap.SugaredLogger, ctx context.Context, repo port.UrlRepository) *Usecase {
	return &Usecase{
		cfg:  cfg,
		log:  log,
		ctx:  ctx,
		repo: repo,
	}
}

func (u *Usecase) Url(ctx context.Context, url *entities.Url) (*entities.Url, error) {
	result, err := u.repo.Url(ctx, url)
	if err != nil {
		if errors.Is(err, repoResponse.ErrNotFoundUrl) {
			return nil, response.ErrNotFoundUrl
		}
		u.log.Errorw("failed to get url", "error", err)
		return nil, response.ErrUnknown
	}

	return result, nil
}

func (u *Usecase) UrlHistory(ctx context.Context, urlID int64) ([]entities.UrlHistory, error) {
	if _, err := u.repo.Url(ctx, &entities.Url{Id: &urlID}); err != nil {
		if errors.Is(err, repoResponse.ErrNotFoundUrl) {
			return nil, response.ErrNotFoundUrl
		}
		u.log.Errorw("failed to load url before history", "error", err, "url_id", urlID)
		return nil, response.ErrUnknown
	}

	history, err := u.repo.UrlHistory(ctx, urlID)
	if err != nil {
		if errors.Is(err, repoResponse.ErrHistoryList) {
			return nil, response.ErrListHistory
		}
		u.log.Errorw("failed to list url history", "error", err, "url_id", urlID)
		return nil, response.ErrUnknown
	}

	return history, nil
}
