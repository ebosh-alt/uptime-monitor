package response

import "errors"

var (
	ErrCreateUrl      = errors.New("error creating url")
	ErrAlreadyExists  = errors.New("url already exists")
	ErrRequestBody    = errors.New("invalid request body")
	ErrDeleteUrl      = errors.New("error deleting url")
	ErrDeactivateUrl  = errors.New("error deactivating url")
	ErrNotFoundUrl    = errors.New("url not found")
	ErrActivateUrl    = errors.New("error activating url")
	ErrListUrls       = errors.New("error listing urls")
	ErrRecordListUrls = errors.New("error scan record by list urls")
	ErrHistorySave    = errors.New("error saving url history")
	ErrHistoryList    = errors.New("error listing url history")
)
