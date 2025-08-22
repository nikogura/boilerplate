# {{.ProjectName}} Service - Design Document

## Overview

Simple headless Go service for code generation tool demonstration. Runs, logs, and provides Prometheus metrics on port {{.DefaultServerPort}}.

## Architecture

```
┌─────────────────────────────────────────┐
│         {{.ProjectName}} Service       │
├─────────────────────────────────────────┤
│  HTTP Server (:{{.DefaultServerPort}})                    │
│  ├── /metrics (Prometheus)              │
│  ├── /healthz (Health Check)            │
│  └── /readyz (Readiness Check)          │
├─────────────────────────────────────────┤
│  Simple Worker Loop                     │
│  ├── Periodic Tasks                     │
│  └── Metrics Collection                 │
└─────────────────────────────────────────┘
```

## Package Layout (MVC-ish Pattern)

**Models (Libraries) - pkg/{{.ProjectName}}/**
- General-purpose, reusable business logic
- No transport assumptions
- Exhaustively tested

```
pkg/{{.ProjectName}}/
├── config.go          # Configuration management
├── metrics.go         # Metrics collection  
├── server.go          # HTTP server for metrics
├── service.go         # Core service logic
├── *_test.go          # Exhaustive tests
└── mocks/             # Generated mocks
```

**Views (Applications) - cmd/{{.ProjectName}}/**
- Application-specific interface
- Transport mapping only
- Integration tested

```
cmd/{{.ProjectName}}/
└── main.go            # Cobra CLI + service orchestration
```

## Data Models

### Config
```go
type Config struct {
    Server  ServerConfig
    Logging LoggingConfig  
    Metrics MetricsConfig
}
```

### Service
Simple service that:
1. Starts HTTP server on :{{.DefaultServerPort}}
2. Runs periodic work loop
3. Collects and exposes metrics
4. Handles graceful shutdown

## MVC Layering

**Models (pkg/{{.ProjectName}}/):**
- Configuration loading/validation
- Metrics collection interfaces
- HTTP server implementation
- Service orchestration logic
- All business logic with no transport assumptions

**Views (cmd/{{.ProjectName}}/):**
- CLI argument parsing (Cobra)
- Application startup/shutdown
- Signal handling
- Maps CLI to Model calls

## Validation & Error Handling

**Views handle:**
- CLI argument validation
- Configuration file errors
- Startup/shutdown errors

**Models handle:**  
- Configuration validation
- Service operation errors
- Clean error types for Views

## Context & Timeouts

- HTTP server: 30s read/write timeouts
- Graceful shutdown: 30s timeout
- All I/O operations use context.Context

## Viper Environment Keys + Defaults

```
SERVER_PORT={{.DefaultServerPort}}                    # HTTP server port
SERVER_READ_TIMEOUT=30s             # Read timeout  
SERVER_WRITE_TIMEOUT=30s            # Write timeout
LOGGING_LEVEL=info                  # Log level
LOGGING_FORMAT=json                 # Log format
METRICS_ENABLED=true                # Enable metrics
WORKER_INTERVAL=10s                 # Worker loop interval
```

## Observability

**Structured Logging (Zap):**
- JSON format (configurable)
- Structured fields
- Request IDs and durations

**Prometheus Metrics:**
- `{{.ProjectName}}_requests_total{endpoint,method}`
- `{{.ProjectName}}_request_errors_total{endpoint,method,error_type}`  
- `{{.ProjectName}}_request_duration_seconds{endpoint,method}`
- `{{.ProjectName}}_worker_operations_total{operation}`
- Standard Go/process metrics

**Health Endpoints:**
- `/metrics` - Prometheus metrics
- `/healthz` - Liveness probe  
- `/readyz` - Readiness probe

## Testing Strategy

**Models (≥85% coverage):**
- Exhaustive unit testing
- All error paths
- Table-driven tests
- No network dependencies
- Mock external interfaces

**Views (≥70% coverage):**
- Integration testing
- CLI execution tests
- End-to-end scenarios
- Signal handling

**Test Organization:**
```
pkg/{{.ProjectName}}/
├── config_test.go
├── metrics_test.go  
├── server_test.go
├── service_test.go
└── mocks/
    └── generated_mocks.go

cmd/{{.ProjectName}}/
└── main_test.go       # Integration tests
```

## Security Considerations

- No secrets in logs
- Input validation at boundaries  
- Timeouts prevent DoS
- Structured error responses
- No sensitive data exposure

## Deployment

**Single Binary:**
- Self-contained executable
- Configuration via environment variables
- Graceful shutdown on SIGTERM/SIGINT

**Resource Requirements:**
- Minimal CPU/memory
- Single port (:{{.DefaultServerPort}})
- No external dependencies