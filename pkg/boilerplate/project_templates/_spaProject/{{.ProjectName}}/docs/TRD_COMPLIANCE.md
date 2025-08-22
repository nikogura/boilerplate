# TRD Compliance Matrix

This document maps each TRD requirement and General Requirement to code implementation and tests.

## TRD Requirements

| Req ID | Description | Priority | Code Location(s) | Test(s) | Status |
|--------|-------------|----------|------------------|---------|--------|
| TRD-1 | Create Single-Page Application from embedded filesystem | HIGH | pkg/ui/ui.go, pkg/ui/static/ | pkg/ui/ui_test.go | ✅ IMPLEMENTED |
| TRD-2 | Service must run, log, and provide metrics (keep it simple) | HIGH | cmd/server.go, pkg/{{.ProjectPackageName}}/server.go | pkg/{{.ProjectPackageName}}/server_test.go | ✅ IMPLEMENTED |
| TRD-3 | Use Cobra for CLI | HIGH | cmd/root.go, cmd/server.go | cmd/server_test.go | ✅ IMPLEMENTED |
| TRD-4 | Use Viper for config with automatic env | HIGH | pkg/{{.ProjectPackageName}}/config.go | pkg/{{.ProjectPackageName}}/config_test.go | ✅ IMPLEMENTED |
| TRD-5 | Use Prometheus for metrics (requests, errors, durations) | HIGH | pkg/{{.ProjectPackageName}}/metrics.go | pkg/{{.ProjectPackageName}}/metrics_test.go | ✅ IMPLEMENTED |
| TRD-6 | Service name "{{.ProjectName}}" as metrics prefix | HIGH | pkg/{{.ProjectPackageName}}/metrics.go | pkg/{{.ProjectPackageName}}/metrics_test.go | ✅ IMPLEMENTED |
| TRD-7 | Authentication via OIDC | HIGH | pkg/auth/oidc.go | pkg/auth/oidc_test.go | ✅ IMPLEMENTED |

## General Requirements

| Req ID | Description | Priority | Code Location(s) | Test(s) | Status |
|--------|-------------|----------|------------------|---------|--------|
| GR-1 | Viper Automatic Env | HIGH | pkg/{{.ProjectPackageName}}/config.go | pkg/{{.ProjectPackageName}}/config_test.go | ✅ IMPLEMENTED |
| GR-2 | Coding Standards & Lint | HIGH | .golangci.yml, all Go files | make lint | ✅ IMPLEMENTED |
| GR-3 | MVC-ish Pattern | HIGH | pkg/ (Models), cmd/ (Views) | All tests | ✅ IMPLEMENTED |
| GR-4 | Globals & init for Cobra/Viper/Prometheus | MED | cmd/root.go (with //nolint) | N/A | ✅ IMPLEMENTED |
| GR-5 | PostgreSQL Access Patterns | LOW | N/A (Simple example) | N/A | ❌ NOT APPLICABLE |
| GR-6 | No ORM | LOW | N/A (No DB) | N/A | ✅ COMPLIANT |
| GR-7 | Metrics Port (:8080 /metrics /healthz /readyz) | HIGH | pkg/{{.ProjectPackageName}}/metrics.go | pkg/{{.ProjectPackageName}}/metrics_test.go | ✅ IMPLEMENTED |
| GR-8 | gRPC Port (:50001) | LOW | N/A (HTTP only) | N/A | ❌ NOT APPLICABLE |
| GR-9 | HTTP/WS Port (:9999) | HIGH | pkg/{{.ProjectPackageName}}/server.go | pkg/{{.ProjectPackageName}}/server_test.go | ✅ IMPLEMENTED |
| GR-10 | Named Returns Pattern | HIGH | All function signatures | All tests | ✅ IMPLEMENTED |

## Implementation Summary

This {{.ProjectName}} service implements a minimal Single-Page Application server with the following key features:

1. **SPA Serving**: Serves static assets from embedded filesystem with client-side routing support
2. **Configuration**: Viper-based configuration with automatic environment variable binding
3. **Metrics**: Prometheus metrics on :8080 with basic request tracking
4. **Authentication**: OIDC integration for secure access
5. **Logging**: Structured logging with zap
6. **Health Checks**: Standard /healthz and /readyz endpoints

The service follows the MVC-ish pattern with:
- **Models** (pkg/{{.ProjectPackageName}}): Reusable business logic and configuration
- **Views** (cmd/): Application-specific CLI and server interfaces
- **Embedded Assets** (pkg/ui/): Static SPA files served via embedded filesystem

All requirements are satisfied while maintaining simplicity as specified in the TRD.