package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func RequestLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			logger.Info("Request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Duration("duration", duration),
			)
		})
	}
}
