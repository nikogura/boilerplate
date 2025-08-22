package {{.ProjectPackageName}}

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewMetricsServer(t *testing.T) {
	logger := zap.NewNop()
	metrics := NewMetrics()

	server, err := NewMetricsServer("127.0.0.1:0", logger, metrics)
	require.NoError(t, err)
	assert.NotNil(t, server)
	assert.NotNil(t, server.server)
	assert.NotNil(t, server.logger)
	assert.NotNil(t, server.metrics)
}

func TestHealthzHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	healthzHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestReadyzHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr := httptest.NewRecorder()

	readyzHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "READY", rr.Body.String())
}

func TestResponseWriterWrapper(t *testing.T) {
	rr := httptest.NewRecorder()
	wrapper := &responseWriterWrapper{
		ResponseWriter: rr,
		statusCode:     http.StatusOK,
	}

	// Test WriteHeader
	wrapper.WriteHeader(http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, wrapper.statusCode)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Test Write
	data := []byte("test data")
	n, err := wrapper.Write(data)
	require.NoError(t, err)
	assert.Equal(t, len(data), n)
	assert.Equal(t, "test data", rr.Body.String())
}

func TestMetricsServer_ShutdownNilServer(t *testing.T) {
	logger := zap.NewNop()

	// Test shutdown with valid server
	httpServer := &http.Server{Addr: "127.0.0.1:0"}
	server := &MetricsServer{
		server: httpServer,
		logger: logger,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	// Error is expected since server wasn't started, but not a panic
	_ = err // We just want to test it doesn't panic
}
