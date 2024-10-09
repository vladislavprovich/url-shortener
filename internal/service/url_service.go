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
	short := shortener.GeneratorShortURL()
	log.Println(short)
	_, err := s.repo.GetURL(short)
	if err != nil {
		return "", fmt.Errorf("short url already exists %s ", short)
	}
	return short, nil
}

func (s urlService) GetOriginalURL(shortURL string) (string, error) {

	originalUrl, err := s.repo.GetURL(shortURL)
	if err != nil {
		return "", fmt.Errorf("short url not found %s", shortURL)
	}
	return originalUrl.OriginalURL, nil

}

func (s urlService) LogRedirect(shortURL, referrer string) error {

	id := uuid.NewString()

	logEntry := models.RedirectLog{
		ID:         id,
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
		CreatedAt:     status.AccessedAt,
		LasrAccessed:  status.AccessedAt,
		Referrers:     referrers,
	}

	//var referrers []string
	//// TODO maybe remove if...
	//if status.Referrer != nil {
	//	referrers = append(referrers, *status.Referrer)
	//}
	////------------------------
	//response := models.StatsResponce{
	//	RedirectCount: 1, // ?????
	//	CreatedAt:     status.AccessedAt,
	//	LasrAccessed:  status.AccessedAt,
	//	Referrers:     referrers,
	//}

	return response, nil
}

// ShortenRequest -> http handler -> service -> repo db
// repo db -> service -> http handler -> response
