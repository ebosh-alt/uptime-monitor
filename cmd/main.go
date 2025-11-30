package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"uptime-monitor/config"
	"uptime-monitor/internal/delivery/http/server"
	"uptime-monitor/internal/delivery/http/server/handler"
	"uptime-monitor/internal/delivery/http/server/middleware"
	"uptime-monitor/internal/monitor"
	"uptime-monitor/internal/repository"
	"uptime-monitor/internal/usecase"
	"uptime-monitor/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log, err := logger.New("debug")
	if err != nil {
		log.Errorw("Error initializing the logger:",
			"error", err.Error(),
		)
		os.Exit(1)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Errorw("Configuration initialization error:",
			"error", err.Error(),
		)
		os.Exit(1)
	}

	repo, err := repository.New("postgres", log, cfg, ctx)
	if err != nil {
		log.Errorw("Repository initialization error:",
			"error", err.Error(),
		)
		os.Exit(1)
	}

	err = repo.OnStart(ctx)
	if err != nil {
		log.Errorw("Configuration error:",
			"error", err.Error(),
		)
		os.Exit(1)
	}
	log.Infof("cfg: %v", cfg)
	uc := usecase.New(cfg, log, ctx, repo)
	mon := monitor.New(cfg, log, repo)
	mon.Start(ctx)
	engine := gin.Default()
	mdl := middleware.New(cfg, log, engine)
	mdl.Register()

	handlers := handler.New(log, engine, uc)

	srv := server.New(log, cfg, engine, handlers)

	if err := srv.OnStart(); err != nil {
		log.Errorw("Server startup error:",
			"error", err.Error(),
		)
		os.Exit(1)
	}

	<-ctx.Done()

	log.Infof("📦 Completion...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.OnStop(shutdownCtx); err != nil {
		log.Errorf("Server shutdown error %s", err.Error())
	}

	if err := repo.OnStop(shutdownCtx); err != nil {
		log.Errorf("Repository shutdown error %s", err.Error())
	}

	log.Infof("✅ Completed.")
}
