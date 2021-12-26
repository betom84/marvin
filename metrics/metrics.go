package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	registry = prometheus.DefaultRegisterer

	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of http requests by method, endpoint and statusCode",
		},
		[]string{"method", "endpoint", "statusCode"},
	)

	requestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_duration_milliseconds",
			Help:       "Http request duration in milliseconds by method, endpoint and statusCode",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		},
		[]string{"method", "endpoint", "statusCode"},
	)
)

func init() {
	registry.MustRegister(requestsTotal, requestDuration)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wrapped := middleware.NewWrapResponseWriter(w, 0)

		defer func() {
			labels := prometheus.Labels{
				"method":     r.Method,
				"endpoint":   r.URL.Path,
				"statusCode": fmt.Sprint(wrapped.Status()),
			}

			requestsTotal.With(labels).Inc()
			requestDuration.With(labels).Observe(float64(time.Now().UnixMilli() - startTime.UnixMilli()))
		}()

		next.ServeHTTP(wrapped, r)
	})
}

func Handler() http.Handler {
	return promhttp.Handler()
}
