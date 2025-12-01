package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"uptime-monitor/internal/entities"
	"uptime-monitor/internal/repository/response"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const (
	urlCreateSql      = `INSERT INTO urls (url) VALUES ($1) returning id, url, active, created_at;`
	urlDeleteSql      = `DELETE FROM urls where url = $1;`
	urlSelectByUrlSql = `SELECT id, url, active, created_at FROM urls where url = $1;`
	urlSelectByIDSql  = `SELECT id, url, active, created_at FROM urls where id = $1;`
	urlActivateSql    = `UPDATE urls SET active = true WHERE url = $1;`
	urlDeactivateSql  = `UPDATE urls SET active = false WHERE url = $1;`
	urlListSql        = `SELECT id, url, active, created_at FROM urls;`
	urlHistoryInsert  = `INSERT INTO urls_history (url_id, latency_ms, status_code) VALUES ($1, $2, $3);`
	urlHistoryList    = `SELECT id, url_id, latency_ms, status_code, created_at FROM urls_history WHERE url_id = $1 ORDER BY created_at DESC LIMIT $2;`
)

func (p *Postgres) UrlCreate(ctx context.Context, url *entities.Url) (*entities.Url, error) {
	var (
		id        int64
		urlStr    string
		active    bool
		createdAt time.Time
	)

	err := p.db.QueryRow(ctx, urlCreateSql, url.Url).Scan(
		&id,
		&urlStr,
		&active,
		&createdAt,
	)
	if err != nil {
		p.log.Errorw("failed to create url", "error", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, response.ErrAlreadyExists
		}
		return nil, response.ErrCreateUrl
	}

	createdUrl := &entities.Url{
		Id:        &id,
		Url:       &urlStr,
		Active:    &active,
		CreatedAt: &createdAt,
	}

	p.log.Infow("created url", "id", createdUrl.Id, "url", createdUrl.Url, "active", createdUrl.Active, "created_at", createdUrl.CreatedAt)
	return createdUrl, nil
}

func (p *Postgres) UrlDelete(ctx context.Context, url *entities.Url) error {
	cmdTag, err := p.db.Exec(ctx, urlDeleteSql, *url.Url)
	if err != nil {
		p.log.Errorw("failed to delete url", "error", zap.Error(err))
		return response.ErrDeleteUrl
	}

	if cmdTag.RowsAffected() == 0 {
		p.log.Infow("url not found", "url", *url.Url)
		return response.ErrNotFoundUrl
	}

	return nil
}

func (p *Postgres) DeactivateUrl(ctx context.Context, url *entities.Url) error {
	cmdTag, err := p.db.Exec(ctx, urlDeactivateSql, *url.Url)
	if err != nil {
		p.log.Errorw("failed to deactivate url", "error", zap.Error(err))
		return response.ErrDeactivateUrl
	}

	if cmdTag.RowsAffected() == 0 {
		p.log.Infow("url not found", "url", *url.Url)
		return response.ErrNotFoundUrl
	}

	return nil
}

func (p *Postgres) ActivateUrl(ctx context.Context, url *entities.Url) error {
	cmdTag, err := p.db.Exec(ctx, urlActivateSql, *url.Url)
	if err != nil {
		p.log.Errorw("failed to activate url", "error", zap.Error(err))
		return response.ErrActivateUrl
	}

	if cmdTag.RowsAffected() == 0 {
		p.log.Infow("url not found", "url", *url.Url)
		return response.ErrNotFoundUrl
	}

	return nil
}

func (p *Postgres) Url(ctx context.Context, url *entities.Url) (*entities.Url, error) {
	var (
		query string
		arg   any
	)

	switch {
	case url.Id != nil:
		query = urlSelectByIDSql
		arg = *url.Id
	case url.Url != nil:
		query = urlSelectByUrlSql
		arg = *url.Url
	default:
		return nil, response.ErrNotFoundUrl
	}

	var (
		id        int64
		urlStr    string
		active    bool
		createdAt time.Time
	)

	err := p.db.QueryRow(ctx, query, arg).Scan(
		&id,
		&urlStr,
		&active,
		&createdAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Infow("url not found", "url", arg)
			return nil, response.ErrNotFoundUrl
		}
		p.log.Errorw("failed to get url", "error", err)
		return nil, err
	}

	return &entities.Url{
		Id:        &id,
		Url:       &urlStr,
		Active:    &active,
		CreatedAt: &createdAt,
	}, nil
}

func (p *Postgres) ListUrls(ctx context.Context) ([]*entities.Url, error) {
	rows, err := p.db.Query(ctx, urlListSql)

	if err != nil {
		p.log.Errorw("failed to get urls", "error", err)
		return nil, response.ErrListUrls
	}

	defer rows.Close()
	urls := make([]*entities.Url, 0)
	for rows.Next() {
		var (
			id        int64
			urlStr    string
			active    bool
			createdAt time.Time
		)
		if err := rows.Scan(&id, &urlStr, &active, &createdAt); err != nil {
			p.log.Errorw("failed to scan url", "error", err)
			return nil, response.ErrListUrls
		}
		url := &entities.Url{
			Id:        &id,
			Url:       &urlStr,
			Active:    &active,
			CreatedAt: &createdAt,
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (p *Postgres) SaveUrlHistory(ctx context.Context, history *entities.UrlHistory) error {
	if _, err := p.db.Exec(ctx, urlHistoryInsert, history.UrlID, history.LatencyMs, history.StatusCode); err != nil {
		p.log.Errorw("failed to save url history", "error", err, "url_id", history.UrlID)
		return response.ErrHistorySave
	}

	return nil
}

func (p *Postgres) UrlHistory(ctx context.Context, urlID int64) ([]entities.UrlHistory, error) {
	limit := p.cfg.Monitor.HistoryLimit
	if limit <= 0 {
		limit = 100
	}

	rows, err := p.db.Query(ctx, urlHistoryList, urlID, limit)
	if err != nil {
		p.log.Errorw("failed to get url history", "error", err, "url_id", urlID)
		return nil, response.ErrHistoryList
	}
	defer rows.Close()

	history := make([]entities.UrlHistory, 0)
	for rows.Next() {
		var (
			id        int64
			urlId     int64
			latencyMs int64
			status    int
			createdAt time.Time
		)
		if err := rows.Scan(&id, &urlId, &latencyMs, &status, &createdAt); err != nil {
			p.log.Errorw("failed to scan url history", "error", err, "url_id", urlID)
			return nil, response.ErrHistoryList
		}
		h := entities.UrlHistory{
			ID:         &id,
			UrlID:      urlId,
			LatencyMs:  latencyMs,
			StatusCode: status,
			CreatedAt:  createdAt,
		}
		history = append(history, h)
	}

	return history, nil
}
