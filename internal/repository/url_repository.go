package repository

import (
	"database/sql"
	"github.com/vladislavprovich/url-shortener/internal/models"
)

type URLRepository interface {
	SaveURL(url models.URL) error
	GetURL(shortURL string) (models.URL, error)
	SaveRedirectLog(log models.RedirectLog) error
	GetStats(shortURL string) (models.RedirectLog, error)
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

func (u urlRepository) SaveRedirectLog(log models.RedirectLog) error {
	//TODO implement me
	panic("implement me")
}

func (u urlRepository) GetStats(shortURL string) (models.RedirectLog, error) {
	//TODO implement me
	panic("implement me")
}
