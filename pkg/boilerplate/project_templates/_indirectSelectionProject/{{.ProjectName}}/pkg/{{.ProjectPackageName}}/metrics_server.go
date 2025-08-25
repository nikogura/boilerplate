//nolint:staticcheck // Package name matches service name requirement from prompt.xml
package {{.ProjectPackageName}}

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// MetricsServer provides HTTP endpoints for metrics and health checks.
type MetricsServer struct {
	Address string
	Logger  *zap.Logger
}

// NewMetricsServer creates a new metrics server.
func NewMetricsServer(address string, logger *zap.Logger) (server *MetricsServer) {
	server = &MetricsServer{
		Address: address,
		Logger:  logger,
	}

	return server
}

// StartServer starts the metrics HTTP server.
func (s *MetricsServer) StartServer(ctx context.Context) (err error) {
	http.HandleFunc("/healthz", s.LivenessHandler)
	http.HandleFunc("/readyz", s.ReadinessHandler)
	http.Handle("/metrics", promhttp.Handler())

	s.Logger.Info("Starting metrics server", zap.String("address", s.Address))

	err = http.ListenAndServe(s.Address, nil)
	if err != nil {
		s.Logger.Error("Metrics server failed", zap.Error(err))
		return err
	}

	return err
}
