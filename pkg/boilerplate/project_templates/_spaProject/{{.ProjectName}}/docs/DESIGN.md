# {{.ProjectName}} Service Design

## Overview

This service is a simplified Single-Page Application (SPA) server designed as a boilerplate template. It serves static assets from an embedded filesystem while providing metrics, logging, and OIDC authentication.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        {{.ProjectName}} Service                     │
├─────────────────────────────────────────────────────────────────┤
│  Port :9999 (HTTP/SPA)      │  Port :8080 (Metrics)            │
│  ┌─────────────────────────┐ │ ┌─────────────────────────────────┐│
│  │ Static Asset Handler    │ │ │ Prometheus Metrics              ││
│  │ - Embedded FS          │ │ │ - /metrics                      ││
│  │ - SPA Routing          │ │ │ - /healthz                      ││
│  │ - OIDC Auth            │ │ │ - /readyz                       ││
│  └─────────────────────────┘ │ └─────────────────────────────────┘│
├─────────────────────────────────────────────────────────────────┤
│                       Core Libraries (pkg/)                    │
│  ┌─────────────────┐ ┌─────────────┐ ┌─────────────────────────┐│
│  │ Config (Viper)  │ │ Auth (OIDC) │ │ UI (Embedded FS)       ││
│  │ - Auto Env      │ │ - OAuth2    │ │ - Static Assets        ││
│  │ - Validation    │ │ - JWT       │ │ - SPA Handler          ││
│  └─────────────────┘ └─────────────┘ └─────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

## Package Layout

Following MVC-ish pattern per nikogura.com standards:

### Models (pkg/)
- `pkg/{{.ProjectPackageName}}/`: Core business logic, configuration, servers
  - `config.go`: Viper-based configuration with automatic env
  - `server.go`: HTTP server implementation
  - `metrics.go`: Prometheus metrics collection
  - `logging.go`: Structured logging with zap
- `pkg/auth/`: OIDC authentication (reusable)
- `pkg/ui/`: Embedded filesystem and SPA handlers (reusable)

### Views (cmd/)
- `cmd/root.go`: Cobra CLI root command
- `cmd/server.go`: Server startup command
- `main.go`: Entry point

## Data Models

### Configuration
```go
type Config struct {
    ServerPort    int    `mapstructure:"server_port"`
    MetricsPort   int    `mapstructure:"metrics_port"`
    LogLevel      string `mapstructure:"log_level"`
    OIDCClientID  string `mapstructure:"oidc_client_id"`
    // ... other fields
}
```

## Validation
- Input validation at View layer (cmd/)
- Config validation at startup
- OIDC token validation

## Error Taxonomy
- Configuration errors: Fatal exit with code 100
- Runtime errors: Logged with structured context
- Auth errors: HTTP 401/403 responses
- Server errors: HTTP 500 responses with request IDs

## Context & Timeouts
- All I/O operations use context.Context
- HTTP server has configurable timeouts
- OIDC operations have reasonable timeouts

## Viper Configuration

### Environment Variables
All configuration via environment variables with defaults:
```bash
{{.ProjectEnvPrefix}}_SERVER_PORT=9999      # HTTP server port
{{.ProjectEnvPrefix}}_METRICS_PORT=8080     # Metrics server port
{{.ProjectEnvPrefix}}_LOG_LEVEL=info        # Logging level
{{.ProjectEnvPrefix}}_OIDC_CLIENT_ID=""     # OIDC client ID
{{.ProjectEnvPrefix}}_OIDC_CLIENT_SECRET="" # OIDC client secret
{{.ProjectEnvPrefix}}_OIDC_ISSUER_URL=""    # OIDC issuer URL
{{.ProjectEnvPrefix}}_OIDC_REDIRECT_URL=""  # OIDC redirect URL
```

### Defaults
```go
viper.SetDefault("server_port", 9999)
viper.SetDefault("metrics_port", 8080)
viper.SetDefault("log_level", "info")
// ... etc
```

## Observability

### Structured Logging
- JSON format via zap
- Request IDs for correlation
- Configurable levels
- No PII logging

### Prometheus Metrics
Exposed on `:8080/metrics`:
```
{{.ProjectPackageName}}_http_requests_total{method,status}
{{.ProjectPackageName}}_http_request_duration_seconds{method}
{{.ProjectPackageName}}_server_start_time_seconds
```

### Health Endpoints
- `/healthz`: Basic health check
- `/readyz`: Readiness check (includes auth provider connectivity)

## Testing Strategy

### Models (pkg/) - Exhaustive Testing (≥85% coverage)
- Unit tests with table-driven patterns
- Mock external dependencies (OIDC provider)
- Test all error conditions and timeouts
- Race condition testing with `-race` flag

### Views (cmd/) - Integration Testing (≥70% coverage)
- CLI command testing
- HTTP endpoint integration tests
- OIDC flow testing with test provider

### Test Organization
- Tests alongside production code
- `testdata/` for golden files
- Generated mocks in same package

## Security Considerations
- No secrets in logs
- HTTPS-only in production (configurable)
- CSRF protection via OIDC state parameter
- Input validation at boundaries
- Secure cookie settings

## Deployment
- Single binary with embedded assets
- Container-ready with health checks
- Configurable via environment variables only
- Graceful shutdown with context cancellation

## Template Considerations
This service is designed as a boilerplate template ("spa" type) with:
- Minimal complexity while demonstrating patterns
- Clear separation of concerns (MVC-ish)
- Reusable components (auth, config, ui)
- Standard observability practices
- Production-ready structure