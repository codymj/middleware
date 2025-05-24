package middleware

import (
	"net/http"
	"slices"
)

// For setting values in context.Context.
type ContextKey string

const (
	// Set client IP.
	ClientContextKey ContextKey = "client"
)

// For chaining middleware.
type Middleware func(http.HandlerFunc) http.HandlerFunc

type Chain []Middleware

func (c Chain) Then(h http.HandlerFunc) http.HandlerFunc {
	for _, mw := range slices.Backward(c) {
		h = mw(h)
	}
	return h
}

// For wrapping responses with additional data.
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
