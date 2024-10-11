package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vladislavprovich/url-shortener/internal/models"
)

type MockURLRepository struct {
	mock.Mock
}

// TODO add mock methods for save and get url

func (m *MockURLRepository) SaveRedirectLog(log models.RedirectLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockURLRepository) GetStats(shortURL string) (models.StatsResponse, error) {
	args := m.Called(shortURL)
	return args.Get(0).(models.StatsResponse), args.Error(1)
}

// TODO TestCreateShortURL and TestGetOriginalURL

func TestLogRedirect(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	shortURL := "abc123"
	referrer := "https://referrer.com"

	mockRepo.On("SaveRedirectLog", mock.Anything).Return(nil)

	err := service.LogRedirect(shortURL, referrer)
	assert.NoError(t, err)
}

func TestGetStats(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	shortURL := "abc123"
	stats := models.StatsResponse{
		RedirectCount: 10,
		CreatedAt:     time.Now(),
		LastAccessed:  nil,
		Referrers:     []string{"https://referrer.com"},
	}

	mockRepo.On("GetStats", shortURL).Return(stats, nil)

	result, err := service.GetStats(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, stats, result)
}
