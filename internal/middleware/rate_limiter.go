package middleware

import (
	"github.com/go-chi/httprate"
	"net/http"
	"os"
	"strconv"
	"time"
)

var RateLimiter func(next http.Handler) http.Handler

func init() {
	rateLimit := 100 // default 100
	if os.Getenv("RATE_LIMITER") != "" {
		if rl, err := strconv.Atoi(os.Getenv("RATE_LIMIT")); err == nil {
			rateLimit = rl
		}
	}
	RateLimiter = httprate.LimitByIP(rateLimit, time.Minute)
}
