package {{.ProjectPackageName}}

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Server represents the HTTP server for metrics and health endpoints.
type Server struct {
	server  *http.Server
	logger  *zap.Logger
	metrics *Metrics
	config  *Config
}

// NewServer creates a new HTTP server.
func NewServer(cfg *Config, logger *zap.Logger, metrics *Metrics) (server *Server) {
	mux := http.NewServeMux()

	s := &Server{
		logger:  logger,
		metrics: metrics,
		config:  cfg,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			Handler:      mux,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
		},
	}

	// Register routes
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", s.metricsMiddleware("/healthz", "GET", s.healthzHandler))
	mux.HandleFunc("/readyz", s.metricsMiddleware("/readyz", "GET", s.readyzHandler))

	server = s
	return server
}

// Start starts the HTTP server.
func (s *Server) Start() (err error) {
	s.logger.Info("Starting HTTP server", zap.String("addr", s.server.Addr))

	err = s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		err = fmt.Errorf("HTTP server failed to start: %w", err)
		return err
	}

	err = nil
	return err
}

// Stop gracefully stops the HTTP server.
func (s *Server) Stop(ctx context.Context) (err error) {
	s.logger.Info("Stopping HTTP server")
	err = s.server.Shutdown(ctx)
	return err
}

// metricsMiddleware wraps HTTP handlers with metrics collection.
func (s *Server) metricsMiddleware(endpoint, method string, next http.HandlerFunc) (handler http.HandlerFunc) {
	handler = func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Record request
		if s.metrics != nil {
			s.metrics.RecordRequest(endpoint, method)
		}

		// Execute handler
		next(w, r)

		// Record duration
		if s.metrics != nil {
			duration := time.Since(start).Seconds()
			s.metrics.RecordRequestDuration(endpoint, method, duration)
		}

		s.logger.Debug("HTTP request handled",
			zap.String("endpoint", endpoint),
			zap.String("method", method),
			zap.Duration("duration", time.Since(start)),
		)
	}
	return handler
}

// healthzHandler handles liveness probe requests.
func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := fmt.Sprintf(`{"status":"ok","timestamp":"%s"}`, time.Now().UTC().Format(time.RFC3339))
	_, err := w.Write([]byte(response))
	if err != nil {
		s.logger.Error("Failed to write health response", zap.Error(err))
	}
}

// readyzHandler handles readiness probe requests.
func (s *Server) readyzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := fmt.Sprintf(`{"status":"ready","timestamp":"%s"}`, time.Now().UTC().Format(time.RFC3339))
	_, err := w.Write([]byte(response))
	if err != nil {
		s.logger.Error("Failed to write readiness response", zap.Error(err))
	}
}