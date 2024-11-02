package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"time"
)

func PrepareConnection(ctx context.Context, config Config, logger *zap.Logger) (*sql.DB, error) {
	// Validate Config
	if err := config.ValidateWithContext(ctx); err != nil {
		return nil, fmt.Errorf("validate Postgres config: %w", err)
	}

	db, err := sql.Open("postgres", config.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Verify connection
	ctxPing, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()

	if err := db.PingContext(ctxPing); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize Database
	if err := ensureTables(ctx, db, config, logger); err != nil {
		return nil, fmt.Errorf("ensure tables: %w", err)
	}

	return db, nil
}

func ensureTables(ctx context.Context, db *sql.DB, cfg Config, logger *zap.Logger) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, cfg.EnsureIdxTimeout)
	defer cancel()

	queries := []string{
		// Create URLs table
		`
        CREATE TABLE IF NOT EXISTS urls (
            id UUID PRIMARY KEY,
            original_url TEXT NOT NULL,
            short_url VARCHAR(10) UNIQUE NOT NULL,
            custom_alias VARCHAR(30) UNIQUE,
            created_at TIMESTAMP NOT NULL,
            expires_at TIMESTAMP
        );
        `,
		// Create Redirect Logs table
		`
        CREATE TABLE IF NOT EXISTS redirect_logs (
            id UUID PRIMARY KEY,
            short_url VARCHAR(10) NOT NULL,
            accessed_at TIMESTAMP NOT NULL,
            referrer TEXT,
            FOREIGN KEY (short_url) REFERENCES urls(short_url)
        );
        `,
		// Create index on created_at for expiration
		`
        CREATE INDEX IF NOT EXISTS idx_urls_expires_at ON urls (expires_at);
        `,
	}

	for _, query := range queries {
		if _, err := db.ExecContext(ctxTimeout, query); err != nil {
			logger.Error("Failed to execute query", zap.String("query", query), zap.Error(err))
			return fmt.Errorf("execute query: %w", err)
		}
	}

	return nil
}
