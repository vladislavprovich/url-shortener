package service

import (
	"github.com/vladislavprovich/url-shortener/internal/models"
	"github.com/vladislavprovich/url-shortener/internal/repository"
)

type URLService interface {
	CreateShortURL(req models.ShortenRequest) (string, error)
	GetOriginalURL(shortURL string) (string, error)
	// TODO redirect
	// stats
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{repo: repo}
}

func (u urlService) CreateShortURL(req models.ShortenRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (u urlService) GetOriginalURL(shortURL string) (string, error) {
	//TODO implement me
	panic("implement me")
}

// ShortenRequest -> http handler -> service -> repo db
// repo db -> service -> http handler -> response
