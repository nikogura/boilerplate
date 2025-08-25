//nolint:staticcheck // Package name matches service name requirement from prompt.xml
package {{.ProjectPackageName}}

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics for the service.
type Metrics struct {
	RequestsTotal   *prometheus.CounterVec
	ErrorsTotal     *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
}

// NewMetrics creates and registers Prometheus metrics.
func NewMetrics() (metrics *Metrics) {
	metrics = &Metrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "example_indirect_selection_requests_total",
				Help: "Total number of requests processed",
			},
			[]string{"method", "status"},
		),
		ErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "example_indirect_selection_errors_total",
				Help: "Total number of errors encountered",
			},
			[]string{"method", "error_type"},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "example_indirect_selection_request_duration_seconds",
				Help:    "Request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method"},
		),
	}

	return metrics
}
