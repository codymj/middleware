package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

// Middleware for generating a request trace ID.
func Tracer(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a traceId for this request.
		// Set the traceId in the header to pass to other services.
		traceId := r.Header.Get("X-Trace-Id")
		if traceId == "" {
			traceId = uuid.New().String()
			r.Header.Set("X-Trace-Id", traceId)
		}

		next.ServeHTTP(w, r)
	})
}
