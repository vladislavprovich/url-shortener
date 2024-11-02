package handler

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	Port        string `envconfig:"SERVER_PORT" default:"8080"`
	BaseURL     string `envconfig:"BASE_URL" default:"http://localhost:8080"`
	ReadTimeout int    `envconfig:"SERVER_READ_TIMEOUT" default:"15"`
	RateLimit   int    `envconfig:"RATE_LIMIT" default:"100"`
}

func (c Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &c,
		validation.Field(&c.Port, validation.Required),
		validation.Field(&c.BaseURL, validation.Required),
		validation.Field(&c.RateLimit, validation.Required, validation.Min(1)),
	)
}
