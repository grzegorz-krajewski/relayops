package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "relayops",
			Subsystem: "gateway",
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	TasksCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "relayops",
			Subsystem: "gateway",
			Name:      "tasks_created_total",
			Help:      "Total number of accepted tasks.",
		},
	)

	TaskPublishErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "relayops",
			Subsystem: "gateway",
			Name:      "task_publish_errors_total",
			Help:      "Total number of task publish errors.",
		},
	)

	TaskPersistErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "relayops",
			Subsystem: "gateway",
			Name:      "task_persist_errors_total",
			Help:      "Total number of task persist errors.",
		},
	)

	HTTPRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "relayops",
			Subsystem: "gateway",
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request duration in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)

func MustRegister() {
	prometheus.MustRegister(
		HTTPRequestsTotal,
		TasksCreatedTotal,
		TaskPublishErrorsTotal,
		TaskPersistErrorsTotal,
		HTTPRequestDurationSeconds,
	)
}
