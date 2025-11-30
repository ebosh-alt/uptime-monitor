package entities

import "errors"

var (
	ErrorCreateUrl   = errors.New("error creating url")
	ErrorEmptyUrl    = errors.New("url is empty")
	ErrorRequestBody = errors.New("invalid request body")
)

const (
	NotFound ErrorResponseErrorCode = "NOT_FOUND"
)

type ErrorResponseErrorCode string

type Error struct {
	Code    ErrorResponseErrorCode `json:"code"`
	Message string                 `json:"message"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}
