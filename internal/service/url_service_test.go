package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vladislavprovich/url-shortener/internal/models"
)

type MockURLRepository struct {
	mock.Mock
}

// created two mock method for save and ger url address.
func (m *MockURLRepository) SaveURL(_ context.Context, url models.URL) error {
	args := m.Called(url)
	return args.Error(0)
}

func (m *MockURLRepository) GetURL(_ context.Context, shortURL string) (models.URL, error) {
	args := m.Called(shortURL)
	return args.Get(0).(models.URL), args.Error(1)
}
func (m *MockURLRepository) SaveRedirectLog(_ context.Context, log models.RedirectLog) error {
	args := m.Called(log)
	return args.Error(0)
}
func (m *MockURLRepository) GetStats(_ context.Context, shortURL string) (models.StatsResponse, error) {
	args := m.Called(shortURL)
	return args.Get(0).(models.StatsResponse), args.Error(1)
}

func TestSaveURL_ValidURL_CustomAlies_Found(t *testing.T) {
	mockRepo := new(MockURLRepository)
	logger, _ := zap.NewDevelopment()
	defer func() {
		err := logger.Sync()
		if err != nil {
			logger.Warn("Error logger Sync", zap.Error(err))
		}
	}()
	service := NewURLService(mockRepo, logger)
	ctx := context.Background()

	customAlias := "userurl"

	req := models.ShortenRequest{
		URL:         "https://example.com",
		CustomAlias: &customAlias,
	}

	oldurl := models.URL{
		ShortURL:    customAlias,
		OriginalURL: "https://example.com",
	}

	mockRepo.On("GetURL", ctx, customAlias).Return(oldurl, errors.New("URL is found")).Once()

	mockRepo.On("SaveURL", ctx, mock.AnythingOfType("models.URL")).Return(errors.New("save error")).Once()
	url, err := service.CreateShortURL(ctx, req)

	require.Error(t, err)

	assert.Contains(t, err.Error(), "save error")

	assert.Empty(t, url)
	mockRepo.AssertExpectations(t)
}

func TestSaveURL_ValidURL_CustomAlies(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo, nil)
	ctx := context.Background()

	customAlias := "userurl"

	req := models.ShortenRequest{
		URL:         "https://example.com",
		CustomAlias: &customAlias,
	}

	mockRepo.On("GetURL", customAlias).Return(models.URL{}, errors.New("URL not found")).Once()

	mockRepo.On("SaveURL", mock.AnythingOfType("models.URL")).Return(nil).Once()
	url, err := service.CreateShortURL(ctx, req)

	require.NoError(t, err)
	assert.NotEmpty(t, url)
	mockRepo.AssertExpectations(t)
}

func TestSaveURL_ValidURL_NoCustomAlies(t *testing.T) {
	mockRepo := new(MockURLRepository)
	logger, _ := zap.NewDevelopment()
	defer func() {
		err := logger.Sync()
		if err != nil {
			logger.Warn("Error logger Sync", zap.Error(err))
		}
	}()
	service := NewURLService(mockRepo, logger)
	ctx := context.Background()

	req := models.ShortenRequest{
		URL: "https://example.com",
	}

	mockRepo.On("GetURL", mock.Anything).Return(models.URL{}, errors.New("URL not found"))

	mockRepo.On("SaveURL", mock.Anything).Return(nil).Once()

	url, err := service.CreateShortURL(ctx, req)
	require.NoError(t, err)
	assert.NotEmpty(t, url)
	mockRepo.AssertExpectations(t)
}

func TestGetURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo, nil)
	ctx := context.Background()

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
	result, err := service.GetOriginalURL(ctx, shortURL)
	require.NoError(t, err)
	assert.Equal(t, testurl.OriginalURL, result)
}

func TestCreateShortURL_ValidURL_NoCustomAlias(t *testing.T) {
	mockRepo := new(MockURLRepository)
	logger, _ := zap.NewDevelopment()
	defer func() {
		err := logger.Sync()
		if err != nil {
			logger.Warn("Error logger Sync", zap.Error(err))
		}
	}()
	service := NewURLService(mockRepo, logger)
	ctx := context.Background()

	req := models.ShortenRequest{
		URL: "https://example.com",
	}

	// Simulate that the generated short URL does not exist in the repository.
	mockRepo.On("GetURL", mock.Anything).Return(models.URL{}, errors.New("URL not found")).Once()
	// Simulate successful save.
	mockRepo.On("SaveURL", mock.Anything).Return(nil).Once()

	shortURL, err := service.CreateShortURL(ctx, req)

	require.NoError(t, err)
	assert.NotEmpty(t, shortURL)
	mockRepo.AssertExpectations(t)
}

func TestCreateShortURL_ValidURL_WithUniqueCustomAlias(t *testing.T) {
	mockRepo := new(MockURLRepository)
	logger, _ := zap.NewDevelopment()
	defer func() {
		err := logger.Sync()
		if err != nil {
			logger.Warn("Error logger Sync", zap.Error(err))
		}
	}()
	service := NewURLService(mockRepo, logger)
	ctx := context.Background()

	customAlias := "uniquealias"
	req := models.ShortenRequest{
		URL:         "https://example.com",
		CustomAlias: &customAlias,
	}

	// Simulate that the custom alias does not exist in the repository.
	mockRepo.On("GetURL", customAlias).Return(models.URL{}, errors.New("URL not found")).Once()
	// Simulate successful save.
	mockRepo.On("SaveURL", mock.Anything).Return(nil).Once()

	shortURL, err := service.CreateShortURL(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, customAlias, shortURL)
	mockRepo.AssertExpectations(t)
}

func TestCreateShortURL_ValidURL_WithExistingCustomAlias(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo, nil)
	ctx := context.Background()

	customAlias := "existingalias"
	req := models.ShortenRequest{
		URL:         "https://example.com",
		CustomAlias: &customAlias,
	}

	// Simulate that the custom alias already exists in the repository.
	existingURL := models.URL{
		ShortURL:    customAlias,
		OriginalURL: "https://other.com",
	}
	mockRepo.On("GetURL", customAlias).Return(existingURL, nil).Once()

	shortURL, err := service.CreateShortURL(ctx, req)

	require.Error(t, err)
	require.EqualError(t, err, "custom alias already in use")
	assert.Empty(t, shortURL)
	mockRepo.AssertExpectations(t)
}

func TestCreateShortURL_InvalidURLFormat(t *testing.T) {
	// Since URL validation is handled by the validator, and the service expects valid input.
	// this test would normally be in the handler or validator tests.
	// However, we can simulate the service handling an invalid URL if validation was missed.

	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo, nil)
	ctx := context.Background()

	req := models.ShortenRequest{
		URL: "invalid-url",
	}

	// The service does not validate the URL format; it's expected to be valid.
	// So we simulate the normal flow, but perhaps the repository returns an error.

	// Simulate that the generated short URL does not exist in the repository.
	mockRepo.On("GetURL", mock.Anything).Return(models.URL{}, errors.New("URL not found")).Once()
	// Simulate successful save.
	mockRepo.On("SaveURL", mock.Anything).Return(nil).Once()

	shortURL, err := service.CreateShortURL(ctx, req)

	// Since the service does not validate the URL, it will proceed as normal.
	require.NoError(t, err)
	assert.NotEmpty(t, shortURL)
	mockRepo.AssertExpectations(t)
}

func TestCreateShortURL_RepositoryErrorOnCheck(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo, nil)
	ctx := context.Background()

	req := models.ShortenRequest{
		URL: "https://example.com",
	}

	// Simulate an error when checking for existing short URL.
	mockRepo.On("GetURL", mock.Anything).Return(models.URL{}, errors.New("database error")).Once()

	shortURL, err := service.CreateShortURL(ctx, req)

	require.Error(t, err)
	require.EqualError(t, err, "database error")
	assert.Empty(t, shortURL)
	mockRepo.AssertExpectations(t)
}

func TestCreateShortURL_RepositoryErrorOnSave(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo, nil)
	ctx := context.Background()

	req := models.ShortenRequest{
		URL: "https://example.com",
	}

	// Simulate that the generated short URL does not exist.
	mockRepo.On("GetURL", mock.Anything).Return(models.URL{}, errors.New("URL not found")).Once()
	// Simulate an error when saving the URL.
	mockRepo.On("SaveURL", mock.Anything).Return(errors.New("save error")).Once()

	shortURL, err := service.CreateShortURL(ctx, req)

	require.Error(t, err)

	// fixed.
	// use contains.
	assert.Contains(t, err.Error(), "save error")

	assert.Empty(t, shortURL)
	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL_ExistingShortURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo, nil)
	ctx := context.Background()

	shortURL := "abc123"
	originalURL := "https://example.com"

	// Simulate that the short URL exists.
	mockRepo.On("GetURL", shortURL).Return(models.URL{
		OriginalURL: originalURL,
		ExpiredAt:   nil,
	}, nil).Once()

	url, err := service.GetOriginalURL(ctx, shortURL)

	require.NoError(t, err)
	assert.Equal(t, originalURL, url)
	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL_NonExistingShortURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo, nil)
	ctx := context.Background()

	shortURL := "nonexistent"

	// Simulate that the short URL does not exist.
	mockRepo.On("GetURL", shortURL).Return(models.URL{}, errors.New("URL not found")).Once()

	url, err := service.GetOriginalURL(ctx, shortURL)

	require.Error(t, err)

	// fixed.
	// assert.EqualError(t, err, "get short url, get url err:, URL not found").
	assert.Contains(t, err.Error(), "URL not found")

	assert.Empty(t, url)
	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL_ExpiredURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo, nil)
	ctx := context.Background()

	shortURL := "expired123"
	expiresAt := time.Now().Add(-time.Hour) // URL expired an hour ago.

	// Simulate that the short URL exists but has expired.
	mockRepo.On("GetURL", shortURL).Return(models.URL{
		OriginalURL: "https://example.com",
		ExpiredAt:   &expiresAt,
	}, nil).Once()

	url, err := service.GetOriginalURL(ctx, shortURL)

	require.Error(t, err)
	require.EqualError(t, err, "URL has expired")
	assert.Empty(t, url)
	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL_RepositoryError(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo, nil)
	ctx := context.Background()

	shortURL := "abc123"

	// Simulate a repository error.
	mockRepo.On("GetURL", shortURL).Return(models.URL{}, errors.New("database error")).Once()

	url, err := service.GetOriginalURL(ctx, shortURL)
	require.Error(t, err)

	require.EqualError(t, err, "get short url, get url err:, database error")
	// error have "database error" text.
	assert.Contains(t, err.Error(), "database error")
	// checks that the url is empty.
	assert.Empty(t, url)
	mockRepo.AssertExpectations(t)
}

func TestLogRedirect(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo, nil)
	ctx := context.Background()

	shortURL := "abc123"
	referrer := "https://referrer.com"
	mockRepo.On("SaveRedirectLog", mock.Anything).Return(nil)
	err := service.LogRedirect(ctx, shortURL, referrer)
	require.NoError(t, err)
}

func TestGetStats(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo, nil)
	ctx := context.Background()

	shortURL := "abc123"
	stats := models.StatsResponse{
		RedirectCount: 10,
		CreatedAt:     time.Now(),
		LastAccessed:  nil,
		Referrers:     []string{"https://referrer.com"},
	}
	mockRepo.On("GetStats", shortURL).Return(stats, nil)
	result, err := service.GetStats(ctx, shortURL)
	require.NoError(t, err)
	assert.Equal(t, stats, result)
}
