package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTPRequestsTotal counts how many HTTP requests the backend has handled.
var HTTPRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		// This is the metric name Prometheus will see.
		Name: "http_requests_total",

		// This explains what the metric means.
		Help: "Total number of HTTP requests handled by the backend.",
	},
	[]string{"method", "path", "status"},
)

// HTTPRequestDurationSeconds measures how long HTTP requests take.
var HTTPRequestDurationSeconds = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		// This is the metric name Prometheus will see.
		Name: "http_request_duration_seconds",

		// This explains what the metric means.
		Help: "Duration of HTTP requests handled by the backend in seconds.",

		// These buckets group requests by speed.
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "path", "status"},
)

// HealthChecksTotal counts how many background health checks have been completed.
var HealthChecksTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		// This is the metric name Prometheus will see.
		Name: "health_checks_total",

		// This explains what the metric means.
		Help: "Total number of service health checks completed by the background worker.",
	},
	[]string{"status"},
)

// HealthCheckDurationSeconds measures how long monitored services take to respond.
var HealthCheckDurationSeconds = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		// This is the metric name Prometheus will see.
		Name: "health_check_duration_seconds",

		// This explains what the metric means.
		Help: "Response time of monitored service health checks in seconds.",

		// These buckets are useful for response times.
		Buckets: prometheus.DefBuckets,
	},
	[]string{"status"},
)

// RecordHTTPRequest stores metrics for one backend HTTP request.
func RecordHTTPRequest(method string, path string, statusCode int, duration time.Duration) {
	// This protects the metric from having an empty path.
	if path == "" {
		path = "unknown"
	}

	// This converts the numeric HTTP status code into text because Prometheus labels are strings.
	status := strconv.Itoa(statusCode)

	// This increases the request counter by 1.
	HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()

	// This records how long the request took in seconds.
	HTTPRequestDurationSeconds.WithLabelValues(method, path, status).Observe(duration.Seconds())
}

// RecordHealthCheck stores metrics for one background service health check.
func RecordHealthCheck(status string, duration time.Duration) {
	// This protects the metric from having an empty status.
	if status == "" {
		status = "unknown"
	}

	// This increases the health check counter by 1.
	HealthChecksTotal.WithLabelValues(status).Inc()

	// This records how long the monitored service took to respond in seconds.
	HealthCheckDurationSeconds.WithLabelValues(status).Observe(duration.Seconds())
}
