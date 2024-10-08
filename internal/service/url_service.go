package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/vladislavprovich/url-shortener/internal/models"
	"github.com/vladislavprovich/url-shortener/internal/repository"
	"github.com/vladislavprovich/url-shortener/pkg/shortener"
	"log"
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
	if req.URL == "" {
		return "", fmt.Errorf("no long URL")
	}
	short := shortener.GeneratorShortURL()
	log.Println(short)
	_, err := s.repo.GetURL(short)
	if err == nil {
		return "", fmt.Errorf("short URL already exists")
	}

	panic("implement me")
}

func (s urlService) GetOriginalURL(shortURL string) (string, error) {
	if shortURL == "" {
		return "", fmt.Errorf("no short URL")
	}
	originalUrl, err := s.repo.GetURL(shortURL)
	if err != nil {
		return "", err
	}
	log.Println(originalUrl)
	panic("implement me")
}

func (s urlService) LogRedirect(shortURL, referrer string) error {
	if shortURL == "" {
		return fmt.Errorf("no short URL")
	}
	//TODO referrer delete...
	if referrer == "" {
		return fmt.Errorf("no referrer")
	}

	res, err := s.repo.GetStats(shortURL)
	if err != nil {
		return err
	}
	log.Println(res)

	var id = uuid.NewString()

	logEntry := models.RedirectLog{
		ID:         id,
		ShortURL:   shortURL,
		Referrer:   &referrer,
		AccessedAt: time.Now(),
	}
	err = s.repo.SaveRedirectLog(logEntry)
	if err != nil {
		return fmt.Errorf("error saving redirect log: %w", err)
	}
	panic("implement me")
	return nil
}

func (s urlService) GetStats(shortURL string) (models.StatsResponce, error) {
	status, err := s.repo.GetStats(shortURL)
	if err != nil {
		return models.StatsResponce{}, fmt.Errorf("error getting stats: %w", err)
	}
	log.Println(status)
	panic("implement me")
}

// ShortenRequest -> http handler -> service -> repo db
// repo db -> service -> http handler -> response
