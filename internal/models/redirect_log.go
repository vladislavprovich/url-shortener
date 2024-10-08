package models

import "time"

type RedirectLog struct {
	ID         string    `json:"id"`
	ShortURL   string    `json:"short_url"`
	AccessedAt time.Time `json:"accessed_at"`
	Referrer   *string   `json:"referrer,omitempty"`
}
