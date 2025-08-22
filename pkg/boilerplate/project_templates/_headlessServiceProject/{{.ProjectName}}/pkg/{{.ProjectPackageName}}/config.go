package {{.ProjectPackageName}}

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config holds all configuration for the service.
type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Logging LoggingConfig `mapstructure:"logging"`
	Metrics MetricsConfig `mapstructure:"metrics"`
	Worker  WorkerConfig  `mapstructure:"worker"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// LoggingConfig holds logging configuration.
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// MetricsConfig holds metrics configuration.
type MetricsConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Namespace string `mapstructure:"namespace"`
}

// WorkerConfig holds worker loop configuration.
type WorkerConfig struct {
	Interval time.Duration `mapstructure:"interval"`
}

// LoadConfig loads configuration using Viper with automatic environment variable binding.
func LoadConfig() (cfg *Config, err error) {
	v := viper.New()

	// Set up environment variable handling
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	// Set defaults
	setDefaults(v)

	// Unmarshal into config struct
	var config Config
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	err = validateConfig(&config)
	if err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	cfg = &config
	return cfg, err
}

// setDefaults sets default values for all configuration keys.
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 30*time.Second)
	v.SetDefault("server.write_timeout", 30*time.Second)

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")

	// Metrics defaults
	v.SetDefault("metrics.enabled", true)
	v.SetDefault("metrics.namespace", "{{.ProjectName}}")

	// Worker defaults
	v.SetDefault("worker.interval", 10*time.Second)
}

// validateConfig validates the loaded configuration.
func validateConfig(cfg *Config) (err error) {
	// Validate server port
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		err = fmt.Errorf("server.port must be between 1 and 65535, got %d", cfg.Server.Port)
		return err
	}

	// Validate timeouts
	if cfg.Server.ReadTimeout <= 0 {
		err = errors.New("server.read_timeout must be positive")
		return err
	}
	if cfg.Server.WriteTimeout <= 0 {
		err = errors.New("server.write_timeout must be positive")
		return err
	}

	// Validate log level
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
		"dpanic": true, "panic": true, "fatal": true,
	}
	if !validLevels[cfg.Logging.Level] {
		err = errors.New("logging.level must be one of: debug, info, warn, error, dpanic, panic, fatal")
		return err
	}

	// Validate log format
	validFormats := map[string]bool{"json": true, "console": true}
	if !validFormats[cfg.Logging.Format] {
		err = errors.New("logging.format must be one of: json, console")
		return err
	}

	// Validate worker interval
	if cfg.Worker.Interval <= 0 {
		err = errors.New("worker.interval must be positive")
		return err
	}

	return err
}

// LogConfig logs the current configuration (without sensitive data).
func (c *Config) LogConfig(logger *zap.Logger) {
	logger.Info("Configuration loaded",
		zap.Int("server.port", c.Server.Port),
		zap.Duration("server.read_timeout", c.Server.ReadTimeout),
		zap.Duration("server.write_timeout", c.Server.WriteTimeout),
		zap.String("logging.level", c.Logging.Level),
		zap.String("logging.format", c.Logging.Format),
		zap.Bool("metrics.enabled", c.Metrics.Enabled),
		zap.String("metrics.namespace", c.Metrics.Namespace),
		zap.Duration("worker.interval", c.Worker.Interval),
	)
}