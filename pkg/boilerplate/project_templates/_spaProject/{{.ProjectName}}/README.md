# {{.ProjectName}}

{{.ProjectShortDesc}}

## Features

- 🚀 **Single-Page Application**: Serves static assets from embedded filesystem
- ⚙️ **Configuration**: Viper-based with automatic environment variable binding  
- 📊 **Metrics**: Prometheus metrics on :8080 (/metrics, /healthz, /readyz)
- 🔐 **Authentication**: Optional OIDC integration
- 📝 **Logging**: Structured JSON logging with zap
- 🧪 **Testing**: Comprehensive test coverage with mocks
- 🔧 **Standards**: Follows MVC-ish patterns and coding standards

## Quick Start

### Build and Run
```bash
make build
./{{.ProjectName}} server
```

### Build with Semantic Version
To inject a semantic version into the UI footer:
```bash
SEMVER=v1.2.3 make build
./{{.ProjectName}} server
```

The service will start:
- **Main application**: http://localhost:9999 
- **Metrics/health**: http://localhost:8080

### With Docker
```bash
make docker-build
docker run -p 9999:9999 -p 8080:8080 {{.ProjectName}}
```

## Configuration

Configure via environment variables (see `configs/.env.example`):

```bash
export EXAMPLE_SPA_LOG_LEVEL=debug
export EXAMPLE_SPA_OIDC_CLIENT_ID="your-google-client-id"
export EXAMPLE_SPA_OIDC_CLIENT_SECRET="your-google-client-secret"
./{{.ProjectName}} server
```

## Architecture

Following MVC-ish patterns:
- **Models** (`pkg/{{.ProjectPackageName}}/`): Reusable business logic
- **Views** (`cmd/`): Application-specific CLI and server interfaces
- **UI** (`pkg/ui/`): Embedded static assets with SPA routing

## Development

### Make Targets
```bash
make help          # Show available commands
make build         # Build binary
make test          # Run tests with coverage
make lint          # Run linters
make ci            # Full CI pipeline
make run-debug     # Run with debug logging
```

### Testing
```bash
make test          # Run all tests
make coverage      # View coverage report
go test -race ./...  # Race detection
```

## Endpoints

### Application (Port 9999)
- `GET /` - Single-Page Application
- `GET /api/status` - Service status
- `GET /api/user` - Current user info (requires auth if enabled)
- `GET /auth/login` - OIDC login (if auth enabled)

### Metrics (Port 8080) 
- `GET /metrics` - Prometheus metrics
- `GET /healthz` - Health check
- `GET /readyz` - Readiness check

## Metrics

- `example_spa_http_requests_total{method,route,status}` - HTTP request counter
- `example_spa_http_request_duration_seconds{method,route}` - Request duration
- `example_spa_server_start_time_seconds` - Server startup time

## Authentication

Optional OIDC authentication via Google or compatible providers:

1. Set up OAuth2 application  
2. Configure environment variables
3. Service automatically enables auth middleware

Static bearer tokens also supported for API access.

## Documentation

- [Design Document](docs/DESIGN.md) - Architecture and patterns
- [Runbook](docs/RUNBOOK.md) - Operations and troubleshooting
- [TRD Compliance](docs/TRD_COMPLIANCE.md) - Requirements traceability

## License

Apache License 2.0