package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vladislavprovich/url-shortener/internal/models"
	"github.com/vladislavprovich/url-shortener/internal/repository"
	"github.com/vladislavprovich/url-shortener/pkg/shortener"
)

type URLService interface {
	CreateShortURL(req models.ShortenRequest) (string, error)
	GetOriginalURL(shortURL string) (string, error)
	LogRedirect(shortURL, referrer string) error
	GetStats(shortURL string) (models.StatsResponse, error)
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{repo: repo}
}

func (s *urlService) CreateShortURL(req models.ShortenRequest) (string, error) {
	var shortURL string
	if req.CustomAlias != nil && *req.CustomAlias != "" {
		// Check if Custom Alias is unique
		_, err := s.repo.GetURL(*req.CustomAlias)
		if err == nil {
			return "", errors.New("custom alias already in use")
		}

		shortURL = *req.CustomAlias
	} else {
		// Generate unique short URL
		for {
			shortURL = shortener.GeneratorShortURL()
			_, err := s.repo.GetURL(shortURL)
			if err != nil {
				if strings.Contains(err.Error(), "URL not found") {
					break
				}

				return "", err
			}
		}
	}

	url := models.URL{
		ID:          uuid.New().String(),
		OriginalURL: req.URL,
		ShortURL:    shortURL,
		CustomAlias: req.CustomAlias,
		CreatedAt:   time.Now(),
	}

	err := s.repo.SaveURL(url)
	if err != nil {
		return "", fmt.Errorf("create short url, get url err: %s ", err)
	}

	return url.ShortURL, nil
}

func (s *urlService) GetOriginalURL(shortURL string) (string, error) {
	originalUrl, err := s.repo.GetURL(shortURL)
	if err != nil {
		return "", fmt.Errorf("get short url, get url err:, %s", err)
	}
	// Check if URL has expired
	if originalUrl.ExpiredAt != nil && time.Now().After(*originalUrl.ExpiredAt) {
		return "", errors.New("URL has expired")
	}

	return originalUrl.OriginalURL, nil
}

func (s *urlService) LogRedirect(shortURL, referrer string) error {
	var referrerPtr *string
	if referrer != "" {
		referrerPtr = &referrer
	}

	log := models.RedirectLog{
		ID:         uuid.New().String(),
		ShortURL:   shortURL,
		AccessedAt: time.Now(),
		Referrer:   referrerPtr,
	}

	return s.repo.SaveRedirectLog(log)
}

func (s *urlService) GetStats(shortURL string) (models.StatsResponse, error) {
	return s.repo.GetStats(shortURL)
}
