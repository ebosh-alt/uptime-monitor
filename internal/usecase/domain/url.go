package domain

import (
	"context"
	"errors"

	"uptime-monitor/internal/entities"
	"uptime-monitor/internal/usecase/response"

	repoResponse "uptime-monitor/internal/repository/response"
)

func (u *Usecase) UrlCreate(ctx context.Context, url *entities.Url) (*entities.Url, error) {
	url, err := u.repo.UrlCreate(ctx, url)
	if err != nil {
		switch {
		case errors.Is(err, repoResponse.ErrAlreadyExists):
			return nil, response.ErrAlreadyExists
		case errors.Is(err, repoResponse.ErrNotFoundUrl):
			return nil, response.ErrNotFoundUrl
		case errors.Is(err, repoResponse.ErrCreateUrl):
			return nil, response.ErrCreateUrl
		default:
			u.log.Errorw("an unknown error while creating url", "error", err)
			return nil, response.ErrUnknown
		}
	}
	return url, nil
}

func (u *Usecase) UrlDelete(ctx context.Context, url *entities.Url) error {
	err := u.repo.UrlDelete(ctx, url)
	if err != nil {
		if errors.Is(err, repoResponse.ErrNotFoundUrl) {
			return response.ErrNotFoundUrl
		}
		if errors.Is(err, repoResponse.ErrDeleteUrl) {
			return response.ErrDeleteUrl
		} else {
			u.log.Errorw("an unknown error", "error", err)
			return response.ErrUnknown
		}
	}
	return nil
}

func (u *Usecase) ActivateUrl(ctx context.Context, url *entities.Url) error {
	err := u.repo.ActivateUrl(ctx, url)
	if err != nil {
		if errors.Is(err, repoResponse.ErrNotFoundUrl) {
			return response.ErrNotFoundUrl
		}
		if errors.Is(err, repoResponse.ErrActivateUrl) {
			return response.ErrActivateUrl
		} else {
			u.log.Errorw("an unknown error", "error", err)
			return response.ErrUnknown
		}
	}
	return nil
}

func (u *Usecase) DeactivateUrl(ctx context.Context, url *entities.Url) error {
	err := u.repo.DeactivateUrl(ctx, url)
	if err != nil {
		if errors.Is(err, repoResponse.ErrNotFoundUrl) {
			return response.ErrNotFoundUrl
		}
		if errors.Is(err, repoResponse.ErrDeactivateUrl) {
			return response.ErrDeactivateUrl
		} else {
			u.log.Errorw("an unknown error", "error", err)
			return response.ErrUnknown
		}
	}
	return nil
}

func (u *Usecase) ListUrls(ctx context.Context) ([]*entities.Url, error) {
	urls, err := u.repo.ListUrls(ctx)
	if err != nil {
		if errors.Is(err, repoResponse.ErrListUrls) {
			return nil, response.ErrListUrls
		} else {
			u.log.Errorw("an unknown error", "error", err)
			return nil, response.ErrUnknown
		}
	}
	return urls, nil
}
