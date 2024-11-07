package handler

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/vladislavprovich/url-shortener/internal/middleware"
	"go.uber.org/zap"
)

func InitRouter(urlHandler *URLHandler, logger *zap.Logger, cfg Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(middleware.Recoverer(logger))
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.CORS)
	r.Use(middleware.RateLimiter(cfg.RateLimit))

	r.Post("/shorten", urlHandler.ShortenURL)
	r.Get("/{shortURL}", urlHandler.Redirect)
	r.Get("/{shortURL}/stats", urlHandler.GetStats)

	return r
}
