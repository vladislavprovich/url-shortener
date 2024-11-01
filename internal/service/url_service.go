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
	s.logger.Debug("service, request url", zap.String("req_original_url", req.URL))

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
		s.logger.Debug("service, create custom alias or init", zap.String("custom_alias", *req.CustomAlias))
		s.logger.Debug("service, creating new URL", zap.String("", *req.CustomAlias))
	} else {
		s.logger.Info("service, generating unique short URL")
		// Generate unique short URL
		for {
			shortURL = shortener.GeneratorShortURL()
			s.logger.Debug("Short URL generated", zap.String("short_url", shortURL))
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

	s.logger.Debug("info URL model", zap.Any("url", url))

	err := s.repo.SaveURL(ctx, url)
	if err != nil {
		s.logger.Error("service, failed to save URL", zap.Error(err))
		return "", fmt.Errorf("create short url, get url err: %s", err)
	}

	s.logger.Info("service, short URL created successfully", zap.String("short_url", shortURL))
	s.logger.Debug("service short URL created CALLED", zap.String("short_url", shortURL))
	return url.ShortURL, nil
}

func (s *urlService) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	s.logger.Info("service.GetOriginalURL", zap.String("short_url", shortURL))
	originalUrl, err := s.repo.GetURL(ctx, shortURL)
	s.logger.Debug("service.GetOriginURL short URL ", zap.String("short_url", shortURL))
	s.logger.Debug("service.GetOriginURL originalURL ", zap.Any("original_url", originalUrl))
	if err != nil {
		s.logger.Info("service, failed to get original URL", zap.String("short_url", shortURL))
		return "", fmt.Errorf("get short url, get url err:, %s", err)
	}
	// Check if URL has expired
	if originalUrl.ExpiredAt != nil && time.Now().After(*originalUrl.ExpiredAt) {
		s.logger.Info("service, storage time has expired, URL has expired", zap.String("short_url", shortURL))
		return "", errors.New("URL has expired")
	}
	s.logger.Info("service, origin URL retrieved successfully", zap.String("original_url", originalUrl.OriginalURL))
	return originalUrl.OriginalURL, nil
}

func (s *urlService) LogRedirect(ctx context.Context, shortURL, referrer string) error {
	s.logger.Info("service.LogRedirect", zap.String("short_url", shortURL), zap.String("referrer", referrer))
	var referrerPtr *string
	if referrer != "" {
		referrerPtr = &referrer
		s.logger.Info("referrer is nil, ", zap.String("referrerPtr", *referrerPtr))
	} else {
		s.logger.Info("Referrer is empty, not assigning to referrerPtr")
	}
	s.logger.Info("Saving RedirectLog to repository", zap.String("short_url", shortURL), zap.Time("accessed_at", time.Now()))
	log := models.RedirectLog{
		ID:         uuid.New().String(),
		ShortURL:   shortURL,
		AccessedAt: time.Now(),
		Referrer:   referrerPtr,
	}
	err := s.repo.SaveRedirectLog(ctx, log)
	if err != nil {
		s.logger.Error("Error saving RedirectLog", zap.String("log_id", log.ID), zap.Error(err))
		return err
	}

	return nil
}

func (s *urlService) GetStats(ctx context.Context, shortURL string) (models.StatsResponse, error) {
	s.logger.Info("service GetStatus", zap.String("shortURL", shortURL))

	status, err := s.repo.GetStats(ctx, shortURL)
	if err != nil {
		s.logger.Info("Error getting stats", zap.Error(err))
		return models.StatsResponse{}, err
	}
	return status, nil
}
