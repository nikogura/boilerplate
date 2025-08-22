package {{.ProjectPackageName}}

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected func(*testing.T, *Config)
		wantErr  bool
	}{
		{
			name:    "default configuration",
			envVars: map[string]string{},
			expected: func(t *testing.T, cfg *Config) {
				assert.Equal(t, 8080, cfg.Server.Port)
				assert.Equal(t, 30*time.Second, cfg.Server.ReadTimeout)
				assert.Equal(t, 30*time.Second, cfg.Server.WriteTimeout)
				assert.Equal(t, "info", cfg.Logging.Level)
				assert.Equal(t, "json", cfg.Logging.Format)
				assert.True(t, cfg.Metrics.Enabled)
				assert.Equal(t, "{{.ProjectName}}", cfg.Metrics.Namespace)
				assert.Equal(t, 10*time.Second, cfg.Worker.Interval)
			},
		},
		{
			name: "custom configuration via env vars",
			envVars: map[string]string{
				"SERVER_PORT":         "9090",
				"SERVER_READ_TIMEOUT": "45s",
				"LOGGING_LEVEL":       "debug",
				"LOGGING_FORMAT":      "console",
				"METRICS_ENABLED":     "false",
				"WORKER_INTERVAL":     "5s",
			},
			expected: func(t *testing.T, cfg *Config) {
				assert.Equal(t, 9090, cfg.Server.Port)
				assert.Equal(t, 45*time.Second, cfg.Server.ReadTimeout)
				assert.Equal(t, "debug", cfg.Logging.Level)
				assert.Equal(t, "console", cfg.Logging.Format)
				assert.False(t, cfg.Metrics.Enabled)
				assert.Equal(t, 5*time.Second, cfg.Worker.Interval)
			},
		},
		{
			name: "invalid port",
			envVars: map[string]string{
				"SERVER_PORT": "70000",
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			envVars: map[string]string{
				"LOGGING_LEVEL": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid log format",
			envVars: map[string]string{
				"LOGGING_FORMAT": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid timeout",
			envVars: map[string]string{
				"SERVER_READ_TIMEOUT": "0s",
			},
			wantErr: true,
		},
		{
			name: "invalid worker interval",
			envVars: map[string]string{
				"WORKER_INTERVAL": "0s",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			cfg, err := LoadConfig()

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cfg)

			if tt.expected != nil {
				tt.expected(t, cfg)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid configuration",
			cfg: &Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "json",
				},
				Worker: WorkerConfig{
					Interval: 10 * time.Second,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid port - too low",
			cfg: &Config{
				Server: ServerConfig{Port: 0},
			},
			wantErr: true,
			errMsg:  "server.port must be between 1 and 65535",
		},
		{
			name: "invalid port - too high",
			cfg: &Config{
				Server: ServerConfig{Port: 70000},
			},
			wantErr: true,
			errMsg:  "server.port must be between 1 and 65535",
		},
		{
			name: "invalid read timeout",
			cfg: &Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  0,
					WriteTimeout: 30 * time.Second,
				},
			},
			wantErr: true,
			errMsg:  "server.read_timeout must be positive",
		},
		{
			name: "invalid write timeout",
			cfg: &Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 0,
				},
			},
			wantErr: true,
			errMsg:  "server.write_timeout must be positive",
		},
		{
			name: "invalid log level",
			cfg: &Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
				},
				Logging: LoggingConfig{
					Level:  "invalid",
					Format: "json",
				},
			},
			wantErr: true,
			errMsg:  "logging.level must be one of",
		},
		{
			name: "invalid log format",
			cfg: &Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "invalid",
				},
			},
			wantErr: true,
			errMsg:  "logging.format must be one of",
		},
		{
			name: "invalid worker interval",
			cfg: &Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "json",
				},
				Worker: WorkerConfig{
					Interval: 0,
				},
			},
			wantErr: true,
			errMsg:  "worker.interval must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.cfg)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}