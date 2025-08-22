# {{.ProjectName}} Service Runbook

## Service Overview

{{.ProjectName}} is a minimal Single-Page Application server that demonstrates standard patterns for Go services. It serves static assets from an embedded filesystem and provides metrics, logging, and optional OIDC authentication.

## Quick Start

### Build and Run
```bash
make build
./{{.ProjectName}} server
```

### With Docker
```bash
make docker-build
docker run -p 9999:9999 -p 8080:8080 {{.ProjectName}}
```

## Configuration

### Environment Variables

All configuration is done via environment variables with the `{{.ProjectEnvPrefix}}_` prefix:

| Variable | Default | Description |
|----------|---------|-------------|
| `{{.ProjectEnvPrefix}}_SERVER_ADDRESS` | `0.0.0.0:9999` | HTTP server bind address |
| `{{.ProjectEnvPrefix}}_METRICS_ADDRESS` | `0.0.0.0:8080` | Metrics server bind address |
| `{{.ProjectEnvPrefix}}_LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| `{{.ProjectEnvPrefix}}_OIDC_CLIENT_ID` | `""` | OAuth2 client ID (optional) |
| `{{.ProjectEnvPrefix}}_OIDC_CLIENT_SECRET` | `""` | OAuth2 client secret (optional) |
| `{{.ProjectEnvPrefix}}_OIDC_ISSUER_URL` | `https://accounts.google.com` | OIDC provider issuer |
| `{{.ProjectEnvPrefix}}_OIDC_REDIRECT_URL` | `http://localhost:9999/auth/callback` | OAuth2 callback URL |
| `{{.ProjectEnvPrefix}}_OIDC_COOKIE_SECURE` | `false` | Use secure cookies (set true for HTTPS) |
| `{{.ProjectEnvPrefix}}_OIDC_STATIC_TOKEN` | `""` | Static bearer token for API access |

### Example Configuration
```bash
export {{.ProjectEnvPrefix}}_LOG_LEVEL=debug
export {{.ProjectEnvPrefix}}_OIDC_CLIENT_ID="your-client-id.apps.googleusercontent.com"
export {{.ProjectEnvPrefix}}_OIDC_CLIENT_SECRET="your-client-secret"
```

## Service Endpoints

### Main Application (Port 9999)
- **/** - Single-Page Application (HTML/CSS/JS)
- **/api/user** - Current user information (JSON)
- **/api/status** - Service status (JSON)
- **/auth/login** - Initiate OIDC login (if auth enabled)
- **/auth/callback** - OIDC callback handler
- **/auth/logout** - Logout endpoint

### Metrics/Health (Port 8080)
- **/metrics** - Prometheus metrics
- **/healthz** - Health check (returns 200 OK)
- **/readyz** - Readiness check (returns 200 READY)

## Authentication

### OIDC Setup (Optional)
1. Configure OAuth2 application in your provider (Google, etc.)
2. Set callback URL to: `http://your-domain:9999/auth/callback`
3. Set environment variables:
   ```bash
   export {{.ProjectEnvPrefix}}_OIDC_CLIENT_ID="your-client-id"
   export {{.ProjectEnvPrefix}}_OIDC_CLIENT_SECRET="your-client-secret"
   ```
4. Restart service

### Static Token Access (Optional)
For API access without web authentication:
```bash
export {{.ProjectEnvPrefix}}_OIDC_STATIC_TOKEN="your-secret-token"
```

Access API with:
```bash
curl -H "Authorization: Bearer your-secret-token" http://localhost:9999/api/user
```

## Monitoring

### Prometheus Metrics
- `{{.ProjectPackageName}}_http_requests_total{method,route,status}` - Request counter
- `{{.ProjectPackageName}}_http_request_duration_seconds{method,route}` - Request duration histogram
- `{{.ProjectPackageName}}_server_start_time_seconds` - Server start timestamp

### Health Checks
```bash
# Health check
curl http://localhost:8080/healthz

# Readiness check  
curl http://localhost:8080/readyz

# Metrics
curl http://localhost:8080/metrics
```

## Troubleshooting

### Common Issues

#### Service won't start
1. Check if ports are available:
   ```bash
   netstat -tulpn | grep -E ':(8080|9999)'
   ```
2. Check configuration:
   ```bash
   ./{{.ProjectName}} server --help
   ```

#### Authentication issues
1. Verify OIDC credentials are correct
2. Check callback URL matches provider configuration
3. Enable debug logging:
   ```bash
   {{.ProjectEnvPrefix}}_LOG_LEVEL=debug ./{{.ProjectName}} server
   ```

#### High memory usage
1. Check embedded assets size
2. Monitor goroutine leaks via metrics

### Logging

Structured JSON logs to stdout:
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "info",
  "message": "Server started",
  "address": "0.0.0.0:9999"
}
```

### Performance Tuning

#### Resource Limits
- Memory: ~50MB base + static assets
- CPU: Minimal under normal load
- File descriptors: Standard Go HTTP server limits

#### Scaling
- Stateless service - can run multiple instances
- Load balance across instances
- Use external session storage for multi-instance auth

## Deployment

### Docker
```dockerfile
FROM your-base-image
COPY {{.ProjectName}} /app/
ENTRYPOINT ["/app/{{.ProjectName}}", "server"]
```

### Kubernetes
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ProjectName}}
spec:
  replicas: 2
  selector:
    matchLabels:
      app: {{.ProjectName}}
  template:
    metadata:
      labels:
        app: {{.ProjectName}}
    spec:
      containers:
      - name: {{.ProjectName}}
        image: {{.ProjectName}}:latest
        ports:
        - containerPort: 9999
          name: http
        - containerPort: 8080
          name: metrics
        env:
        - name: {{.ProjectEnvPrefix}}_LOG_LEVEL
          value: "info"
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8080
          initialDelaySeconds: 1
          periodSeconds: 10
```

## Development

### Local Development
```bash
# Run with debug logging
make run-debug

# Run tests
make test

# Run linting
make lint

# Full CI pipeline
make ci
```

### Testing
```bash
# Unit tests
go test ./...

# Coverage report
make coverage

# Race detection
go test -race ./...
```