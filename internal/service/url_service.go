package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/vladislavprovich/url-shortener/internal/models"
	"github.com/vladislavprovich/url-shortener/internal/repository"
	"github.com/vladislavprovich/url-shortener/pkg/shortener"
	"time"
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

	response := models.URL{
		ID:          uuid.NewString(),
		OriginalURL: req.URL,
		ShortURL:    shortener.GeneratorShortURL(),
		CustomAlias: req.CustomAlias,
		CreatedAt:   time.Now(),
	}

	err := s.repo.SaveURL(response)
	if err != nil {
		return "", fmt.Errorf("create short url, get url err: %s ", err)
	}
	return shortener.GeneratorShortURL(), nil
}

func (s urlService) GetOriginalURL(shortURL string) (string, error) {

	originalUrl, err := s.repo.GetURL(shortURL)
	if err != nil {
		return "", fmt.Errorf("get short url, get url err:, %s", err)
	}
	return originalUrl.OriginalURL, nil

}

func (s urlService) LogRedirect(shortURL, referrer string) error {

	logEntry := models.RedirectLog{
		ID:         uuid.NewString(),
		ShortURL:   shortURL,
		AccessedAt: time.Now(),
		Referrer:   &referrer,
	}

	err := s.repo.SaveRedirectLog(logEntry)
	if err != nil {
		return fmt.Errorf("error saving redirect log: %w", err)
	}

	return nil
}

func (s urlService) GetStats(shortURL string) (models.StatsResponce, error) {
	status, err := s.repo.GetStats(shortURL)
	if err != nil {
		return models.StatsResponce{}, fmt.Errorf("error getting stats: %w", err)
	}

	var referrers []string
	if status.Referrer != nil {
		referrers = append(referrers, *status.Referrer)
	}

	response := models.StatsResponce{
		RedirectCount: 1,
		CreatedAt:     status.AccessedAt, // TODO createdAt...
		LasrAccessed:  status.AccessedAt,
		Referrers:     referrers,
	}

	return response, nil
}

// ShortenRequest -> http handler -> service -> repo db
// repo db -> service -> http handler -> response
