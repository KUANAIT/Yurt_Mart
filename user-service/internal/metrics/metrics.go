package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestCount tracks the total number of requests
	RequestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_service_requests_total",
			Help: "Total number of requests by endpoint and status",
		},
		[]string{"endpoint", "status"},
	)

	// RequestDuration tracks the duration of requests
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_service_request_duration_seconds",
			Help:    "Duration of requests by endpoint",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	// ActiveUsers tracks the number of active users
	ActiveUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "user_service_active_users",
			Help: "Number of active users",
		},
	)

	// ErrorCount tracks the number of errors
	ErrorCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_service_errors_total",
			Help: "Total number of errors by endpoint and type",
		},
		[]string{"endpoint", "error_type"},
	)
) 