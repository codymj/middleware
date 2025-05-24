package middleware

import (
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/rs/zerolog"
)

func Recover(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Catch any panics.
		defer func() {
			// Get file and line number.
			_, file, line, ok := runtime.Caller(4)
			if !ok {
				file = "unknown"
				line = 0
			}

			if rec := recover(); rec != nil {
				stacktrace := debug.Stack()
				zerolog.Ctx(r.Context()).Error().
					Interface("recover", rec).
					Str("file", file).
					Int("line", line).
					Bytes("stacktrace", stacktrace).
					Msg("panic")
			}
		}()

		// Do work.
		next.ServeHTTP(w, r)
	}
}
