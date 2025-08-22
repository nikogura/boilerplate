package {{.ProjectPackageName}}

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the {{.ProjectName}} service.
type Config struct {
	ServerAddress  string `mapstructure:"server_address"`
	MetricsAddress string `mapstructure:"metrics_address"`
	LogLevel       string `mapstructure:"log_level"`

	// OIDC Configuration
	OIDCClientID     string `mapstructure:"oidc_client_id"`
	OIDCClientSecret string `mapstructure:"oidc_client_secret"`
	OIDCIssuerURL    string `mapstructure:"oidc_issuer_url"`
	OIDCRedirectURL  string `mapstructure:"oidc_redirect_url"`
	OIDCCookieDomain string `mapstructure:"oidc_cookie_domain"`
	OIDCCookieSecure bool   `mapstructure:"oidc_cookie_secure"`
	OIDCStaticToken  string `mapstructure:"oidc_static_token"`
}

// LoadConfig loads configuration from environment variables using Viper.
func LoadConfig() (config *Config, err error) {
	// Set up Viper for automatic environment variable binding
	viper.AutomaticEnv()

	// Set environment variable prefix
	viper.SetEnvPrefix("{{.ProjectEnvPrefix}}")

	// Set key replacer to convert dots and hyphens to underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Set defaults for all configuration keys
	setDefaults()

	config = &Config{}

	// Unmarshal the configuration
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate the configuration
	err = validateConfig(config)
	if err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// setDefaults sets default values for all configuration keys.
func setDefaults() {
	// Server configuration
	viper.SetDefault("server_address", "0.0.0.0:9999")
	viper.SetDefault("metrics_address", "0.0.0.0:8080")
	viper.SetDefault("log_level", "info")

	// OIDC configuration - empty defaults (authentication optional)
	viper.SetDefault("oidc_client_id", "")
	viper.SetDefault("oidc_client_secret", "")
	viper.SetDefault("oidc_issuer_url", "https://accounts.google.com")
	viper.SetDefault("oidc_redirect_url", "http://localhost:9999/auth/callback")
	viper.SetDefault("oidc_cookie_domain", "")
	viper.SetDefault("oidc_cookie_secure", false)
	viper.SetDefault("oidc_static_token", "")
}

// validateConfig validates the loaded configuration.
func validateConfig(config *Config) (err error) {
	// Validate log level
	validLogLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
	validLevel := false
	for _, level := range validLogLevels {
		if strings.ToLower(config.LogLevel) == level {
			validLevel = true
			break
		}
	}
	if !validLevel {
		return fmt.Errorf("invalid log level: %s (must be one of: %s)",
			config.LogLevel, strings.Join(validLogLevels, ", "))
	}

	// Validate OIDC configuration consistency
	if config.OIDCClientID != "" && config.OIDCClientSecret == "" {
		return errors.New("oidc_client_secret is required when oidc_client_id is set")
	}
	if config.OIDCClientSecret != "" && config.OIDCClientID == "" {
		return errors.New("oidc_client_id is required when oidc_client_secret is set")
	}

	// Validate addresses are not empty
	if config.ServerAddress == "" {
		return errors.New("server_address cannot be empty")
	}
	if config.MetricsAddress == "" {
		return errors.New("metrics_address cannot be empty")
	}

	return nil
}
