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
			logger.Info("HTTP Request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Int("status", getStatusCode(w)),
				zap.Duration("duration", duration),
			)
		})
	}
}

// Helper function to get status code from ResponseWriter
func getStatusCode(w http.ResponseWriter) int {
	if rw, ok := w.(interface{ Status() int }); ok {
		return rw.Status()
	}
	return http.StatusOK
}
