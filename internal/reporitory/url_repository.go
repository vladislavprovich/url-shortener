package reporitory

import "github.com/vladislavprovich/url-shortener/internal/models"

type URLRepository interface {
	SaveURL(url models.URL) error
	GetURL(shortURL string) (models.URL, error)
	// TODO redirect log
	// stats
}
