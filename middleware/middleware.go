package middleware

import "net/http"

// For setting values in context.Context.
type ContextKey string

const (
	ClientContextKey ContextKey = "client"
)

// Middleware type for chaining HTTP handlers.
type Middleware func(http.HandlerFunc) http.HandlerFunc

type statusRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	size, err := sr.ResponseWriter.Write(b)
	sr.size += size

	return size, err
}
