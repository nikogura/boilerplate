//nolint:staticcheck // Package name matches service name requirement from prompt.xml
package {{.ProjectPackageName}}

import (
	"testing"
)

func TestNewMetrics(t *testing.T) {
	// Create a shared metrics instance for all tests to avoid duplicate registration
	metrics := getTestMetrics()

	if metrics == nil {
		t.Fatal("NewMetrics() returned nil")
	}

	if metrics.RequestsTotal == nil {
		t.Error("RequestsTotal metric is nil")
	}

	if metrics.ErrorsTotal == nil {
		t.Error("ErrorsTotal metric is nil")
	}

	if metrics.RequestDuration == nil {
		t.Error("RequestDuration metric is nil")
	}
}

func TestMetricsUsage(t *testing.T) {
	// Use the shared metrics instance to avoid duplicate registration
	metrics := getTestMetrics()

	// Test that we can increment counters
	metrics.RequestsTotal.WithLabelValues("test_method", "success").Inc()
	metrics.ErrorsTotal.WithLabelValues("test_method", "validation_error").Inc()

	// Test that we can observe durations
	metrics.RequestDuration.WithLabelValues("test_method").Observe(0.1)

	// If we get here without panics, the metrics work correctly
}

// Shared metrics instance to avoid duplicate registration in tests.
//
//nolint:gochecknoglobals // Needed to avoid Prometheus duplicate registration in tests
var testMetrics *Metrics

func getTestMetrics() *Metrics {
	if testMetrics == nil {
		testMetrics = NewMetrics()
	}
	return testMetrics
}
