package middleware

import (
	"net/http"

	"github.com/rs/zerolog"
)

// Parameters for the tracer middleware.
type LoggerParams struct {
	ServiceName string
}

// Middleware for logging request and response information.
func Logger(params LoggerParams) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client's IP.
			client := r.Context().Value(ClientContextKey).(string)

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
			ctx := logger.WithContext(r.Context())
			next.ServeHTTP(sr, r.WithContext(ctx))

			// Note end of request.
			logger.Info().
				Int("status", sr.status).
				Int("size", sr.size).
				Msg("done")
		})
	}
}
