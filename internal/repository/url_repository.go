package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/vladislavprovich/url-shortener/internal/models"
)

type URLRepository interface {
	SaveURL(ctx context.Context, url models.URL) error
	GetURL(ctx context.Context, shortURL string) (models.URL, error)
	SaveRedirectLog(ctx context.Context, log models.RedirectLog) error
	GetStats(ctx context.Context, shortURL string) (models.StatsResponse, error)
}

type urlRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) URLRepository {
	return &urlRepository{db: db}
}

func (repo *urlRepository) SaveURL(ctx context.Context, url models.URL) error {
	query := `
        INSERT INTO urls (id, original_url, short_url, custom_alias, created_at, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := repo.db.ExecContext(ctx, query, url.ID, url.OriginalURL, url.ShortURL, url.CustomAlias, url.CreatedAt, url.ExpiredAt)
	return err
}

func (repo *urlRepository) GetURL(ctx context.Context, shortURL string) (models.URL, error) {
	var url models.URL
	query := `
        SELECT id, original_url, short_url, custom_alias, created_at, expires_at
        FROM urls WHERE short_url = $1 OR custom_alias = $1`
	row := repo.db.QueryRowContext(ctx, query, shortURL)
	err := row.Scan(&url.ID, &url.OriginalURL, &url.ShortURL, &url.CustomAlias, &url.CreatedAt, &url.ExpiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return url, errors.New("URL not found")
		}
		return url, err
	}
	return url, nil
}

func (repo *urlRepository) SaveRedirectLog(ctx context.Context, log models.RedirectLog) error {
	query := `
        INSERT INTO redirect_logs (id, short_url, accessed_at, referrer)
        VALUES ($1, $2, $3, $4)`
	_, err := repo.db.ExecContext(ctx, query, log.ID, log.ShortURL, log.AccessedAt, log.Referrer)
	return err
}

func (repo *urlRepository) GetStats(ctx context.Context, shortURL string) (models.StatsResponse, error) {
	var stats models.StatsResponse

	query := `
        SELECT created_at FROM urls WHERE short_url = $1 OR custom_alias = $1
    `
	row := repo.db.QueryRowContext(ctx, query, shortURL)
	err := row.Scan(&stats.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return stats, errors.New("URL not found")
		}
		return stats, err
	}

	query = `
        SELECT COUNT(*), MAX(accessed_at) FROM redirect_logs WHERE short_url = $1`

	row = repo.db.QueryRowContext(ctx, query, shortURL)
	var lastAccessed sql.NullTime
	err = row.Scan(&stats.RedirectCount, &lastAccessed)
	if err != nil {
		return stats, err
	}
	if lastAccessed.Valid {
		stats.LastAccessed = &lastAccessed.Time
	}

	query = `
        SELECT DISTINCT referrer FROM redirect_logs WHERE short_url = $1 AND referrer IS NOT NULL`

	rows, err := repo.db.QueryContext(ctx, query, shortURL)
	if err != nil {
		return stats, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			panic(err)
		}
	}()

	var referrers []string
	for rows.Next() {
		var referrer sql.NullString
		if err = rows.Scan(&referrer); err != nil {
			return stats, err
		}
		if referrer.Valid {
			referrers = append(referrers, referrer.String)
		}
	}
	stats.Referrers = referrers

	return stats, nil
}
