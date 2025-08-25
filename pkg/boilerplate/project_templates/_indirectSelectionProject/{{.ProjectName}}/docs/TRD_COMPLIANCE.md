# TRD Compliance Matrix

## TRD Requirements

| Req ID | Description | Priority | Code Location(s) | Test(s) | Status |
|--------|-------------|----------|------------------|---------|--------|
| TRD-1 | Create gRPC client/server with indirect selection example | High | `pkg/example-indirect-selection/`, `cmd/example-indirect-selection/` | `pkg/example-indirect-selection/*_test.go` | Pending |
| TRD-2 | Use Cobra for CLI | High | `cmd/example-indirect-selection/main.go`, `cmd/root.go` | `cmd/*_test.go` | Pending |
| TRD-3 | Use Viper for config with automatic env | High | `pkg/example-indirect-selection/config.go` | `pkg/example-indirect-selection/config_test.go` | Pending |
| TRD-4 | Use Prometheus metrics (requests, errors, durations) | High | `pkg/example-indirect-selection/metrics.go`, `pkg/metrics/` | `pkg/metrics/*_test.go` | Pending |
| TRD-5 | Service name prefix "example-indirect-selection" for metrics | High | `pkg/example-indirect-selection/metrics.go` | `pkg/example-indirect-selection/metrics_test.go` | Pending |
| TRD-6 | JWT-SSH-Agent authentication | High | `pkg/example-indirect-selection/auth.go` | `pkg/example-indirect-selection/auth_test.go` | Pending |
| TRD-7 | gRPC reflection toggle via environment config | Medium | `pkg/example-indirect-selection/server.go` | `pkg/example-indirect-selection/server_test.go` | Pending |
| TRD-8 | Follow all lint rules including namedreturns | High | All Go files | Lint checks | Pending |

## General Requirements

| Req ID | Description | Priority | Code Location(s) | Test(s) | Status |
|--------|-------------|----------|------------------|---------|--------|
| GR-1 | Viper Automatic Env with typed Config struct | High | `pkg/example-indirect-selection/config.go` | `pkg/example-indirect-selection/config_test.go` | Pending |
| GR-2 | Coding Standards & Lint compliance | High | All Go files | `make lint` | Pending |
| GR-3 | MVC-ish Pattern (Models in /pkg, Views in /cmd) | High | `pkg/example-indirect-selection/`, `cmd/example-indirect-selection/` | All tests | Pending |
| GR-4 | Cobra/Viper/Prometheus globals with //nolint | Medium | `cmd/root.go`, metrics files | Lint checks | Pending |
| GR-5 | PostgreSQL Access Patterns (if needed) | Low | N/A - Simple example service | N/A | Not Required |
| GR-6 | No ORM | N/A | N/A | N/A | Not Applicable |
| GR-7 | Metrics port 8080 with /metrics, /healthz, /readyz | High | `pkg/metrics/server.go` | `pkg/metrics/server_test.go` | Pending |
| GR-8 | gRPC port 50001 with reflection and health | High | `pkg/example-indirect-selection/server.go` | `pkg/example-indirect-selection/server_test.go` | Pending |
| GR-9 | HTTP/WS port 9999 (if applicable) | Low | N/A - gRPC only service | N/A | Not Required |
| GR-10 | Named Returns Pattern | High | All function signatures | Manual inspection + lint | Pending |

## Implementation Phases Status

| Phase | Description | Status | Notes |
|-------|-------------|--------|-------|
| A | Setup & Tooling | Pending | go.mod, .golangci.yml analysis complete |
| B | Configuration via Viper | Pending | Config struct with env binding |
| C | Servers & Ports | Pending | Metrics (8080), gRPC (50001) |
| D | MVC-ish Enforcement | Pending | pkg/ for models, cmd/ for views |
| E | Database Layer | Not Required | Simple example service |
| F | External Integrations | Not Required | Focus on gRPC example |
| G | Observability | Pending | Structured logging, Prometheus metrics |
| H | Validation & Data Flow | Pending | Input validation in views |
| I | Testing Strategy | Pending | Models ≥85%, Views ≥70% coverage |
| J | Lint & Static Analysis | Pending | golangci-lint clean, named returns |
| J-ITER | Continuous Validation | In Progress | Test/lint after every change |
| K | Docs & Runbook | Pending | This file, DESIGN.md, RUNBOOK.md |
| L | CI Pipeline | Pending | Use existing .github/workflows/ci.yml |

## Assumptions and Design Decisions

1. **Simple Service**: As per TRD requirement "don't get elaborate, keep it simple", this service will have minimal business logic and focus on demonstrating gRPC patterns with indirect selection.

2. **No Database**: Since TRD emphasizes simplicity and just "run, log, and provide metrics", no PostgreSQL integration will be implemented unless specifically required.

3. **JWT-SSH-Agent**: Will use `github.com/nikogura/jwt-ssh-agent-go` as specified in TRD-6.

4. **Indirect Selection**: The service will demonstrate indirect selection patterns through gRPC service methods that can be called dynamically.

5. **Metrics Focus**: Will implement core Prometheus metrics for request count, errors, and duration as specified.

## Traceability Notes

- All code locations will be updated as implementation progresses
- Test coverage will be verified with `go test -coverprofile=coverage.out`
- Lint compliance verified with existing `.golangci.yml` configuration
- Named returns pattern will be manually verified in addition to automated checks