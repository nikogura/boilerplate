package {{.ProjectPackageName}}

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Metrics holds all Prometheus metrics for the service.
type Metrics struct {
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
	ServerStartTime     prometheus.Gauge
}

// NewMetrics creates and registers all Prometheus metrics.
func NewMetrics() (metrics *Metrics) {
	metrics = &Metrics{
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "example_spa_http_requests_total",
				Help: "The total number of HTTP requests processed",
			},
			[]string{"method", "route", "status"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "example_spa_http_request_duration_seconds",
				Help:    "The duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "route"},
		),
		ServerStartTime: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "example_spa_server_start_time_seconds",
				Help: "The Unix timestamp when the server was started",
			},
		),
	}

	// Set server start time
	metrics.ServerStartTime.SetToCurrentTime()

	return metrics
}

// MetricsServer represents a Prometheus metrics HTTP server.
type MetricsServer struct {
	server  *http.Server
	logger  *zap.Logger
	metrics *Metrics
}

// NewMetricsServer creates a new metrics server.
func NewMetricsServer(address string, logger *zap.Logger, metrics *Metrics) (srv *MetricsServer, err error) {
	router := mux.NewRouter()

	// Register Prometheus metrics endpoint
	router.Handle("/metrics", promhttp.Handler())

	// Register health check endpoints
	router.HandleFunc("/healthz", healthzHandler).Methods("GET")
	router.HandleFunc("/readyz", readyzHandler).Methods("GET")

	httpServer := &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	srv = &MetricsServer{
		server:  httpServer,
		logger:  logger,
		metrics: metrics,
	}

	return srv, nil
}

// Start starts the metrics server.
func (s *MetricsServer) Start() (err error) {
	s.logger.Info("Starting metrics server",
		zap.String("address", s.server.Addr),
	)

	err = s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown gracefully shuts down the metrics server.
func (s *MetricsServer) Shutdown(ctx context.Context) (err error) {
	s.logger.Info("Shutting down metrics server")
	return s.server.Shutdown(ctx)
}

// healthzHandler handles health check requests.
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// readyzHandler handles readiness check requests.
func readyzHandler(w http.ResponseWriter, r *http.Request) {
	// Add any readiness checks here (e.g., database connectivity)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("READY"))
}

// InstrumentHandler wraps an HTTP handler with metrics instrumentation.
func (m *Metrics) InstrumentHandler(route string, handler http.HandlerFunc) (instrumentedHandler http.HandlerFunc) {
	instrumentedHandler = func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create response writer wrapper to capture status code
		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the handler
		handler(wrapper, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		method := r.Method
		status := http.StatusText(wrapper.statusCode)

		m.HTTPRequestsTotal.WithLabelValues(method, route, status).Inc()
		m.HTTPRequestDuration.WithLabelValues(method, route).Observe(duration)
	}

	return instrumentedHandler
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code.
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code.
func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
