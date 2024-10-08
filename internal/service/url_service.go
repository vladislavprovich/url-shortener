package service

import (
	"github.com/vladislavprovich/url-shortener/internal/models"
	"github.com/vladislavprovich/url-shortener/internal/repository"
)

type URLService interface {
	CreateShortURL(req models.ShortenRequest) (string, error)
	GetOriginalURL(shortURL string) (string, error)
	LogRedirect(shortURL, referrer string) error
	GetStats(shortURL string) (models.StatsResponce, error)
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{repo: repo}
}

func (s urlService) CreateShortURL(req models.ShortenRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s urlService) GetOriginalURL(shortURL string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s urlService) LogRedirect(shortURL, referrer string) error {
	//TODO implement me
	panic("implement me")
}

func (s urlService) GetStats(shortURL string) (models.StatsResponce, error) {
	//TODO implement me
	panic("implement me")
}

// ShortenRequest -> http handler -> service -> repo db
// repo db -> service -> http handler -> response
