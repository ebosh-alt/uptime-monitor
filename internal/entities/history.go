package entities

import "time"

// UrlHistory represents a single availability check result for a URL.
// UrlID is omitted from JSON because it's an internal reference back to the urls table.
type UrlHistory struct {
	ID         *int64    `json:"id,omitempty"`
	UrlID      int64     `json:"-"`
	StatusCode int       `json:"status"`
	LatencyMs  int64     `json:"latency_ms"`
	CreatedAt  time.Time `json:"timestamp"`
}

type UrlHistoryResponse struct {
	History []UrlHistory `json:"history"`
}
