package services

import "github.com/vladislavprovich/url-shortener/internal/models"

type URLService interface {
	CreateShortURL(req models.ShortenRequest) (string, error)
	GetOriginalURL(shortURL string) (string, error)
	// TODO redirect
	// stats
}

// ShortenRequest -> http handler -> service -> repo db
// repo db -> service -> http handler -> response
