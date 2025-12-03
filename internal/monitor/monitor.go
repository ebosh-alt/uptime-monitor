package monitor

import (
	"context"
	"net/http"
	"sync"
	"time"

	"uptime-monitor/config"
	"uptime-monitor/internal/entities"
	"uptime-monitor/internal/usecase/port"

	"go.uber.org/zap"
)

type Monitor struct {
	cfg    config.MonitorConfig
	log    *zap.SugaredLogger
	repo   port.UrlRepository
	client *http.Client
	sem    chan struct{}
	maxC   int
}

func New(cfg *config.Config, log *zap.SugaredLogger, repo port.UrlRepository) *Monitor {
	monCfg := cfg.Monitor

	timeout := time.Duration(monCfg.RequestTimeoutSeconds) * time.Second
	maxConcurrency := monCfg.MaxConcurrency
	if maxConcurrency <= 0 {
		maxConcurrency = 1
	}

	return &Monitor{
		cfg:    monCfg,
		log:    log,
		repo:   repo,
		client: &http.Client{Timeout: timeout},
		sem:    make(chan struct{}, maxConcurrency),
		maxC:   maxConcurrency,
	}
}

func (m *Monitor) Start(ctx context.Context) {
	interval := time.Duration(m.cfg.IntervalSeconds) * time.Second
	if interval <= 0 {
		m.log.Warn("monitor interval is not configured; background checker disabled")
		return
	}

	m.log.Infow("starting monitor loop", "interval", interval, "timeout", m.client.Timeout, "max_concurrency", m.maxC)

	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				m.log.Info("monitor loop stopped")
				return
			case <-ticker.C:
				m.runIteration(ctx)
			}
		}
	}()
}

func (m *Monitor) runIteration(ctx context.Context) {
	urls, err := m.repo.ListUrls(ctx)
	if err != nil {
		m.log.Errorw("monitor: failed to list urls", "error", err)
		return
	}
	m.log.Infow("monitor: running iteration", "urls_count", len(urls), "max_concurrency", m.maxC)

	var wg sync.WaitGroup

	for _, url := range urls {
		if url == nil || url.Url == nil || url.Id == nil {
			m.log.Warnw("monitor: skipping url without required fields", "url", url)
			continue
		}
		if url.Active != nil && !*url.Active {
			m.log.Debugw("monitor: skipping inactive url", "url", url)
			continue
		}

		m.sem <- struct{}{}
		inFlight := len(m.sem)
		wg.Add(1)
		u := url
		go func() {
			defer wg.Done()
			m.log.Debugw("monitor: acquired slot", "url", u.Url, "in_flight", inFlight, "max", m.maxC)
			defer func() {
				<-m.sem
				m.log.Debugw("monitor: released slot", "url", u.Url, "in_flight", len(m.sem), "max", m.maxC)
			}()
			m.probe(ctx, u)
		}()
	}

	wg.Wait()
}

func (m *Monitor) probe(ctx context.Context, url *entities.Url) {
	start := time.Now()
	statusCode := 0

	reqCtx, cancel := context.WithTimeout(ctx, m.client.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, *url.Url, nil)
	if err != nil {
		m.log.Warnw("monitor: failed to build request", "error", err, "url", url.Url)
		return
	}

	resp, err := m.client.Do(req)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		m.log.Warnw("monitor: request failed", "error", err, "url", url.Url)
	} else {
		statusCode = resp.StatusCode
		_ = resp.Body.Close()
	}

	history := entities.UrlHistory{
		UrlID:      *url.Id,
		StatusCode: statusCode,
		LatencyMs:  latency,
	}

	if err := m.repo.SaveUrlHistory(ctx, &history); err != nil {
		m.log.Errorw("monitor: failed to save url history", "error", err, "url_id", url.Id)
		return
	}

	m.log.Infow("monitor: probe complete", "url", url.Url, "status_code", statusCode, "latency_ms", latency)
}
