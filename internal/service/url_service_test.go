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

// created two mock method for save and ger url address
func (m *MockURLRepository) SaveURL(url models.URL) error {
	args := m.Called(url)
	return args.Error(1)
}

func (m *MockURLRepository) GetURL(shortURL string) (models.URL, error) {
	args := m.Called(shortURL)
	return args.Get(0).(models.URL), args.Error(1)
}

func (m *MockURLRepository) SaveRedirectLog(log models.RedirectLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockURLRepository) GetStats(shortURL string) (models.StatsResponse, error) {
	args := m.Called(shortURL)
	return args.Get(0).(models.StatsResponse), args.Error(1)
}

func TestSaveURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	OriginalURL := "http://testurl.com"
	castomalias := "http://myurl123"
	shortURL := "http://abc123"

	shortrequesr := models.ShortenRequest{
		URL:         OriginalURL,
		CustomAlias: &castomalias,
	}

	testURL := models.URL{
		ID:          "testone",
		OriginalURL: shortrequesr.URL,
		ShortURL:    shortURL,
		CustomAlias: nil,
		CreatedAt:   time.Now(),
		ExpiredAt:   nil,
	}
	// TODO FIX THIS....
	if shortrequesr.CustomAlias != nil && *shortrequesr.CustomAlias != "" {
		mockRepo.On("GetURL", *shortrequesr.CustomAlias).Return(testURL, nil)
		shortURL = *shortrequesr.CustomAlias
		mockRepo.On("SaveURL", testURL).Return(nil)
	} else {
		mockRepo.On("GetURL", shortURL).Return(testURL, nil)
		_, err := service.CreateShortURL(shortrequesr)
		assert.NoError(t, err)
		mockRepo.On("SaveURL", testURL).Return(nil)
	}

}

func TestGetURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	castomalias := "http://myurl123"
	shortURL := "http://abc123"

	testurl := models.URL{
		ID:          "testone",
		OriginalURL: "http://testurl.com",
		ShortURL:    shortURL,
		CustomAlias: &castomalias,
		CreatedAt:   time.Now(),
		ExpiredAt:   nil,
	}

	mockRepo.On("GetURL", shortURL).Return(testurl, nil)
	result, err := service.GetOriginalURL(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, testurl.OriginalURL, result)

}

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
