package prometheus

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type customMetrics struct {
	requestDuration *prometheus.HistogramVec
	requestSize     *prometheus.HistogramVec
}

//Metrics Add customer promtheus metrics
func Metrics(next http.Handler) http.Handler {
	var c customMetrics
	c.requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_milliseconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{10, 50, 100, 200, 500, 1000, 5000, 11000, 15000},
		},
		[]string{"code", "path", "method"},
	)
	c.requestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{10, 50, 100, 500, 1000, 5000, 11000, 15000},
		},
		[]string{"code", "path", "method"},
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		c.requestDuration.WithLabelValues("200", r.URL.Path, r.Method).Observe(float64(time.Since(start).Milliseconds()))
		c.requestSize.WithLabelValues("200", r.URL.Path, r.Method).Observe(float64(ww.BytesWritten()))
	})
}
