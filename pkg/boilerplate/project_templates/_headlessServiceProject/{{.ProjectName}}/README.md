# {{.ProjectName}}

{{.ProjectLongDesc}}

## Description

{{.ProjectShortDesc}}

## Usage

### Basic Usage

```bash
./{{.ProjectName}}
```

### Configuration

The service can be configured using:

1. Command line flags
2. Environment variables (prefixed with `{{.EnvPrefix}}_`)
3. Configuration file (YAML format)

### Command Line Options

```bash
./{{.ProjectName}} --help
```

### Environment Variables

- `{{.EnvPrefix}}_SERVER_ADDRESS` - Server bind address (default: 0.0.0.0)
- `{{.EnvPrefix}}_SERVER_PORT` - Server port (default: {{.DefaultServerPort}})
- `{{.EnvPrefix}}_LOG_LEVEL` - Log level (debug, info, warn, error) (default: info)

### Configuration File

Create a `.{{.ProjectName}}.yaml` file in your home directory or current directory:

```yaml
server:
  address: 0.0.0.0
  port: {{.DefaultServerPort}}
log:
  level: info
```

## API Endpoints

- `GET /health` - Health check endpoint
- `GET /api/v1/status` - Service status and configuration

## Building

```bash
go build -o {{.ProjectName}} .
```

## Development

```bash
go run . --log-level=debug
```

## Author

{{.MaintainerName}} <{{.MaintainerEmail}}>