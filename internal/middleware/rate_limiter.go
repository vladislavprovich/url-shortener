package middleware

import (
	"github.com/go-chi/httprate"
	"net/http"
	"time"
)

const defaultRateLimit = 100

func RateLimiter(rateLimit int) func(next http.Handler) http.Handler {
	if rateLimit == 0 {
		rateLimit = defaultRateLimit
	}
	return httprate.LimitByIP(rateLimit, time.Minute)
}
