package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vladislavprovich/url-shortener/internal/models"
)

func TestSaveURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewURLRepository(db)

	url := models.URL{
		ID:          "uuid",
		OriginalURL: "https://example.com",
		ShortURL:    "abc123",
		CreatedAt:   time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`
        INSERT INTO urls (id, original_url, short_url, custom_alias, created_at, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `)).
		WithArgs(url.ID, url.OriginalURL, url.ShortURL, url.CustomAlias, url.CreatedAt, url.ExpiredAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveURL(nil, url)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewURLRepository(db)

	shortURL := "abc123"
	url := models.URL{
		ID:          "uuid",
		OriginalURL: "https://example.com",
		ShortURL:    shortURL,
		CreatedAt:   time.Now(),
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT id, original_url, short_url, custom_alias, created_at, expires_at
        FROM urls WHERE short_url = $1 OR custom_alias = $1
    `)).
		WithArgs(shortURL).
		WillReturnRows(sqlmock.NewRows([]string{"id", "original_url", "short_url", "custom_alias", "created_at", "expires_at"}).
			AddRow(url.ID, url.OriginalURL, url.ShortURL, url.CustomAlias, url.CreatedAt, url.ExpiredAt))

	result, err := repo.GetURL(nil, shortURL)
	assert.NoError(t, err)
	assert.Equal(t, url, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveRedirectLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewURLRepository(db)

	logEntry := models.RedirectLog{
		ID:         "uuid",
		ShortURL:   "abc123",
		AccessedAt: time.Now(),
		Referrer:   nil,
	}

	mock.ExpectExec(regexp.QuoteMeta(`
        INSERT INTO redirect_logs (id, short_url, accessed_at, referrer)
        VALUES ($1, $2, $3, $4)
    `)).
		WithArgs(logEntry.ID, logEntry.ShortURL, logEntry.AccessedAt, logEntry.Referrer).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveRedirectLog(nil, logEntry)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewURLRepository(db)

	shortURL := "abc123"
	createdAt := time.Now()
	redirectCount := 5
	lastAccessed := time.Now().Add(-time.Hour)
	referrers := []string{"https://referrer1.com", "https://referrer2.com"}

	// Mock for created_at
	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT created_at FROM urls WHERE short_url = $1 OR custom_alias = $1
    `)).
		WithArgs(shortURL).
		WillReturnRows(sqlmock.NewRows([]string{"created_at"}).
			AddRow(createdAt))

	// Mock for redirect logs
	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT COUNT(*), MAX(accessed_at) FROM redirect_logs WHERE short_url = $1
    `)).
		WithArgs(shortURL).
		WillReturnRows(sqlmock.NewRows([]string{"count", "max"}).
			AddRow(redirectCount, lastAccessed))

	// Mock for referrers
	rows := sqlmock.NewRows([]string{"referrer"})
	for _, ref := range referrers {
		rows.AddRow(ref)
	}
	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT DISTINCT referrer FROM redirect_logs WHERE short_url = $1 AND referrer IS NOT NULL
    `)).
		WithArgs(shortURL).
		WillReturnRows(rows)

	stats, err := repo.GetStats(nil, shortURL)
	assert.NoError(t, err)
	assert.Equal(t, redirectCount, stats.RedirectCount)
	assert.Equal(t, createdAt, stats.CreatedAt)
	assert.Equal(t, &lastAccessed, stats.LastAccessed)
	assert.Equal(t, referrers, stats.Referrers)
	assert.NoError(t, mock.ExpectationsWereMet())
}
