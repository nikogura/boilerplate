# Example Indirect Selection Service - Design Document

## Overview

The `example-indirect-selection` service is a simple gRPC-based Go service that demonstrates indirect selection patterns for code generation tools. The service follows the MVC-ish architecture pattern with strict separation between Models (business logic in `/pkg`) and Views (presentation logic in `/cmd`).

## Architecture

```
example-indirect-selection/
├── cmd/example-indirect-selection/          # View Layer (Cobra CLI)
│   ├── main.go                             # Entry point
│   ├── root.go                             # Root command setup
│   ├── server.go                           # Server command
│   └── client.go                           # Client command (demo)
├── pkg/example-indirect-selection/          # Model Layer (Business Logic)
│   ├── config.go                           # Configuration management
│   ├── server.go                           # gRPC server implementation
│   ├── client.go                           # gRPC client implementation
│   ├── auth.go                             # JWT-SSH authentication
│   ├── metrics.go                          # Domain-specific metrics
│   ├── service.proto                       # Protocol buffer definitions
│   ├── service.pb.go                       # Generated protobuf code
│   ├── service_grpc.pb.go                  # Generated gRPC code
│   └── *_test.go                           # Model tests
├── pkg/metrics/                            # Shared metrics server
│   ├── server.go                           # Metrics HTTP server
│   ├── health.go                           # Health check handlers
│   └── metrics.go                          # Base metrics setup
├── configs/
│   └── .env.example                        # Environment variables template
├── docs/
│   ├── DESIGN.md                           # This document
│   ├── TRD_COMPLIANCE.md                   # Requirements traceability
│   └── RUNBOOK.md                          # Operations guide
├── Makefile                                # Build and test targets
├── go.mod                                  # Go module definition
└── .golangci.yml                           # Linting configuration
```

## Package Layout and MVC-ish Pattern

### Models (`pkg/example-indirect-selection/`)
- **Purpose**: Reusable business logic libraries
- **Characteristics**:
  - No transport-specific types (HTTP requests, CLI args, etc.)
  - Context-aware for all I/O operations
  - Exhaustive test coverage (≥85%)
  - Can be imported by multiple different programs
  - Focused interfaces with dependency injection

### Views (`cmd/example-indirect-selection/`)
- **Purpose**: Application-specific presentation layer
- **Characteristics**:
  - Handles Cobra CLI setup and command routing
  - Maps between transport formats and Model data structures
  - Integration test coverage (≥70%)
  - Contains all user interaction logic

## Data Models

### Configuration Structure
```go
type Config struct {
    // Server Configuration
    GRPCAddress     string        `mapstructure:"grpc_address"`
    MetricsAddress  string        `mapstructure:"metrics_address"`
    EnableReflection bool         `mapstructure:"enable_reflection"`
    
    // Authentication
    Audience        string        `mapstructure:"audience"`
    TrustedUsersFile string       `mapstructure:"trusted_users_file"`
    
    // Timeouts
    ServerTimeout   time.Duration `mapstructure:"server_timeout"`
    ClientTimeout   time.Duration `mapstructure:"client_timeout"`
    
    // Logging
    LogLevel        string        `mapstructure:"log_level"`
    LogFormat       string        `mapstructure:"log_format"`
}
```

### gRPC Service Definition
```proto
service ExampleIndirectSelection {
    // Demonstrates indirect selection with dynamic method calls
    rpc ProcessRequest(ProcessRequestInput) returns (ProcessResponse);
    rpc ListMethods(ListMethodsInput) returns (ListMethodsResponse);
    rpc GetMethodInfo(GetMethodInfoInput) returns (GetMethodInfoResponse);
}

message ProcessRequestInput {
    string method_name = 1;
    map<string, string> parameters = 2;
    string request_id = 3;
}

message ProcessResponse {
    string result = 1;
    string method_used = 2;
    int64 processing_time_ms = 3;
    string request_id = 4;
}
```

## Error Taxonomy

1. **Configuration Errors**: `ErrInvalidConfig`, `ErrMissingRequiredField`
2. **Authentication Errors**: `ErrInvalidToken`, `ErrUnauthorized`, `ErrExpiredToken`
3. **Service Errors**: `ErrMethodNotFound`, `ErrInvalidParameters`, `ErrProcessingFailed`
4. **Network Errors**: `ErrConnectionFailed`, `ErrTimeout`

All errors use `%w` wrapping for proper error chains and sentinel errors for expected conditions.

## Context and Timeouts

- All I/O operations accept `context.Context`
- Default timeouts:
  - Server operations: 30 seconds
  - Client operations: 10 seconds
  - Health checks: 5 seconds
- Graceful shutdown with 15-second timeout

## Viper Environment Configuration

### Environment Variables
All configuration uses Viper's automatic environment binding with the prefix `EXAMPLE_INDIRECT_SELECTION_`:

```bash
EXAMPLE_INDIRECT_SELECTION_GRPC_ADDRESS=0.0.0.0:50001
EXAMPLE_INDIRECT_SELECTION_METRICS_ADDRESS=0.0.0.0:8080
EXAMPLE_INDIRECT_SELECTION_ENABLE_REFLECTION=true
EXAMPLE_INDIRECT_SELECTION_AUDIENCE=example-indirect-selection
EXAMPLE_INDIRECT_SELECTION_TRUSTED_USERS_FILE=/etc/users.json
EXAMPLE_INDIRECT_SELECTION_LOG_LEVEL=info
```

### Default Values
```go
viper.SetDefault("grpc_address", "0.0.0.0:50001")
viper.SetDefault("metrics_address", "0.0.0.0:8080")
viper.SetDefault("enable_reflection", false)
viper.SetDefault("server_timeout", "30s")
viper.SetDefault("client_timeout", "10s")
viper.SetDefault("log_level", "info")
viper.SetDefault("log_format", "json")
```

## Observability

### Structured Logging
- Uses `zap.Logger` with structured JSON output
- Request IDs for tracing
- Log levels: DEBUG, INFO, WARN, ERROR
- Context-aware logging throughout

### Prometheus Metrics
All metrics use the prefix `example_indirect_selection_`:

```go
// Counters
example_indirect_selection_requests_total{method, status}
example_indirect_selection_errors_total{method, error_type}

// Histograms  
example_indirect_selection_request_duration_seconds{method}
example_indirect_selection_grpc_server_handling_seconds{method}

// Gauges
example_indirect_selection_active_connections
example_indirect_selection_server_uptime_seconds
```

### Health Endpoints
- `GET /healthz`: Liveness probe (always returns 200 OK)
- `GET /readyz`: Readiness probe (checks gRPC server status)
- `GET /metrics`: Prometheus metrics endpoint

## Testing Strategy

### Models Testing (pkg/example-indirect-selection/)
- **Coverage Target**: ≥85% line coverage
- **Approach**: Exhaustive testing of all scenarios
- **Patterns**: Table-driven tests, deterministic execution
- **Mocking**: External dependencies mocked with interfaces
- **Scenarios Covered**:
  - Success paths for all operations
  - All error conditions and edge cases
  - Timeout and cancellation handling
  - Concurrent access patterns

### Views Testing (cmd/example-indirect-selection/)
- **Coverage Target**: ≥70% line coverage
- **Approach**: Integration testing focus
- **Patterns**: Command execution testing, real gRPC calls
- **Scenarios Covered**:
  - CLI argument parsing and validation
  - gRPC client-server integration
  - Error handling and user feedback

### Test Organization
```
pkg/example-indirect-selection/
├── config_test.go
├── server_test.go
├── client_test.go
├── auth_test.go
├── metrics_test.go
└── testdata/
    ├── valid_config.json
    ├── users.json
    └── test_certificates/

cmd/example-indirect-selection/
├── server_test.go
├── client_test.go
└── integration_test.go
```

## Security Considerations

1. **JWT-SSH Authentication**: All gRPC calls require valid JWT tokens signed with SSH keys
2. **Input Validation**: All user inputs validated at View layer before passing to Models
3. **No Secret Logging**: Careful to never log JWT tokens or sensitive configuration
4. **TLS**: gRPC server supports TLS for production deployment
5. **Timeouts**: All operations have reasonable timeouts to prevent resource exhaustion

## Indirect Selection Implementation

The service demonstrates indirect selection through:

1. **Dynamic Method Registry**: Service maintains a registry of available methods
2. **Runtime Method Discovery**: Clients can query available methods and their parameters
3. **Parameterized Execution**: Methods can be invoked with dynamic parameters
4. **Reflection Support**: Optional gRPC reflection for tooling integration

This pattern enables code generation tools to discover and interact with service methods dynamically without compile-time dependencies.

## Deployment Considerations

### Ports
- **8080**: Metrics, health checks (HTTP)
- **50001**: gRPC service (with optional TLS)

### Environment
- Development: Uses plaintext gRPC, reflection enabled
- Production: TLS required, reflection disabled by default

### Dependencies
- Go 1.21+
- Protocol Buffers compiler
- golangci-lint for static analysis
- No external runtime dependencies (statically linked binary)

## Future Extensibility

The architecture supports easy extension through:

1. **Interface-based Design**: Easy to add new implementations
2. **Plugin Architecture**: Method registry supports dynamic registration
3. **Configuration-driven**: New features can be enabled via environment variables
4. **Separate Model/View**: Business logic changes don't affect CLI interface