# TRD Compliance Matrix

| Req ID | Description | Priority | Code Location(s) | Test(s) | Status |
|--------|-------------|----------|------------------|---------|--------|
| TRD-1  | Provide headless service example (runs, logs, provides metrics) | HIGH | cmd/{{.ProjectName}}/, pkg/{{.ProjectName}}/ | All tests | PLANNED |
| TRD-2  | Use Cobra for CLI | HIGH | cmd/{{.ProjectName}}/main.go | Integration tests | PLANNED |
| TRD-3  | Use Viper for config with automatic env | HIGH | pkg/{{.ProjectName}}/config.go | pkg/{{.ProjectName}}/config_test.go | PLANNED |
| TRD-4  | Use Prometheus for metrics | HIGH | pkg/{{.ProjectName}}/metrics.go | pkg/{{.ProjectName}}/metrics_test.go | PLANNED |
| TRD-5  | Serve metrics for requests, errors, durations | HIGH | pkg/{{.ProjectName}}/metrics.go | pkg/{{.ProjectName}}/metrics_test.go | PLANNED |
| GR-1   | Viper Automatic Env binding | HIGH | pkg/{{.ProjectName}}/config.go | pkg/{{.ProjectName}}/config_test.go | PLANNED |
| GR-2   | Coding Standards & Lint compliance | HIGH | All .go files | CI pipeline | PLANNED |
| GR-3   | MVC-ish Pattern enforcement | HIGH | pkg/{{.ProjectName}}/ (Models), cmd/{{.ProjectName}}/ (Views) | All tests | PLANNED |
| GR-4   | Globals & init for Cobra/Viper/Prometheus | MEDIUM | cmd/{{.ProjectName}}/main.go | N/A | PLANNED |
| GR-7   | Metrics port {{.DefaultServerPort}} with /metrics, /healthz, /readyz | HIGH | pkg/{{.ProjectName}}/server.go | pkg/{{.ProjectName}}/server_test.go | PLANNED |

## Service Shape Decision

Based on TRD analysis:
- **Service Type**: Headless service (no external clients beyond Prometheus)
- **Required Transports**: None (CLI for control only)  
- **Required Ports**: {{.DefaultServerPort}} (metrics, health endpoints)
- **Optional Ports**: None needed
- **Key Requirements**: Simple, runs, logs, provides metrics

## Assumptions Made

1. **Simple Implementation**: TRD explicitly states "Don't get elaborate. Keep it simple"
2. **No Database**: No data persistence requirements mentioned
3. **No External Integrations**: Headless service with minimal functionality
4. **Metrics Focus**: Primary function is to demonstrate metrics collection
5. **Example Service**: For code generation tool demonstration purposes