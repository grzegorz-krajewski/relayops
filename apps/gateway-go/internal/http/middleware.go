package http

import (
	"net/http"
	"strconv"
	"time"

	"relayops/apps/gateway-go/internal/metrics"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func newStatusRecorder(w http.ResponseWriter) *statusRecorder {
	return &statusRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := newStatusRecorder(w)

		next.ServeHTTP(rec, r)

		status := strconv.Itoa(rec.statusCode)
		path := r.URL.Path

		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		metrics.HTTPRequestDurationSeconds.WithLabelValues(r.Method, path, status).
			Observe(time.Since(start).Seconds())
	})
}
