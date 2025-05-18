package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
)

const ClientContextKey ContextKey = "client"

// Parameters for the tracer middleware.
type LoggerParams struct {
	ServiceName string
}

// Middleware for logging request and response information.
func Logger(params LoggerParams) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// For logging client IP and setting it in context.
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

			// Create a logger to pass into the request context.
			logger := zerolog.Ctx(r.Context()).With().
				Timestamp().
				Str("traceId", r.Header.Get("X-Trace-Id")).
				Str("service", params.ServiceName).
				Str("client", client).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Logger()

			// Note beginning of request.
			logger.Info().Msg("init")

			// Wrap the response writer to capture the status code.
			sr := &statusRecorder{
				ResponseWriter: w,
				status:         200,
			}

			// Do work.
			ctx = logger.WithContext(ctx)
			next.ServeHTTP(sr, r.WithContext(ctx))

			// Note end of request.
			logger.Info().
				Int("status", sr.status).
				Int("size", sr.size).
				Msg("done")
		})
	}
}
