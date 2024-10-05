package repository

import (
	"database/sql"
	"github.com/vladislavprovich/url-shortener/internal/models"
)

type URLRepository interface {
	SaveURL(url models.URL) error
	GetURL(shortURL string) (models.URL, error)
	// TODO redirect log
	// stats
}

type urlRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) URLRepository {
	return &urlRepository{db: db}
}

func (u urlRepository) SaveURL(url models.URL) error {
	//TODO implement me
	panic("implement me")
}

func (u urlRepository) GetURL(shortURL string) (models.URL, error) {
	//TODO implement me
	panic("implement me")
}
