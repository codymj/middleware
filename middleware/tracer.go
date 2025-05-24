package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"

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

		// Store client IP in context.
		client := r.Header.Get("X-Forwarded-For")
		if client == "" {
			client = r.Header.Get("X-Real-IP")
		}
		if client == "" {
			client = r.RemoteAddr
		}
		if strings.Contains(client, ":") {
			client, _, _ = net.SplitHostPort(client)
		}
		ctx := context.WithValue(r.Context(), ClientContextKey, client)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
