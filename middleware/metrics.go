package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsParams struct {
	HttpRequestsTotal   *prometheus.CounterVec
	HttpRequestDuration *prometheus.HistogramVec
	MemoryUsage         prometheus.Gauge
}

func MetricsMiddleware(params MetricsParams) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the response writer to capture the status code.
			sr := &statusRecorder{ResponseWriter: w, status: 200}

			// Do work.
			next.ServeHTTP(sr, r)

			// Record request duration.
			duration := time.Since(start).Seconds()
			params.HttpRequestDuration.WithLabelValues(
				r.URL.Path,
				r.Method,
			).Observe(duration)

			// Record request count.
			params.HttpRequestsTotal.WithLabelValues(
				r.URL.Path,
				r.Method,
				fmt.Sprintf("%d", sr.status),
			).Inc()
		})
	}
}
