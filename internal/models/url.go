package models

import (
	"time"
)

type URL struct {
	ID          string     `json:"id"`
	OriginalURL string     `json:"original_url"`
	ShortURL    string     `json:"short_url"`
	CustomAlias *string    `json:"custom_alias,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiredAt   *time.Time `json:"expired_at,omitempty"`
}

type ShortenRequest struct {
	URL         string  `json:"url" validate:"required,url"`
	CustomAlias *string `json:"custom_alias,omitempty" validate:"omitempty,alphanum"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type StatsResponse struct {
	RedirectCount int        `json:"redirect_count"`
	CreatedAt     time.Time  `json:"created_at"`
	LastAccessed  *time.Time `json:"last_accessed,omitempty"`
	Referrers     []string   `json:"referrers,omitempty"`
}
