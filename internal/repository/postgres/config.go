package postgres

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	Driver             string        `envconfig:"DB_DRIVER" default:"postgres"`
	ConnectionString   string        `envconfig:"DATABASE_URL"`
	MaxOpenConnections int           `envconfig:"DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConnections int           `envconfig:"DB_MAX_IDLE_CONNS" default:"25"`
	ConnMaxLifetime    time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"5m"`
	EnsureIdxTimeout   time.Duration `envconfig:"ENSURE_IDX_TIMEOUT" default:"30s"`
}

func (c Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &c,
		validation.Field(&c.Driver, validation.Required),
		validation.Field(&c.ConnectionString, validation.Required),
		validation.Field(&c.EnsureIdxTimeout, validation.Required),
	)
}
