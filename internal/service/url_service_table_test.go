package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vladislavprovich/url-shortener/internal/models"
)

func TestSaveURL_Table(t *testing.T) {
	castom := "userurl"
	tests := []struct {
		name          string
		customAlias   *string
		reqURL        string
		getURLError   error
		existingURL   models.URL
		saveURLError  error
		expectError   bool
		expectedError string
	}{
		{
			name:          "Custom alias already exists and SaveURL fails",
			customAlias:   &castom,
			reqURL:        "https://example.com",
			getURLError:   errors.New("URL is found"),
			existingURL:   models.URL{ShortURL: "userurl", OriginalURL: "https://example.com"},
			saveURLError:  errors.New("save error"),
			expectError:   true,
			expectedError: "save error",
		},
		{
			name:          "Custom alias not found and SaveURL succeeds",
			customAlias:   &castom,
			reqURL:        "https://example.com",
			getURLError:   errors.New("URL not found"),
			existingURL:   models.URL{},
			saveURLError:  nil,
			expectError:   false,
			expectedError: "",
		},
		{
			name:          "No custom alias and SaveURL succeeds",
			customAlias:   nil,
			reqURL:        "https://example.com",
			getURLError:   errors.New("URL not found"),
			existingURL:   models.URL{},
			saveURLError:  nil,
			expectError:   false,
			expectedError: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockURLRepository)
			service := NewURLService(mockRepo, nil)

			req := models.ShortenRequest{
				URL:         tt.reqURL,
				CustomAlias: tt.customAlias,
			}

			if tt.customAlias != nil {
				mockRepo.On("GetURL", *tt.customAlias).Return(tt.existingURL, tt.getURLError).Once()
			} else {
				mockRepo.On("GetURL", mock.Anything).Return(tt.existingURL, tt.getURLError).Once()
			}

			mockRepo.On("SaveURL", mock.AnythingOfType("models.URL")).Return(tt.saveURLError).Once()

			url, err := service.CreateShortURL(context.TODO(), req)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, url)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, url)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateShortURL_Table(t *testing.T) {
	castom := "userurl"
	shortURL := "https://short.com"
	tests := []struct {
		name          string
		customAlias   *string
		reqURL        string
		getURLError   error
		existingURL   models.URL
		saveURLError  error
		expectError   bool
		expectedError string
	}{
		{
			name:          "CreateShortURL_ValidURL_NoCustomAlias",
			customAlias:   nil,
			reqURL:        "https://example.com",
			getURLError:   errors.New("URL not found"),
			existingURL:   models.URL{},
			saveURLError:  nil,
			expectError:   false,
			expectedError: "",
		},
		{
			name:          "CreateShortURL_ValidURL_WithUniqueCustomAlias",
			customAlias:   &castom,
			reqURL:        "https://example.com",
			getURLError:   errors.New("URL not found"),
			existingURL:   models.URL{},
			saveURLError:  nil,
			expectError:   false,
			expectedError: "",
		},
		{
			name:          "CreateShortURL_ValidURL_WithExistingCustomAlias",
			customAlias:   &castom,
			reqURL:        "https://example.com",
			existingURL:   models.URL{ShortURL: castom, OriginalURL: "https://example.com"},
			saveURLError:  errors.New("custom alias already in use"),
			getURLError:   errors.New("custom alias already in use"),
			expectError:   true,
			expectedError: "custom alias already in use",
		},
		{
			name:          "CreateShortURL_InvalidURLFormat",
			customAlias:   nil,
			reqURL:        "invalid-url",
			existingURL:   models.URL{},
			getURLError:   errors.New("URL not found"),
			saveURLError:  nil,
			expectError:   false,
			expectedError: "",
		},
		{
			name:          "CreateShortURL_RepositoryErrorOnCheck",
			customAlias:   &shortURL,
			reqURL:        "https://example.com",
			existingURL:   models.URL{},
			getURLError:   errors.New("database error"),
			saveURLError:  errors.New("database error"),
			expectError:   true,
			expectedError: "database error",
		},
		{
			name:          "CreateShortURL_RepositoryErrorOnCheck",
			customAlias:   &castom,
			reqURL:        "https://example.com",
			existingURL:   models.URL{},
			getURLError:   errors.New("database error"),
			saveURLError:  errors.New("database error"),
			expectError:   true,
			expectedError: "database error",
		},
		{
			name:          "CreateShortURL_RepositoryErrorOnSave",
			customAlias:   nil,
			reqURL:        "https://example.com",
			existingURL:   models.URL{},
			getURLError:   errors.New("URL not found"),
			saveURLError:  errors.New("save error"),
			expectError:   true,
			expectedError: "save error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockURLRepository)
			service := NewURLService(mockRepo, nil)

			req := models.ShortenRequest{
				URL:         tt.reqURL,
				CustomAlias: tt.customAlias,
			}

			if tt.customAlias != nil {
				mockRepo.On("GetURL", *tt.customAlias).Return(tt.existingURL, tt.getURLError).Once()
			} else {
				mockRepo.On("GetURL", mock.Anything).Return(tt.existingURL, tt.getURLError).Once()
			}

			mockRepo.On("SaveURL", mock.Anything).Return(tt.saveURLError).Once()

			url, err := service.CreateShortURL(context.TODO(), req)

			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, url)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, url)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetOriginalURL_Table(t *testing.T) {
	short := "shorturl"
	orig := "https://example.com"
	expiresAt := time.Now().Add(-time.Hour)
	tests := []struct {
		name          string
		shortURL      string
		existingURL   models.URL
		getURLError   error
		expectError   bool
		expectedError string
	}{
		{
			name:          "GetOriginalURL_ExistingShortURL",
			shortURL:      short,
			existingURL:   models.URL{OriginalURL: orig, ShortURL: short},
			getURLError:   nil,
			expectError:   false,
			expectedError: "",
		},
		{
			name:          "GetOriginalURL_NonExistingShortURL",
			shortURL:      short,
			existingURL:   models.URL{},
			getURLError:   errors.New("URL not found"),
			expectError:   true,
			expectedError: "URL not found",
		},
		{
			name:          "TestGetOriginalURL_ExpiredURL",
			shortURL:      short,
			existingURL:   models.URL{OriginalURL: orig, ExpiredAt: &expiresAt},
			getURLError:   nil,
			expectError:   true,
			expectedError: "URL has expired",
		},
		{
			name:          "GetOriginalURL_RepositoryError",
			shortURL:      short,
			existingURL:   models.URL{},
			getURLError:   errors.New("database error"),
			expectError:   true,
			expectedError: "database error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockURLRepository)
			service := NewURLService(mockRepo, nil)

			mockRepo.On("GetURL", short).Return(tt.existingURL, tt.getURLError).Once()

			url, err := service.GetOriginalURL(context.TODO(), short)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, url)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, url)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
