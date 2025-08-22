package {{.ProjectPackageName}}

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Service represents the main service logic.
type Service struct {
	logger   *zap.Logger
	interval time.Duration
}

// NewService creates a new service instance.
func NewService(logger *zap.Logger, interval time.Duration) (service *Service) {
	if interval == 0 {
		interval = 5 * time.Second // Default to 5 seconds
	}

	service = &Service{
		logger:   logger,
		interval: interval,
	}
	return service
}

// Start starts the service.
func (s *Service) Start(ctx context.Context) (err error) {
	s.logger.Info("Starting {{.ProjectName}} service")

	// Start Prometheus server in goroutine on :8080
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/healthz", s.healthHandler)
		mux.HandleFunc("/readyz", s.readinessHandler)
		mux.Handle("/metrics", promhttp.Handler())
		s.logger.Info("Metrics server starting on :8080")
		serverErr := http.ListenAndServe(":8080", mux)
		if serverErr != nil {
			s.logger.Error("Metrics server error", zap.Error(serverErr))
		}
	}()

	s.logger.Info("{{.ProjectName}} service started", zap.Duration("log_interval", s.interval))

	// Main service loop
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Service shutdown requested")
			return err
		case <-ticker.C:
			// TODO: Replace this with your actual service logic
			// This is where developers should implement their business logic
			s.logger.Info("service running")
		}
	}
}

// healthHandler handles health check requests.
func (s *Service) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok\n")
}

// readinessHandler handles readiness check requests.
func (s *Service) readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ready\n")
}