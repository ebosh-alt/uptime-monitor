package response

import "errors"

const (
	CodeNotFound         ErrorResponseErrorCode = "NOT_FOUND"
	CodeErrCreateUrl     ErrorResponseErrorCode = "ERROR_CREATE_URL"
	CodeInvalidArgument  ErrorResponseErrorCode = "INVALID_ARGUMENT"
	CodeErrUrlExists     ErrorResponseErrorCode = "URL_EXISTS"
	CodeEmptyUrl         ErrorResponseErrorCode = "EMPTY_URL"
	CodeErrActivateUrl   ErrorResponseErrorCode = "ERROR_ACTIVATE_URL"
	CodeErrDeactivateUrl ErrorResponseErrorCode = "ERROR_DEACTIVATE_URL"
	CodeErrListUrls      ErrorResponseErrorCode = "ERROR_LIST_URLS"
	CodeUnknown          ErrorResponseErrorCode = "UNKNOWN"
	CodeErrDeleteUrl     ErrorResponseErrorCode = "ERROR_DELETE_URL"
	CodeErrListHistory   ErrorResponseErrorCode = "ERROR_LIST_URL_HISTORY"
)

var (
	ErrUnknown       = errors.New("unknown error")
	ErrRequestBody   = errors.New("invalid request body")
	ErrUrlEmpty      = errors.New("url is empty")
	ErrNotFoundUrl   = errors.New("not found url")
	ErrActivateUrl   = errors.New("error activate url")
	ErrDeactivateUrl = errors.New("error deactivate url")
	ErrCreateUrl     = errors.New("error create url")
	ErrDeleteUrl     = errors.New("error delete url")
	ErrAlreadyExists = errors.New("url already exists")
	ErrListUrls      = errors.New("error list urls")
	ErrListHistory   = errors.New("error list url history")
)

type ErrorResponseErrorCode string

type Error struct {
	Code    ErrorResponseErrorCode `json:"code"`
	Message string                 `json:"message"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}
