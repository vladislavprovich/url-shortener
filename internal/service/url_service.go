package service

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vladislavprovich/url-shortener/internal/models"
	"github.com/vladislavprovich/url-shortener/internal/repository"
	"github.com/vladislavprovich/url-shortener/pkg/shortener"
)

type URLService interface {
	CreateShortURL(ctx context.Context, req models.ShortenRequest) (string, error)
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	LogRedirect(ctx context.Context, shortURL, referrer string) error
	GetStats(ctx context.Context, shortURL string) (models.StatsResponse, error)
}

type urlService struct {
	repo   repository.URLRepository
	logger *zap.Logger
}

func NewURLService(repo repository.URLRepository, logger *zap.Logger) URLService {
	return &urlService{
		repo:   repo,
		logger: logger,
	}
}

func (s *urlService) CreateShortURL(ctx context.Context, req models.ShortenRequest) (string, error) {
	s.logger.Info("service.CreateShortURL", zap.String("original_url", req.URL))

	var shortURL string
	if req.CustomAlias != nil && *req.CustomAlias != "" {
		s.logger.Info("service, custom alias provided", zap.String("custom_alias", *req.CustomAlias))
		// Check if Custom Alias is unique
		_, err := s.repo.GetURL(ctx, *req.CustomAlias)
		if err == nil {
			s.logger.Warn("service, custom alias already in use", zap.String("custom_alias", *req.CustomAlias))
			return "", errors.New("custom alias already in use")
		}

		shortURL = *req.CustomAlias
	} else {
		s.logger.Info("service, generating unique short URL")
		// Generate unique short URL
		for {
			shortURL = shortener.GeneratorShortURL()
			_, err := s.repo.GetURL(ctx, shortURL)
			if err != nil {
				if strings.Contains(err.Error(), "URL not found") {
					s.logger.Info("service, unique short URL generated", zap.String("short_url", shortURL))
					break
				}
				s.logger.Error("service, error checking short URL uniqueness", zap.Error(err))
				return "", err
			}
		}
	}

	url := models.URL{
		ID:          uuid.New().String(), //uuid.NewString() ???
		OriginalURL: req.URL,
		ShortURL:    shortURL,
		CustomAlias: req.CustomAlias,
		CreatedAt:   time.Now(),
	}

	err := s.repo.SaveURL(ctx, url)
	if err != nil {
		s.logger.Error("service, failed to save URL", zap.Error(err))
		return "", fmt.Errorf("create short url, get url err: %s", err)
	}

	s.logger.Info("service, short URL created successfully", zap.String("short_url", shortURL))
	return url.ShortURL, nil
}

func (s *urlService) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	originalUrl, err := s.repo.GetURL(ctx, shortURL)
	if err != nil {
		return "", fmt.Errorf("get short url, get url err:, %s", err)
	}
	// Check if URL has expired
	if originalUrl.ExpiredAt != nil && time.Now().After(*originalUrl.ExpiredAt) {
		return "", errors.New("URL has expired")
	}

	return originalUrl.OriginalURL, nil
}

func (s *urlService) LogRedirect(ctx context.Context, shortURL, referrer string) error {
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

	return s.repo.SaveRedirectLog(ctx, log)
}

func (s *urlService) GetStats(ctx context.Context, shortURL string) (models.StatsResponse, error) {
	return s.repo.GetStats(ctx, shortURL)
}
