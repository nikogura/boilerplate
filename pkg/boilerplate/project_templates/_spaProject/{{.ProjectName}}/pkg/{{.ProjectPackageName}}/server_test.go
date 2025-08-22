package {{.ProjectPackageName}}

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_StatusHandler(t *testing.T) {
	// Test the handler directly without creating full server to avoid metrics conflicts
	config := &Config{
		ServerAddress:  "127.0.0.1:0",
		MetricsAddress: "127.0.0.1:0",
		LogLevel:       "info",
	}

	server := &Server{
		config: config,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rr := httptest.NewRecorder()

	server.statusHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Contains(t, rr.Body.String(), "\"status\":\"ok\"")
	assert.Contains(t, rr.Body.String(), "\"service\":\"{{.ProjectName}}\"")
}

func TestServer_UserAPIHandler_NoAuth(t *testing.T) {
	// Test the handler directly
	server := &Server{}

	req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
	rr := httptest.NewRecorder()

	server.userAPIHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Contains(t, rr.Body.String(), "\"user\":\"anonymous\"")
}

// Server configuration and integration tests are simplified
// to avoid Prometheus metrics registration conflicts in test suite
