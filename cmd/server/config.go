package main

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kelseyhightower/envconfig"
	"github.com/vladislavprovich/url-shortener/internal/handler"
	"github.com/vladislavprovich/url-shortener/internal/repository/postgres"
)

type Config struct {
	Server   handler.Config
	Database postgres.Config
	Logger   LoggerConfig
	// Add other configs as needed
}

type LoggerConfig struct {
	Level string `envconfig:"LOG_LEVEL" default:"development"`
}

func LoadConfig(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	if err := cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Server),
		validation.Field(&c.Database),
		validation.Field(&c.Logger),
		// Add other configs as needed
	)
}

func (c LoggerConfig) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &c,
		validation.Field(&c.Level, validation.Required),
	)
}
