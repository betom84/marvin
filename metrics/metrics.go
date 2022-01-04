package metrics

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Meter interface {
	Register(r prometheus.Registerer)
}

type meter struct {
	sync.Once

	registerer              prometheus.Registerer
	requestsTotal           *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
	deviceOperationDuration *prometheus.HistogramVec
}

var singleton = meter{}

func Instance() Meter {
	return instance()
}

func instance() *meter {
	singleton.Do(func() {
		singleton.requestsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of http requests by method, endpoint and statusCode",
			},
			[]string{"method", "endpoint", "statusCode"},
		)

		singleton.requestDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_milliseconds",
				Help:    "Http request duration in milliseconds by method, endpoint and statusCode",
				Buckets: prometheus.ExponentialBuckets(5, 2, 10),
			},
			[]string{"method", "endpoint", "statusCode"},
		)

		singleton.deviceOperationDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "application_device_operation_duration_milliseconds",
				Help:    "Device operation duration in milliseconds by device and operation",
				Buckets: prometheus.ExponentialBuckets(5, 2, 10),
			},
			[]string{"device", "operation", "failed"},
		)
		singleton.Register(prometheus.DefaultRegisterer)
	})

	return &singleton
}

func (m *meter) Register(r prometheus.Registerer) {
	m.registerer = r
	m.registerer.MustRegister(
		m.requestsTotal,
		m.requestDuration,
		m.deviceOperationDuration,
	)
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

			instance().requestsTotal.With(labels).Inc()
			instance().requestDuration.With(labels).Observe(float64(time.Now().UnixMilli() - startTime.UnixMilli()))
		}()

		next.ServeHTTP(wrapped, r)
	})
}

func Handler() http.Handler {
	return promhttp.Handler()
}

func CollectDeviceOperationDuration(device, op string, err error) func() {
	opStartMs := time.Now().UnixMilli()

	return func() {
		instance().deviceOperationDuration.With(prometheus.Labels{
			"device":    device,
			"operation": op,
			"failed":    fmt.Sprint(bool(err != nil)),
		}).Observe(float64(time.Now().UnixMilli() - opStartMs))
	}
}
