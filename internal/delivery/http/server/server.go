package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"uptime-monitor/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlersInterface interface {
	RegisterRoutes()
}

type Server struct {
	http     *http.Server
	log      *zap.SugaredLogger
	cfg      *config.Config
	handlers HandlersInterface
}

func New(log *zap.SugaredLogger, cfg *config.Config, engine *gin.Engine, handlers HandlersInterface) *Server {
	httpSrv := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: engine,
	}

	return &Server{
		http:     httpSrv,
		log:      log,
		cfg:      cfg,
		handlers: handlers,
	}
}

func (s *Server) OnStart() error {
	s.handlers.RegisterRoutes()

	go func() {
		s.log.Infow("server started", "host", s.cfg.Server.Host, "port", s.cfg.Server.Port)

		err := s.http.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Errorw("failed to serve", "error", err.Error())
		}
	}()
	return nil
}

func (s *Server) OnStop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		return err
	}
	return nil
}
