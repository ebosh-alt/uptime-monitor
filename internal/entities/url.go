package entities

import "time"

type Url struct {
	Id        *int64     `json:"id"`
	Url       *string    `json:"url"`
	Active    *bool      `json:"active"`
	CreatedAt *time.Time `json:"created_at"`
}

type CreateUrlRequest struct {
	Url string `json:"url"`
}

type DeleteUrlRequest struct {
	Url string `json:"url"`
}

type DeactivateUrlRequest struct {
	Url string `json:"url"`
}

type ActivateUrlRequest struct {
	Url string `json:"url"`
}

type CreateUrlResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Url       string    `json:"url"`
	Id        int64     `json:"id"`
	Active    bool      `json:"active"`
}

type ListUrlsResponse []Url
