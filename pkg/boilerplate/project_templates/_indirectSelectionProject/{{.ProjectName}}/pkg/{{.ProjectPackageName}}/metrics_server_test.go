//nolint:staticcheck // Package name matches service name requirement from prompt.xml
package {{.ProjectPackageName}}

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap/zaptest"
)

func TestNewMetricsServer(t *testing.T) {
	logger := zaptest.NewLogger(t)
	address := "localhost:8080"

	server := NewMetricsServer(address, logger)

	if server.Address != address {
		t.Errorf("Expected address %s, got %s", address, server.Address)
	}

	if server.Logger != logger {
		t.Errorf("Expected logger to be set")
	}
}

func TestHealthHandlers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	server := NewMetricsServer("localhost:8080", logger)

	tests := []struct {
		name     string
		handler  http.HandlerFunc
		path     string
		expected string
	}{
		{
			name:     "liveness handler",
			handler:  server.LivenessHandler,
			path:     "/healthz",
			expected: "alive",
		},
		{
			name:     "readiness handler",
			handler:  server.ReadinessHandler,
			path:     "/readyz",
			expected: "ready",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			tt.handler(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			// Check if response contains expected status
			body := w.Body.String()
			if !containsSubstring(body, tt.expected) {
				t.Errorf("Expected response to contain %s, got %s", tt.expected, body)
			}
		})
	}
}

// Helper function.
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr)) ||
			findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
