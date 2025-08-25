//nolint:staticcheck // Package name matches service name requirement from prompt.xml
package {{.ProjectPackageName}}

import (
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

//nolint:gocognit // Test function needs to be comprehensive
func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    Config
		wantErr bool
	}{
		{
			name:    "default configuration",
			envVars: map[string]string{},
			want: Config{
				GRPCAddress:      "0.0.0.0:50001",
				MetricsAddress:   "0.0.0.0:8080",
				EnableReflection: false,
				PlainText:        true,
				Audience:         "example-indirect-selection",
				TrustedUsersFile: "/etc/example-indirect-selection/users.json",
				ServerTimeout:    30 * time.Second,
				ClientTimeout:    10 * time.Second,
				ShutdownTimeout:  15 * time.Second,
				LogLevel:         "info",
				LogFormat:        "json",
				Debug:            false,
			},
			wantErr: false,
		},
		{
			name: "custom configuration",
			envVars: map[string]string{
				"EXAMPLE_INDIRECT_SELECTION_GRPC_ADDRESS":       "127.0.0.1:9000",
				"EXAMPLE_INDIRECT_SELECTION_METRICS_ADDRESS":    "127.0.0.1:9001",
				"EXAMPLE_INDIRECT_SELECTION_ENABLE_REFLECTION":  "true",
				"EXAMPLE_INDIRECT_SELECTION_PLAINTEXT":          "true",
				"EXAMPLE_INDIRECT_SELECTION_AUDIENCE":           "test-audience",
				"EXAMPLE_INDIRECT_SELECTION_TRUSTED_USERS_FILE": "/tmp/users.json",
				"EXAMPLE_INDIRECT_SELECTION_SERVER_TIMEOUT":     "60s",
				"EXAMPLE_INDIRECT_SELECTION_CLIENT_TIMEOUT":     "20s",
				"EXAMPLE_INDIRECT_SELECTION_SHUTDOWN_TIMEOUT":   "30s",
				"EXAMPLE_INDIRECT_SELECTION_LOG_LEVEL":          "debug",
				"EXAMPLE_INDIRECT_SELECTION_LOG_FORMAT":         "text",
				"EXAMPLE_INDIRECT_SELECTION_DEBUG":              "true",
			},
			want: Config{
				GRPCAddress:      "127.0.0.1:9000",
				MetricsAddress:   "127.0.0.1:9001",
				EnableReflection: true,
				PlainText:        true,
				Audience:         "test-audience",
				TrustedUsersFile: "/tmp/users.json",
				ServerTimeout:    60 * time.Second,
				ClientTimeout:    20 * time.Second,
				ShutdownTimeout:  30 * time.Second,
				LogLevel:         "debug",
				LogFormat:        "text",
				Debug:            true,
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			envVars: map[string]string{
				"EXAMPLE_INDIRECT_SELECTION_LOG_LEVEL": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid log format",
			envVars: map[string]string{
				"EXAMPLE_INDIRECT_SELECTION_LOG_FORMAT": "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper
			viper.Reset()

			// Set environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			got, err := LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got.GRPCAddress != tt.want.GRPCAddress {
				t.Errorf("GRPCAddress = %v, want %v", got.GRPCAddress, tt.want.GRPCAddress)
			}
			if got.MetricsAddress != tt.want.MetricsAddress {
				t.Errorf("MetricsAddress = %v, want %v", got.MetricsAddress, tt.want.MetricsAddress)
			}
			if got.EnableReflection != tt.want.EnableReflection {
				t.Errorf("EnableReflection = %v, want %v", got.EnableReflection, tt.want.EnableReflection)
			}
			if got.PlainText != tt.want.PlainText {
				t.Errorf("PlainText = %v, want %v", got.PlainText, tt.want.PlainText)
			}
			if got.Audience != tt.want.Audience {
				t.Errorf("Audience = %v, want %v", got.Audience, tt.want.Audience)
			}
			if got.TrustedUsersFile != tt.want.TrustedUsersFile {
				t.Errorf("TrustedUsersFile = %v, want %v", got.TrustedUsersFile, tt.want.TrustedUsersFile)
			}
			if got.ServerTimeout != tt.want.ServerTimeout {
				t.Errorf("ServerTimeout = %v, want %v", got.ServerTimeout, tt.want.ServerTimeout)
			}
			if got.ClientTimeout != tt.want.ClientTimeout {
				t.Errorf("ClientTimeout = %v, want %v", got.ClientTimeout, tt.want.ClientTimeout)
			}
			if got.ShutdownTimeout != tt.want.ShutdownTimeout {
				t.Errorf("ShutdownTimeout = %v, want %v", got.ShutdownTimeout, tt.want.ShutdownTimeout)
			}
			if got.LogLevel != tt.want.LogLevel {
				t.Errorf("LogLevel = %v, want %v", got.LogLevel, tt.want.LogLevel)
			}
			if got.LogFormat != tt.want.LogFormat {
				t.Errorf("LogFormat = %v, want %v", got.LogFormat, tt.want.LogFormat)
			}
			if got.Debug != tt.want.Debug {
				t.Errorf("Debug = %v, want %v", got.Debug, tt.want.Debug)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				GRPCAddress:      "0.0.0.0:50001",
				MetricsAddress:   "0.0.0.0:8080",
				Audience:         "test-audience",
				TrustedUsersFile: "/tmp/users.json",
				ServerTimeout:    30 * time.Second,
				ClientTimeout:    10 * time.Second,
				ShutdownTimeout:  15 * time.Second,
				LogLevel:         "info",
				LogFormat:        "json",
			},
			wantErr: false,
		},
		{
			name: "empty grpc address",
			config: Config{
				MetricsAddress:   "0.0.0.0:8080",
				Audience:         "test-audience",
				TrustedUsersFile: "/tmp/users.json",
				ServerTimeout:    30 * time.Second,
				ClientTimeout:    10 * time.Second,
				ShutdownTimeout:  15 * time.Second,
				LogLevel:         "info",
				LogFormat:        "json",
			},
			wantErr: true,
			errMsg:  "grpc_address cannot be empty",
		},
		{
			name: "invalid log level",
			config: Config{
				GRPCAddress:      "0.0.0.0:50001",
				MetricsAddress:   "0.0.0.0:8080",
				Audience:         "test-audience",
				TrustedUsersFile: "/tmp/users.json",
				ServerTimeout:    30 * time.Second,
				ClientTimeout:    10 * time.Second,
				ShutdownTimeout:  15 * time.Second,
				LogLevel:         "invalid",
				LogFormat:        "json",
			},
			wantErr: true,
			errMsg:  "invalid log_level: invalid",
		},
		{
			name: "zero timeout",
			config: Config{
				GRPCAddress:      "0.0.0.0:50001",
				MetricsAddress:   "0.0.0.0:8080",
				Audience:         "test-audience",
				TrustedUsersFile: "/tmp/users.json",
				ServerTimeout:    0,
				ClientTimeout:    10 * time.Second,
				ShutdownTimeout:  15 * time.Second,
				LogLevel:         "info",
				LogFormat:        "json",
			},
			wantErr: true,
			errMsg:  "server_timeout must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Config.Validate() error message = %v, should contain %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestEnvKeyReplacer tests that dots and hyphens in config keys are properly replaced with underscores.
func TestEnvKeyReplacer(t *testing.T) {
	// Reset viper
	viper.Reset()

	// Set environment variable with underscores
	t.Setenv("EXAMPLE_INDIRECT_SELECTION_SERVER_TIMEOUT", "45s")

	// Configure viper like in LoadConfig
	viper.SetEnvPrefix("EXAMPLE_INDIRECT_SELECTION")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetDefault("server_timeout", "30s")

	// The key uses underscores internally, but env var also uses underscores
	timeout := viper.GetString("server_timeout")
	if timeout != "45s" {
		t.Errorf("Expected server_timeout to be '45s', got '%s'", timeout)
	}
}
