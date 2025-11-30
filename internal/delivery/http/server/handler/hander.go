package handler

import (
	"uptime-monitor/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InterfaceServer interface {
	AddUrl(c *gin.Context)
	DelUrl(c *gin.Context)
	Url(c *gin.Context)
	UrlHistory(c *gin.Context)
	ListUrls(c *gin.Context)
	CreateController()
}

type Handler struct {
	engine *gin.Engine
	log    *zap.SugaredLogger
	uc     usecase.InterfaceUsecase
}

func New(log *zap.SugaredLogger, engine *gin.Engine, uc usecase.InterfaceUsecase) *Handler {
	return &Handler{log: log, engine: engine, uc: uc}
}

func (h *Handler) RegisterRoutes() {
	api := h.engine.Group("/api")
	{
		apiURL := api.Group("/url")
		{
			apiURL.POST("/", h.AddUrl)
			apiURL.POST("/activate", h.ActivateUrl)
			apiURL.POST("/deactivate", h.DeactivateUrl)
			apiURL.DELETE("/", h.UrlDelete)
			apiURL.GET("/", h.ListUrls)
			apiURL.GET("/:id/history", h.UrlHistory)
		}
	}
}
