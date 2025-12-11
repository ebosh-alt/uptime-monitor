package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	HTTPRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of in-flight HTTP requests",
		},
	)

	MonitorUrlsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "monitor_urls_in_flight",
			Help: "Current number of URLs being processed by monitor",
		},
	)

	MonitorLastStatusCode = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_last_status_code",
			Help: "Last observed HTTP status code per URL",
		},
		[]string{"url_id", "url"},
	)

	MonitorLastLatencyMs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_last_latency_ms",
			Help: "Last observed latency in milliseconds per URL",
		},
		[]string{"url_id", "url"},
	)

	MonitorUrlActive = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_url_active",
			Help: "Whether URL monitoring is active (1) or disabled (0)",
		},
		[]string{"url_id", "url"},
	)
)

func init() {
	prometheus.MustRegister(
		HTTPRequests,
		HTTPRequestDuration,
		HTTPRequestsInFlight,
		MonitorUrlsInFlight,
		MonitorLastStatusCode,
		MonitorLastLatencyMs,
		MonitorUrlActive,
	)
}

// Middleware для всех запросов
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		HTTPRequestsInFlight.Inc()
		defer HTTPRequestsInFlight.Dec()

		start := time.Now()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		path := c.FullPath() // корректный route вместо сырого URL

		if path == "" {
			path = "unknown"
		}

		HTTPRequests.WithLabelValues(method, path, status).Inc()
		HTTPRequestDuration.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
	}
}

// Gin handler для маршрута /metrics
func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
