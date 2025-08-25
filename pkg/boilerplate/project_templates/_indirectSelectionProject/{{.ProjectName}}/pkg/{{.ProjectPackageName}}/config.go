//nolint:staticcheck // Package name matches service name requirement from prompt.xml
package {{.ProjectPackageName}}

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the {{.ProjectName}} service.
type Config struct {
	// Server Configuration
	GRPCAddress      string `mapstructure:"grpc_address"`
	MetricsAddress   string `mapstructure:"metrics_address"`
	EnableReflection bool   `mapstructure:"enable_reflection"`
	PlainText        bool   `mapstructure:"plaintext"`

	// Authentication
	Audience         string `mapstructure:"audience"`
	TrustedUsersFile string `mapstructure:"trusted_users_file"`

	// Timeouts
	ServerTimeout   time.Duration `mapstructure:"server_timeout"`
	ClientTimeout   time.Duration `mapstructure:"client_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`

	// Logging
	LogLevel  string `mapstructure:"log_level"`
	LogFormat string `mapstructure:"log_format"`
	Debug     bool   `mapstructure:"debug"`
}

// LoadConfig loads configuration from environment variables using Viper.
func LoadConfig() (config Config, err error) {
	// Set up Viper for automatic environment variable binding
	viper.SetEnvPrefix("{{.EnvPrefix}}")
	viper.AutomaticEnv()

	// Replace dots and hyphens with underscores for environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Set default values
	setDefaults()

	// Unmarshal into struct
	err = viper.Unmarshal(&config)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal config: %w", err)
		return config, err
	}

	// Validate configuration
	err = config.Validate()
	if err != nil {
		err = fmt.Errorf("invalid configuration: %w", err)
		return config, err
	}

	return config, err
}

// setDefaults sets default configuration values.
func setDefaults() {
	// Server defaults
	viper.SetDefault("grpc_address", "0.0.0.0:50001")
	viper.SetDefault("metrics_address", "0.0.0.0:8080")
	viper.SetDefault("enable_reflection", false)
	viper.SetDefault("plaintext", true)

	// Authentication defaults
	viper.SetDefault("audience", "{{.ProjectName}}")
	viper.SetDefault("trusted_users_file", "/etc/{{.ProjectName}}/users.json")

	// Timeout defaults
	viper.SetDefault("server_timeout", "30s")
	viper.SetDefault("client_timeout", "10s")
	viper.SetDefault("shutdown_timeout", "15s")

	// Logging defaults
	viper.SetDefault("log_level", "info")
	viper.SetDefault("log_format", "json")
	viper.SetDefault("debug", false)
}

// Validate validates the configuration values.
func (c Config) Validate() (err error) {
	if c.GRPCAddress == "" {
		return errors.New("grpc_address cannot be empty")
	}

	if c.MetricsAddress == "" {
		return errors.New("metrics_address cannot be empty")
	}

	if c.Audience == "" {
		return errors.New("audience cannot be empty")
	}

	if c.TrustedUsersFile == "" {
		return errors.New("trusted_users_file cannot be empty")
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[strings.ToLower(c.LogLevel)] {
		return fmt.Errorf("invalid log_level: %s (must be debug, info, warn, or error)", c.LogLevel)
	}

	// Validate log format
	validLogFormats := map[string]bool{
		"json": true,
		"text": true,
	}
	if !validLogFormats[strings.ToLower(c.LogFormat)] {
		return fmt.Errorf("invalid log_format: %s (must be json or text)", c.LogFormat)
	}

	// Validate timeouts
	if c.ServerTimeout <= 0 {
		return errors.New("server_timeout must be positive")
	}

	if c.ClientTimeout <= 0 {
		return errors.New("client_timeout must be positive")
	}

	if c.ShutdownTimeout <= 0 {
		return errors.New("shutdown_timeout must be positive")
	}

	return nil
}