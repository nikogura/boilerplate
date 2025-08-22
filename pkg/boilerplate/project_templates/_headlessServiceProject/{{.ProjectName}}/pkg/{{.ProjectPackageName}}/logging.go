package {{.ProjectPackageName}}

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap logger based on configuration.
func NewLogger(level, format string) (logger *zap.Logger, err error) {
	var config zap.Config

	switch strings.ToLower(format) {
	case "json":
		config = zap.NewProductionConfig()
	case "console":
		config = zap.NewDevelopmentConfig()
	default:
		err = fmt.Errorf("unsupported log format: %s", format)
		return logger, err
	}

	// Parse and set log level
	var zapLevel zapcore.Level
	zapLevel, err = zapcore.ParseLevel(level)
	if err != nil {
		err = fmt.Errorf("invalid log level %s: %w", level, err)
		return logger, err
	}
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	// Build logger
	logger, err = config.Build()
	if err != nil {
		err = fmt.Errorf("failed to build logger: %w", err)
		return logger, err
	}

	return logger, err
}