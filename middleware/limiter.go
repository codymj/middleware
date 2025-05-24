package middleware

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// Parameters for the rate limiting middleware.
type RateLimitParams struct {
	Enabled bool
	Rps     int64
	Rdb     *redis.Client
}

// Middleware for rate limited.
func Limiter(params RateLimitParams) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if params.Enabled {
				// Get client's IP.
				ip := r.Context().Value(ClientContextKey).(string)

				// Increment IP's request count.
				count, err := params.Rdb.Incr(r.Context(), ip).Result()
				if err != nil {
					return
				}

				// Set expiration on first request for this unit of time.
				if count == 1 {
					params.Rdb.Expire(r.Context(), ip, time.Second)
				}

				// Check if rate is exceeded.
				if count > params.Rps {
					rateLimitExceeded(w)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func rateLimitExceeded(w http.ResponseWriter) {
	w.WriteHeader(http.StatusTooManyRequests)
}
