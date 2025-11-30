package handler

import (
	"errors"
	"net/http"
	"uptime-monitor/internal/usecase/response"

	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	code := response.CodeUnknown
	msg := "internal error"

	switch {
	case errors.Is(err, response.ErrCreateUrl):
		status = http.StatusInternalServerError
		code = response.CodeErrCreateUrl
		msg = "Error creating URL"
	case errors.Is(err, response.ErrRequestBody):
		status = http.StatusBadRequest
		code = response.CodeInvalidArgument
		msg = "Invalid request body"
	case errors.Is(err, response.ErrUrlEmpty):
		status = http.StatusUnprocessableEntity
		code = response.CodeEmptyUrl
		msg = "Url is empty"
	case errors.Is(err, response.ErrAlreadyExists):
		status = http.StatusConflict
		code = response.CodeErrUrlExists
		msg = "Url already exists"
	case errors.Is(err, response.ErrNotFoundUrl):
		status = http.StatusNotFound
		code = response.CodeNotFound
		msg = "Url not found"
	case errors.Is(err, response.ErrActivateUrl):
		status = http.StatusInternalServerError
		code = response.CodeErrActivateUrl
		msg = "Error activating URL"
	case errors.Is(err, response.ErrDeactivateUrl):
		status = http.StatusInternalServerError
		code = response.CodeErrDeactivateUrl
		msg = "Error deactivating URL"
	case errors.Is(err, response.ErrListUrls):
		status = http.StatusInternalServerError
		code = response.CodeErrListUrls
		msg = "Error listing URLs"
	case errors.Is(err, response.ErrListHistory):
		status = http.StatusInternalServerError
		code = response.CodeErrListHistory
		msg = "Error listing URL history"
	case errors.Is(err, response.ErrUnknown):
		status = http.StatusInternalServerError
		code = response.CodeUnknown
		msg = "Unknown error"
	case errors.Is(err, response.ErrDeleteUrl):
		status = http.StatusInternalServerError
		code = response.CodeErrDeleteUrl
		msg = "Error deleting URL"

	}

	c.JSON(status, response.ErrorResponse{Error: response.Error{Code: code, Message: msg}})
}
