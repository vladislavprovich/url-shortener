package service

import (
	"errors"
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
	return args.Error(0)
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

func TestCreateShortURL_ValidURL_NoCustomAlias(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	req := models.ShortenRequest{
		URL: "https://example.com",
	}

	// Simulate that the generated short URL does not exist in the repository
	mockRepo.On("GetURL", mock.Anything).Return(models.URL{}, errors.New("URL not found")).Once()
	// Simulate successful save
	mockRepo.On("SaveURL", mock.Anything).Return(nil).Once()

	shortURL, err := service.CreateShortURL(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, shortURL)
	mockRepo.AssertExpectations(t)
}

func TestCreateShortURL_ValidURL_WithUniqueCustomAlias(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	customAlias := "uniquealias"
	req := models.ShortenRequest{
		URL:         "https://example.com",
		CustomAlias: &customAlias,
	}

	// Simulate that the custom alias does not exist in the repository
	mockRepo.On("GetURL", customAlias).Return(models.URL{}, errors.New("URL not found")).Once()
	// Simulate successful save
	mockRepo.On("SaveURL", mock.Anything).Return(nil).Once()

	shortURL, err := service.CreateShortURL(req)

	assert.NoError(t, err)
	assert.Equal(t, customAlias, shortURL)
	mockRepo.AssertExpectations(t)
}

func TestCreateShortURL_ValidURL_WithExistingCustomAlias(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	customAlias := "existingalias"
	req := models.ShortenRequest{
		URL:         "https://example.com",
		CustomAlias: &customAlias,
	}

	// Simulate that the custom alias already exists in the repository
	existingURL := models.URL{
		ShortURL:    customAlias,
		OriginalURL: "https://other.com",
	}
	mockRepo.On("GetURL", customAlias).Return(existingURL, nil).Once()

	shortURL, err := service.CreateShortURL(req)

	assert.Error(t, err)
	assert.EqualError(t, err, "custom alias already in use")
	assert.Empty(t, shortURL)
	mockRepo.AssertExpectations(t)
}

func TestCreateShortURL_InvalidURLFormat(t *testing.T) {
	// Since URL validation is handled by the validator, and the service expects valid input,
	// this test would normally be in the handler or validator tests.
	// However, we can simulate the service handling an invalid URL if validation was missed.

	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	req := models.ShortenRequest{
		URL: "invalid-url",
	}

	// The service does not validate the URL format; it's expected to be valid.
	// So we simulate the normal flow, but perhaps the repository returns an error.

	// Simulate that the generated short URL does not exist in the repository
	mockRepo.On("GetURL", mock.Anything).Return(models.URL{}, errors.New("URL not found")).Once()
	// Simulate successful save
	mockRepo.On("SaveURL", mock.Anything).Return(nil).Once()

	shortURL, err := service.CreateShortURL(req)

	// Since the service does not validate the URL, it will proceed as normal
	assert.NoError(t, err)
	assert.NotEmpty(t, shortURL)
	mockRepo.AssertExpectations(t)
}

func TestCreateShortURL_RepositoryErrorOnCheck(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	req := models.ShortenRequest{
		URL: "https://example.com",
	}

	// Simulate an error when checking for existing short URL
	mockRepo.On("GetURL", mock.Anything).Return(models.URL{}, errors.New("database error")).Once()

	shortURL, err := service.CreateShortURL(req)

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	assert.Empty(t, shortURL)
	mockRepo.AssertExpectations(t)
}

func TestCreateShortURL_RepositoryErrorOnSave(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	req := models.ShortenRequest{
		URL: "https://example.com",
	}

	// Simulate that the generated short URL does not exist
	mockRepo.On("GetURL", mock.Anything).Return(models.URL{}, errors.New("URL not found")).Once()
	// Simulate an error when saving the URL
	mockRepo.On("SaveURL", mock.Anything).Return(errors.New("save error")).Once()

	shortURL, err := service.CreateShortURL(req)

	assert.Error(t, err)
	// TODO fix EqualError must only check if error contain "save error"
	assert.EqualError(t, err, "create short url, get url err: save error")
	assert.Empty(t, shortURL)
	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL_ExistingShortURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	shortURL := "abc123"
	originalURL := "https://example.com"

	// Simulate that the short URL exists
	mockRepo.On("GetURL", shortURL).Return(models.URL{
		OriginalURL: originalURL,
		ExpiredAt:   nil,
	}, nil).Once()

	url, err := service.GetOriginalURL(shortURL)

	assert.NoError(t, err)
	assert.Equal(t, originalURL, url)
	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL_NonExistingShortURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	shortURL := "nonexistent"

	// Simulate that the short URL does not exist
	mockRepo.On("GetURL", shortURL).Return(models.URL{}, errors.New("URL not found")).Once()

	url, err := service.GetOriginalURL(shortURL)

	assert.Error(t, err)
	// TODO fix EqualError must only check if error contain "URL not found"
	assert.EqualError(t, err, "get short url, get url err:, URL not found")
	assert.Empty(t, url)
	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL_ExpiredURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	shortURL := "expired123"
	expiresAt := time.Now().Add(-time.Hour) // URL expired an hour ago

	// Simulate that the short URL exists but has expired
	mockRepo.On("GetURL", shortURL).Return(models.URL{
		OriginalURL: "https://example.com",
		ExpiredAt:   &expiresAt,
	}, nil).Once()

	url, err := service.GetOriginalURL(shortURL)

	assert.Error(t, err)
	assert.EqualError(t, err, "URL has expired")
	assert.Empty(t, url)
	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL_RepositoryError(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	shortURL := "abc123"

	// Simulate a repository error
	mockRepo.On("GetURL", shortURL).Return(models.URL{}, errors.New("database error")).Once()

	url, err := service.GetOriginalURL(shortURL)

	assert.Error(t, err)
	// TODO fix EqualError must only check if error contain "database error"
	assert.EqualError(t, err, "get short url, get url err:, database error")
	assert.Empty(t, url)
	mockRepo.AssertExpectations(t)
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
