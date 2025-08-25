//nolint:staticcheck // Package name matches service name requirement from prompt.xml
package {{.ProjectPackageName}}

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

// LivenessHandler handles liveness probe requests.
func (s *MetricsServer) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status: "alive",
		Time:   r.Header.Get("Date"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		s.Logger.Error("Failed to encode liveness response", zap.Error(err))
	}
}

// ReadinessHandler handles readiness probe requests.
func (s *MetricsServer) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status: "ready",
		Time:   r.Header.Get("Date"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		s.Logger.Error("Failed to encode readiness response", zap.Error(err))
	}
}
