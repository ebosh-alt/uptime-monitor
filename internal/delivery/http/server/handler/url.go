package handler

import (
	"net/http"
	"strconv"

	"uptime-monitor/internal/entities"
	"uptime-monitor/internal/usecase/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) AddUrl(c *gin.Context) {
	var req entities.CreateUrlRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warnw("failed to bind create url request", zap.Error(err))
		writeError(c, response.ErrRequestBody)
		return
	}

	if req.Url == "" {
		h.log.Infow("url is empty")
		writeError(c, response.ErrUrlEmpty)
		return
	}

	created, err := h.uc.UrlCreate(c.Request.Context(), &entities.Url{Url: &req.Url})
	if err != nil {
		h.log.Errorw("failed to create url", zap.Error(err), "url", req.Url)
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, entities.CreateUrlResponse{
		Id:        *created.Id,
		Url:       *created.Url,
		CreatedAt: *created.CreatedAt,
		Active:    *created.Active,
	})
}

func (h *Handler) UrlDelete(c *gin.Context) {
	var req entities.DeleteUrlRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warnw("failed to bind delete url request", zap.Error(err))
		writeError(c, response.ErrRequestBody)
		return
	}

	if req.Url == "" {
		h.log.Infow("url is empty")
		writeError(c, response.ErrUrlEmpty)
		return
	}
	err := h.uc.UrlDelete(c.Request.Context(), &entities.Url{Url: &req.Url})
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) DeactivateUrl(c *gin.Context) {
	var req entities.DeactivateUrlRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warnw("failed to bind deactivate url request", zap.Error(err))
		writeError(c, response.ErrRequestBody)
		return
	}

	if req.Url == "" {
		h.log.Infow("url is empty")
		writeError(c, response.ErrUrlEmpty)
		return
	}

	err := h.uc.DeactivateUrl(c.Request.Context(), &entities.Url{Url: &req.Url})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) ActivateUrl(c *gin.Context) {
	var req entities.ActivateUrlRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warnw("failed to bind activate url request", zap.Error(err))
		writeError(c, response.ErrRequestBody)
		return
	}

	if req.Url == "" {
		h.log.Infow("url is empty")
		writeError(c, response.ErrUrlEmpty)
		return
	}

	err := h.uc.ActivateUrl(c.Request.Context(), &entities.Url{Url: &req.Url})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) UrlHistory(c *gin.Context) {
	urlID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Warnw("invalid url id", zap.Error(err), "id", c.Param("id"))
		writeError(c, response.ErrRequestBody)
		return
	}

	history, err := h.uc.UrlHistory(c.Request.Context(), urlID)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, entities.UrlHistoryResponse{History: history})
}

func (h *Handler) ListUrls(c *gin.Context) {
	urls, err := h.uc.ListUrls(c.Request.Context())
	if err != nil {
		h.log.Errorw("failed to list urls", zap.Error(err))
		writeError(c, err)
		return
	}
	respUrls := make(entities.ListUrlsResponse, 0, len(urls))
	for _, url := range urls {
		respUrls = append(respUrls, *url)
	}

	c.JSON(http.StatusOK, respUrls)
}
