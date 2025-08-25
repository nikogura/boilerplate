# Example Indirect Selection Service - Runbook

## Service Overview

The `example-indirect-selection` service is a gRPC-based demonstration service that showcases indirect method selection patterns for code generation tools. It provides dynamic method discovery and execution capabilities with JWT-SSH authentication.

## Quick Start

### Prerequisites

1. Go 1.21+
2. Protocol Buffers compiler (`protoc`)
3. SSH key pair for authentication

### Building

```bash
# Install dependencies and build
make deps
make build

# Or run full CI pipeline
make ci
```

### Running

1. Create a trusted users file:
```json
[
  {
    "username": "alice",
    "public_key": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExamplePublicKeyHere alice@example.com"
  }
]
```

2. Set environment variables (optional):
```bash
export EXAMPLE_INDIRECT_SELECTION_TRUSTED_USERS_FILE=/path/to/users.json
export EXAMPLE_INDIRECT_SELECTION_ENABLE_REFLECTION=true
export EXAMPLE_INDIRECT_SELECTION_PLAINTEXT=true  # For development only
```

3. Start the server:
```bash
./bin/example-indirect-selection server --users=/path/to/users.json
```

## Service Endpoints

### gRPC Server
- **Address**: `0.0.0.0:50001` (configurable)
- **Reflection**: Configurable via environment variable
- **Authentication**: JWT tokens signed with SSH keys

### Metrics & Health
- **Address**: `0.0.0.0:8080` (configurable)
- **Endpoints**:
  - `GET /metrics` - Prometheus metrics
  - `GET /healthz` - Liveness probe
  - `GET /readyz` - Readiness probe

## Available Methods

The service provides several example methods for indirect selection:

| Method | Description | Required Parameters | Optional Parameters |
|--------|-------------|-------------------|-------------------|
| `echo` | Returns input unchanged | `message` | None |
| `reverse` | Reverses a string | `input` | None |
| `uppercase` | Converts to uppercase | `input` | None |
| `timestamp` | Returns current timestamp | None | `format` (default: RFC3339) |

## Client Usage

### List Available Methods
```bash
./bin/example-indirect-selection client list-methods
```

### Get Method Information
```bash
./bin/example-indirect-selection client method-info echo
```

### Call a Method
```bash
./bin/example-indirect-selection client call echo --message="Hello World"
./bin/example-indirect-selection client call reverse --input="hello"
./bin/example-indirect-selection client call timestamp --format="2006-01-02 15:04:05"
```

## Configuration

All configuration is handled via environment variables with the prefix `EXAMPLE_INDIRECT_SELECTION_`:

### Server Configuration
```bash
EXAMPLE_INDIRECT_SELECTION_GRPC_ADDRESS=0.0.0.0:50001
EXAMPLE_INDIRECT_SELECTION_METRICS_ADDRESS=0.0.0.0:8080
EXAMPLE_INDIRECT_SELECTION_ENABLE_REFLECTION=false
EXAMPLE_INDIRECT_SELECTION_PLAINTEXT=false
```

### Authentication
```bash
EXAMPLE_INDIRECT_SELECTION_AUDIENCE=example-indirect-selection
EXAMPLE_INDIRECT_SELECTION_TRUSTED_USERS_FILE=/etc/users.json
```

### Timeouts
```bash
EXAMPLE_INDIRECT_SELECTION_SERVER_TIMEOUT=30s
EXAMPLE_INDIRECT_SELECTION_CLIENT_TIMEOUT=10s
EXAMPLE_INDIRECT_SELECTION_SHUTDOWN_TIMEOUT=15s
```

### Logging
```bash
EXAMPLE_INDIRECT_SELECTION_LOG_LEVEL=info
EXAMPLE_INDIRECT_SELECTION_LOG_FORMAT=json
EXAMPLE_INDIRECT_SELECTION_DEBUG=false
```

## Monitoring

### Prometheus Metrics

The service exposes the following metrics:

- `example_indirect_selection_requests_total{method, status}` - Total requests processed
- `example_indirect_selection_errors_total{method, error_type}` - Total errors encountered
- `example_indirect_selection_request_duration_seconds{method}` - Request processing duration

### Health Checks

- **Liveness** (`/healthz`): Always returns 200 OK if the service is running
- **Readiness** (`/readyz`): Returns 200 OK if the service is ready to serve requests

### Logs

Structured JSON logs include:
- Request IDs for tracing
- Method names and parameters
- Processing times
- Authentication events
- Error details

## Troubleshooting

### Common Issues

1. **Authentication Failures**
   - Verify SSH key is correctly formatted in users.json
   - Ensure SSH agent is running and has the key loaded
   - Check JWT token generation and audience matching

2. **Connection Issues**
   - Verify ports are available and not blocked by firewall
   - For development, use `--plaintext` flag to disable TLS
   - Check server logs for binding errors

3. **Method Not Found**
   - Use `list-methods` command to see available methods
   - Verify method name spelling and case sensitivity

4. **Configuration Issues**
   - Validate environment variables are properly set
   - Check file permissions for users.json
   - Verify file paths are absolute, not relative

### Debug Mode

Enable debug logging for more detailed output:
```bash
./bin/example-indirect-selection --debug server
```

### Development Setup

For local development, use plaintext connections:
```bash
export EXAMPLE_INDIRECT_SELECTION_PLAINTEXT=true
export EXAMPLE_INDIRECT_SELECTION_ENABLE_REFLECTION=true
./bin/example-indirect-selection server --users=./test-users.json
```

## Production Deployment

### Security Considerations

1. **TLS**: Always use TLS in production (disable `plaintext` mode)
2. **User Management**: Regularly audit and rotate SSH keys in users.json
3. **Network Security**: Use appropriate firewall rules to restrict access
4. **Logging**: Monitor logs for authentication failures and suspicious activity

### Resource Requirements

- **Memory**: ~50MB base + ~10MB per concurrent connection
- **CPU**: Low usage, scales with request volume
- **Disk**: Minimal, only for logs and configuration files
- **Network**: 
  - gRPC: TCP port 50001
  - Metrics: TCP port 8080

### Scaling

The service is stateless and can be horizontally scaled:
- Deploy multiple instances behind a load balancer
- Use health checks for service discovery
- Monitor metrics to determine scaling needs

## Support

For issues and questions:
- Check logs for error details
- Verify configuration against this runbook
- Test with debug mode enabled
- Review prometheus metrics for performance insights