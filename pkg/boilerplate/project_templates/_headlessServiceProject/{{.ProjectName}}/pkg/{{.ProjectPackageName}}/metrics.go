package {{.ProjectPackageName}}

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics holds all Prometheus metrics for the service.
type Metrics struct {
	RequestsTotal      *prometheus.CounterVec
	RequestErrorsTotal *prometheus.CounterVec
	RequestDuration    *prometheus.HistogramVec
	WorkerOperations   *prometheus.CounterVec
}

// NewMetrics creates and registers Prometheus metrics.
func NewMetrics(namespace string) (metrics *Metrics) {
	metrics = NewMetricsWithRegisterer(namespace, prometheus.DefaultRegisterer)
	return metrics
}

// NewMetricsWithRegisterer creates metrics with a specific registerer (useful for testing).
func NewMetricsWithRegisterer(namespace string, reg prometheus.Registerer) (metrics *Metrics) {
	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"endpoint", "method"},
	)

	requestErrorsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "request_errors_total",
			Help:      "Total number of HTTP request errors",
		},
		[]string{"endpoint", "method", "error_type"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration in seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"endpoint", "method"},
	)

	workerOperations := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "worker_operations_total",
			Help:      "Total number of worker operations",
		},
		[]string{"operation"},
	)

	// Register metrics
	if reg != nil {
		reg.MustRegister(requestsTotal, requestErrorsTotal, requestDuration, workerOperations)
	}

	metrics = &Metrics{
		RequestsTotal:      requestsTotal,
		RequestErrorsTotal: requestErrorsTotal,
		RequestDuration:    requestDuration,
		WorkerOperations:   workerOperations,
	}
	return metrics
}

// RecordRequest increments the request counter.
func (m *Metrics) RecordRequest(endpoint, method string) {
	if m != nil && m.RequestsTotal != nil {
		m.RequestsTotal.WithLabelValues(endpoint, method).Inc()
	}
}

// RecordRequestError increments the request error counter.
func (m *Metrics) RecordRequestError(endpoint, method, errorType string) {
	if m != nil && m.RequestErrorsTotal != nil {
		m.RequestErrorsTotal.WithLabelValues(endpoint, method, errorType).Inc()
	}
}

// RecordRequestDuration records the request duration.
func (m *Metrics) RecordRequestDuration(endpoint, method string, duration float64) {
	if m != nil && m.RequestDuration != nil {
		m.RequestDuration.WithLabelValues(endpoint, method).Observe(duration)
	}
}

// RecordWorkerOperation increments the worker operation counter.
func (m *Metrics) RecordWorkerOperation(operation string) {
	if m != nil && m.WorkerOperations != nil {
		m.WorkerOperations.WithLabelValues(operation).Inc()
	}
}