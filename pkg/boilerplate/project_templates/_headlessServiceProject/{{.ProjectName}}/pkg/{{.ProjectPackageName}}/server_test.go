package {{.ProjectPackageName}}

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestNewServer(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Metrics: MetricsConfig{
			Enabled:   true,
			Namespace: "test",
		},
	}

	logger := zaptest.NewLogger(t)
	metrics := NewMetricsWithRegisterer(cfg.Metrics.Namespace, prometheus.NewRegistry())

	server := NewServer(cfg, logger, metrics)

	assert.NotNil(t, server)
	assert.Equal(t, ":8080", server.server.Addr)
	assert.Equal(t, cfg.Server.ReadTimeout, server.server.ReadTimeout)
	assert.Equal(t, cfg.Server.WriteTimeout, server.server.WriteTimeout)
}

func TestHealthzHandler(t *testing.T) {
	cfg := &Config{
		Server:  ServerConfig{Port: 8080},
		Metrics: MetricsConfig{Enabled: true, Namespace: "test"},
	}

	logger := zaptest.NewLogger(t)
	metrics := NewMetricsWithRegisterer(cfg.Metrics.Namespace, prometheus.NewRegistry())
	server := NewServer(cfg, logger, metrics)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	server.healthzHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), `"status":"ok"`)
	assert.Contains(t, w.Body.String(), `"timestamp"`)
}

func TestReadyzHandler(t *testing.T) {
	cfg := &Config{
		Server:  ServerConfig{Port: 8080},
		Metrics: MetricsConfig{Enabled: true, Namespace: "test"},
	}

	logger := zaptest.NewLogger(t)
	metrics := NewMetricsWithRegisterer(cfg.Metrics.Namespace, prometheus.NewRegistry())
	server := NewServer(cfg, logger, metrics)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()

	server.readyzHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), `"status":"ready"`)
	assert.Contains(t, w.Body.String(), `"timestamp"`)
}

func TestMetricsEndpoint(t *testing.T) {
	cfg := &Config{
		Server:  ServerConfig{Port: 8080},
		Metrics: MetricsConfig{Enabled: true, Namespace: "test"},
	}

	logger := zaptest.NewLogger(t)
	metrics := NewMetricsWithRegisterer(cfg.Metrics.Namespace, prometheus.NewRegistry())
	server := NewServer(cfg, logger, metrics)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()

	server.server.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/plain")
	// Should contain some basic metrics
	assert.Contains(t, w.Body.String(), "go_info")
}

func TestMetricsMiddleware(t *testing.T) {
	cfg := &Config{
		Server:  ServerConfig{Port: 8080},
		Metrics: MetricsConfig{Enabled: true, Namespace: "test"},
	}

	logger := zaptest.NewLogger(t)
	metrics := NewMetricsWithRegisterer(cfg.Metrics.Namespace, prometheus.NewRegistry())
	server := NewServer(cfg, logger, metrics)

	// Test handler that the middleware wraps
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test"))
	}

	wrappedHandler := server.metricsMiddleware("/test", "GET", testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	wrappedHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test", w.Body.String())
}

func TestServerStopWithTimeout(t *testing.T) {
	cfg := &Config{
		Server:  ServerConfig{Port: 8080},
		Metrics: MetricsConfig{Enabled: true, Namespace: "test"},
	}

	logger := zaptest.NewLogger(t)
	metrics := NewMetricsWithRegisterer(cfg.Metrics.Namespace, prometheus.NewRegistry())
	server := NewServer(cfg, logger, metrics)

	// Test that Stop works with context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := server.Stop(ctx)
	// Should not error since server wasn't started
	assert.NoError(t, err)
}