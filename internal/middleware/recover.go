package middleware

import (
	"go.uber.org/zap"
	"net/http"
)

func Recoverer(logger *zap.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered", zap.Any("error", err))
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
		})
	}
}
