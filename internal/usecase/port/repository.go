package port

import (
	"context"

	"uptime-monitor/internal/entities"
)

type UrlRepository interface {
	UrlCreate(ctx context.Context, url *entities.Url) (*entities.Url, error)
	UrlDelete(ctx context.Context, url *entities.Url) error
	ActivateUrl(ctx context.Context, url *entities.Url) error
	DeactivateUrl(ctx context.Context, url *entities.Url) error
	Url(ctx context.Context, url *entities.Url) (*entities.Url, error)
	ListUrls(ctx context.Context) ([]*entities.Url, error)
	UrlHistory(ctx context.Context, urlID int64) ([]entities.UrlHistory, error)
	SaveUrlHistory(ctx context.Context, history *entities.UrlHistory) error
}
